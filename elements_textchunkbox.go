package gopdf

import (
	"errors"
	"math"
	"strings"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

type TextChunk struct {
	Text          string
	Font          pdf.FontHandler
	FontSize      float64
	Color         Color
	OutlineColor  Color
	OutlineWidth  Length
	DashPattern   DashPattern
	Bold          bool
	Italic        bool
	Underline     bool
	StrikeThrough bool
	CharSpacing   Length
	TextScaling   float64
}

// TextChunkBoxElement is similar to TextBoxElement, but can have chunks with different format
type TextChunkBoxElement struct {
	Chunks []TextChunk

	Transparency    float64
	LineHeight      float64
	Left, Top       Length
	Width, Height   Length
	Rotate          float64
	TextAlign       HorizontalAlign
	VerticalAlign   VerticalAlign
	HeightBufferRel float64
}

// Build adds the element to the content stream
func (q *TextChunkBoxElement) Build(page *pdf.Page) error {
	// return if no chunks
	if len(q.Chunks) == 0 {
		return nil
	}

	// wrap text
	wrapped := q.wrappedChunks()
	if len(wrapped) == 0 {
		return nil
	}

	// graphics state
	page.GraphicsState_q()
	defer page.GraphicsState_Q()

	// set coordinate system
	y := float64(page.Data.MediaBox.URY) - q.Top.Pt()
	if q.Rotate == 0 {
		page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), y)
	} else {
		r := q.Rotate * math.Pi / 180
		page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), q.Left.Pt(), y)
	}

	// transparency
	if q.Transparency > 0 && q.Transparency <= 1 {
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - q.Transparency),
			"CA": types.Number(1 - q.Transparency),
		})
		page.GraphicsState_gs(n)
	}

	// auto fit font size?
	//if q.FontSize == -1 {
	//	defer func() { q.FontSize = -1 }()
	//
	//	// split by line breaks
	//	lines := strings.Split(strings.ReplaceAll(q.Text, "\r\n", "\n"), "\n")
	//
	//	// adapt to height
	//	h := q.Height.Pt()
	//	if h > 0 {
	//		fh := q.FontHeight().Pt() * float64(len(lines))
	//		q.FontSize *= h / fh * 0.999
	//	}
	//
	//	// adapt to width
	//	w := q.Width.Pt()
	//	if w > 0 {
	//		var maxWidth float64
	//		for _, line := range lines {
	//			v := q.Font.GetWidth(line, q.FontSize)
	//			if q.CharSpacing.Value != 0 {
	//				v += float64(len([]rune(line))-1) * q.CharSpacing.Pt()
	//			}
	//			if q.TextScaling != 0 {
	//				v *= q.TextScaling / 100
	//			}
	//			if v > maxWidth {
	//				maxWidth = v
	//			}
	//		}
	//		if maxWidth > w {
	//			q.FontSize *= w / maxWidth * 0.999
	//		}
	//	}
	//}

	// calculate total height
	var totalHeight float64
	for _, l := range wrapped {
		totalHeight += l.Height
	}

	// set y starting position
	top := -wrapped[0].MaxTop
	if q.Height.Value > 0 {
		switch q.VerticalAlign {
		case VerticalAlignMiddle:
			top -= (q.Height.Pt() - totalHeight) / 2
		case VerticalAlignBottom:
			top -= q.Height.Pt() - totalHeight
		}
	}

	// begin text
	page.TextObjects_BT()
	defer page.TextObjects_ET()

	// iterate over lines
	for _, line := range wrapped {
		// shortcut for empty lines
		if len(line.Chunks) == 0 {
			top -= line.Height
			continue
		}

		// set position
		left := 0.0
		switch q.TextAlign {
		case HorizontalAlignCenter:
			left += (q.Width.Pt() - line.Width) / 2
		case HorizontalAlignRight:
			left += q.Width.Pt() - line.Width
		}

		// iterate over text chunks in this line
		for j, chunk := range line.Chunks {
			// if no font or text given, ignore chunk
			if chunk.Font == nil {
				return errors.New("no font set")
			}

			// set position
			var c float64
			if chunk.Italic {
				c = 0.333
			}
			page.TextPosition_Tm(1, 0, c, 1, left, top)
			left += line.ChunkWidths[j]

			// set color and rendering mode
			color := chunk.Color
			if color == nil && chunk.OutlineColor == nil {
				color = ColorRGBBlack
			}
			if chunk.Bold && chunk.OutlineWidth.Value == 0 && color != nil && chunk.OutlineColor == nil {
				color.Build(page, false)
				color.Build(page, true)
				page.GraphicsState_w(chunk.FontSize * 0.05)
				page.TextState_Tr(types.RenderingModeFillAndStroke)
			} else {
				page.GraphicsState_w(chunk.OutlineWidth.Pt())
				switch {
				case color != nil && chunk.OutlineColor != nil:
					color.Build(page, false)
					chunk.OutlineColor.Build(page, true)
					page.TextState_Tr(types.RenderingModeFillAndStroke)
				case color != nil:
					color.Build(page, false)
					page.TextState_Tr(types.RenderingModeFill)
				case chunk.OutlineColor != nil:
					page.TextState_Tr(types.RenderingModeStroke)
					chunk.OutlineColor.Build(page, true)
				}
			}
			if err := chunk.DashPattern.Build(page); err != nil {
				return err
			}

			// text scaling / char spacing
			page.TextState_Tc(chunk.CharSpacing.Pt())
			if chunk.TextScaling == 0 {
				page.TextState_Tz(100)
			} else {
				page.TextState_Tz(chunk.TextScaling)
			}

			// set font and draw text
			page.TextState_Tf(chunk.Font, chunk.FontSize)
			page.TextShowing_Tj(reverseRTLString(chunk.Text))

			// underline/strike-through text
			if chunk.Underline {
				page.Path_re(left-line.ChunkWidths[j], top+chunk.Font.GetUnderlinePosition(chunk.FontSize), line.ChunkWidths[j], chunk.Font.GetUnderlineThickness(chunk.FontSize))
				page.Path_f()
			}
			if chunk.StrikeThrough {
				page.Path_re(
					left-line.ChunkWidths[j], top+chunk.Font.GetTop(chunk.FontSize)/3,
					line.ChunkWidths[j], chunk.Font.GetUnderlineThickness(chunk.FontSize),
				)
				page.Path_f()
			}
		}
		top -= line.Height
	}
	return nil
}

