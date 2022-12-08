package gopdf

import (
	"math"
	"strings"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

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
	wrapped := q.wrapLines()
	if len(wrapped) == 0 {
		return nil
	}

	// set format of first font before saving graphics state
	if len(wrapped[0].ChunkWidths) != 0 {
		if err := wrapped[0].Chunks[0].setFontAndColor(page); err != nil {
			return err
		}
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

	// calculate total height
	var totalHeight float64
	for _, l := range wrapped {
		totalHeight += l.Height
	}

	// set y starting position
	top := wrapped[0].Height - wrapped[0].MaxTop
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
		top -= line.Height

		// shortcut for empty lines
		if len(line.Chunks) == 0 {
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
			if err := chunk.draw(page, left, top); err != nil {
				return err
			}
			left += line.ChunkWidths[j]
		}
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

// wrapLines returns the wrapped text considering line break, max width and max height
func (q *TextChunkBoxElement) wrapLines() []chunkLine {
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
						chunkWidth := w + spaceWidth + charSpacing
						currLine.Width += chunkWidth
						if currLine.Height < h {
							currLine.Height = h
						}
						if currLine.MaxTop < fontTop {
							currLine.MaxTop = fontTop
						}

						c := q.Chunks[i]
						c.Text = " " + word
						currLine.Chunks = append(currLine.Chunks, c)
						currLine.ChunkWidths = append(currLine.ChunkWidths, chunkWidth)
						continue
					}
				}

				// case 2: current line is empty and current word fits into line
				if w <= boxWidth {
					if word == "" {
						continue
					}

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
	for _, l := range q.wrapLines() {
		h += l.Height
	}
	return Pt(h)
}

// ShrinkToFit reduces the font size so that the entire text fits into the box
func (q *TextChunkBoxElement) ShrinkToFit() {
	// wrap lines
	orgW, orgH := q.Width, q.Height
	q.Width.Value = 0
	q.Height.Value = 0
	lines := q.wrapLines()
	q.Width, q.Height = orgW, orgH

	// check height
	f := 1.0
	h := q.Height.Pt()
	if h > 0 {
		// get total height
		var totalHeight float64
		for _, l := range lines {
			totalHeight += l.Height
		}

		// adapt if exceeds
		if totalHeight > h {
			f = h / totalHeight * 0.999
		}
	}

	// adapt to width
	w := q.Width.Pt()
	if w > 0 {
		// get max width
		var maxWidth float64
		for _, line := range lines {
			if maxWidth < line.Width {
				maxWidth = line.Width
			}
		}

		// adapt if exceeds
		if maxWidth > w {
			f2 := w / maxWidth * 0.999
			if f2 < f {
				f = f2
			}
		}
	}

	// adapt
	if f != 1 {
		for i := range q.Chunks {
			q.Chunks[i].FontSize *= f
		}
	}
}