type chunkLine struct {
	MaxTop      float64
	Height      float64
	Width       float64
	Chunks      []TextChunk
	ChunkWidths []float64
}

// wrappedChunks returns the wrapped text considering line break, max width and max height
func (q *TextChunkBoxElement) wrappedChunks() []chunkLine {
	var chunkLines []chunkLine
	var currLine *chunkLine
	boxWidth := q.Width.Pt()

	// iterate over chunks
	for i, chunk := range q.Chunks {
		// calc spaceWidth in current chunk
		charSpacing := chunk.CharSpacing.Pt()
		if chunk.TextScaling != 0 {
			charSpacing *= chunk.TextScaling / 100
		}
		spaceWidth := chunk.getLineWidth(" ") + charSpacing

		// split chunk texts by line break
		textLines := strings.Split(strings.ReplaceAll(chunk.Text, "\r\n", "\n"), "\n")

		// iterate over lines within chunk
		for j, line := range textLines {
			// create new line
			if len(chunkLines) == 0 || j > 0 {
				chunkLines = append(chunkLines, chunkLine{})
				currLine = &chunkLines[len(chunkLines)-1]
			}

			// get width and height
			w := chunk.getLineWidth(line)
			h := q.LineHeight
			if h == 0 {
				h = chunk.FontHeight().Pt()
			}
			fontTop := chunk.Font.GetTop(chunk.FontSize)

			// does the entire line fit in the current chunk line?
			if boxWidth == 0 || currLine.Width+w <= boxWidth {
				currLine.Width += w
				if currLine.Height < h {
					currLine.Height = h
				}
				if currLine.MaxTop < fontTop {
					currLine.MaxTop = fontTop
				}

				if line != "" {
					c := q.Chunks[i]
					c.Text = line
					currLine.Chunks = append(currLine.Chunks, c)
					currLine.ChunkWidths = append(currLine.ChunkWidths, w)
				}
				continue
			}

			// check how many words fit in the current line
			words := strings.Split(line, " ")
			for _, word := range words {
				w := chunk.getLineWidth(word)

				// case 1: current line is not empty
				if len(currLine.Chunks) != 0 {
					// create new line?
					if currLine.Width+w+spaceWidth+charSpacing > boxWidth {
						chunkLines = append(chunkLines, chunkLine{})
						currLine = &chunkLines[len(chunkLines)-1]

					} else {
						currLine.Width += w + spaceWidth + charSpacing
						if currLine.Height < h {
							currLine.Height = h
						}
						if currLine.MaxTop < fontTop {
							currLine.MaxTop = fontTop
						}

						c := q.Chunks[i]
						c.Text = " " + word
						currLine.Chunks = append(currLine.Chunks, c)
						currLine.ChunkWidths = append(currLine.ChunkWidths, w)
						continue
					}
				}

				// case 2: current line is empty and current word fits into line
				if w <= boxWidth {
					currLine.Width += w
					if currLine.Height < h {
						currLine.Height = h
					}
					if currLine.MaxTop < fontTop {
						currLine.MaxTop = fontTop
					}

					c := q.Chunks[i]
					c.Text = word
					currLine.Chunks = append(currLine.Chunks, c)
					currLine.ChunkWidths = append(currLine.ChunkWidths, w)
					continue
				}

				// case 3: break word
				currLine.Height = h
				currLine.MaxTop = fontTop
				c := q.Chunks[i]
				c.Text = ""
				currLine.Chunks = append(currLine.Chunks, c)
				currLine.ChunkWidths = append(currLine.ChunkWidths, 0)
				currChunk := &currLine.Chunks[len(currLine.Chunks)-1]
				for _, c := range word {
					w := chunk.getLineWidth(string(c))
					if currLine.Width+w+charSpacing > boxWidth-0.1 {
						nc := q.Chunks[i]
						nc.Text = ""
						newChunkLine := chunkLine{
							Height:      h,
							MaxTop:      fontTop,
							Chunks:      []TextChunk{nc},
							ChunkWidths: []float64{0},
						}
						chunkLines = append(chunkLines, newChunkLine)
						currLine = &chunkLines[len(chunkLines)-1]
						currChunk = &currLine.Chunks[0]
					}

					currLine.Width += w + charSpacing
					currLine.ChunkWidths[len(currLine.ChunkWidths)-1] += w
					currChunk.Text += string(c)
				}
			}
		}
	}

	// check height
	h := q.Height.Pt() * (1 + q.HeightBufferRel)
	if h > 0 {
		for i, l := range chunkLines {
			h -= l.Height
			if h < 0 {
				chunkLines = chunkLines[:i]
				break
			}
		}
	}

	return chunkLines
}

// TextHeight returns the height of the text, accounting for line breaks and max width
func (q *TextChunkBoxElement) TextHeight() Length {
	var h float64
	for _, l := range q.wrappedChunks() {
		h += l.Height
	}
	return Pt(h)
}

// getLineWidth returns the width of the given text line consider font, fontsize, text-scaling, char spacing
func (q *TextChunk) getLineWidth(line string) float64 {
	if line == "" {
		return 0
	}
	v := q.Font.GetWidth(line, q.FontSize)
	if q.CharSpacing.Value != 0 {
		v += float64(len([]rune(line))-1) * q.CharSpacing.Pt()
	}
	if q.TextScaling != 0 {
		v *= q.TextScaling / 100
	}
	return v
}

// FontHeight returns the height of the font (bounding box y min to max)
func (q *TextChunk) FontHeight() Length {
	return Pt(q.Font.GetTop(q.FontSize) - q.Font.GetBottom(q.FontSize))
}
