package gopdf

import (
	"math"
	"strings"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// TextBoxElement is similar to TextElement, but can have a maximum width and height
type TextBoxElement struct {
	TextElement
	Width, Height Length
	VerticalAlign VerticalAlign
}

// Build adds the element to the content stream
func (q *TextBoxElement) Build(page *pdf.Page) error {
	// if no font or text given, ignore element
	if q.Font == nil || q.Text == "" {
		return nil
	}

	// auto fit font size?
	if q.FontSize == -1 {
		defer func() { q.FontSize = -1 }()

		// split by line breaks
		lines := strings.Split(strings.ReplaceAll(q.Text, "\r\n", "\n"), "\n")

		// adapt to height
		h := q.Height.Pt()
		if h > 0 {
			fh := q.FontHeight().Pt() * float64(len(lines))
			q.FontSize *= h / fh * 0.999
		}

		// adapt to width
		w := q.Width.Pt()
		if w > 0 {
			var maxWidth float64
			for _, line := range lines {
				v := q.Font.GetWidth(line, q.FontSize)
				if q.CharSpacing.Value != 0 {
					v += float64(len([]rune(line))-1) * q.CharSpacing.Pt()
				}
				if q.TextScaling != 0 {
					v *= q.TextScaling / 100
				}
				if v > maxWidth {
					maxWidth = v
				}
			}
			if maxWidth > w {
				q.FontSize *= w / maxWidth * 0.999
			}
		}
	}

	// set color
	color := q.Color
	if color == nil {
		color = ColorRGBBlack
	}
	color.Build(page, false)
	color.Build(page, true)

	// set bold or outline (bold is done via outline)
	if q.Bold && q.OutlineWidth.Value == 0 && q.OutlineColor == nil {
		page.GraphicsState_w(q.FontSize * 0.05)
		page.TextState_Tr(types.RenderingModeFillAndStroke)
	} else {
		page.GraphicsState_w(q.OutlineWidth.Pt())
		page.TextState_Tr(q.RenderMode)
		if q.OutlineColor != nil {
			q.OutlineColor.Build(page, true)
		}
	}

	// text scaling / char spacing
	page.TextState_Tc(q.CharSpacing.Pt())
	if q.TextScaling == 0 {
		page.TextState_Tz(100)
	} else {
		page.TextState_Tz(q.TextScaling)
	}

	// set font
	page.TextState_Tf(q.Font, q.FontSize)

	// set coordinate system
	page.GraphicsState_q()
	y := float64(page.Data.MediaBox.URY) - q.Top.Pt()
	if q.Rotate == 0 {
		page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), y)
	} else {
		r := q.Rotate * math.Pi / 180
		page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), q.Left.Pt(), y)
	}

	// begin text
	page.TextObjects_BT()

	// calculate some values needed below
	wrapped := q.wrappedText()
	lineHeight := q.lineHeight()
	var c float64
	if q.Italic {
		c = 0.333
	}

	top := -q.Font.GetTop(q.FontSize)
	if q.Height.Value > 0 {
		switch q.VerticalAlign {
		case VerticalAlignMiddle:
			top -= (q.Height.Pt() - lineHeight*float64(len(wrapped))) / 2
		case VerticalAlignBottom:
			top -= q.Height.Pt() - lineHeight*float64(len(wrapped))
		}
	}

	// iterate over lines
	for _, line := range wrapped {
		// shortcut for empty lines
		if line == "" {
			top -= lineHeight
			continue
		}

		// set position
		left := 0.0
		width := q.getLineWidth(line)
		switch q.TextAlign {
		case HorizontalAlignCenter:
			left += (q.Width.Pt() - width) / 2
		case HorizontalAlignRight:
			left += q.Width.Pt() - width
		}
		page.TextPosition_Tm(1, 0, c, 1, left, top)
		page.TextShowing_Tj(reverseRTLString(line))

		// underline/strike-through text
		if q.Underline {
			page.Path_re(left, top+q.Font.GetUnderlinePosition(q.FontSize), width, q.Font.GetUnderlineThickness(q.FontSize))
			page.Path_f()
		}
		if q.StrikeThrough {
			page.Path_re(
				left, top+q.Font.GetTop(q.FontSize)/3,
				width, q.Font.GetUnderlineThickness(q.FontSize),
			)
			page.Path_f()
		}

		top -= lineHeight
	}

	page.TextObjects_ET()
	page.GraphicsState_Q()
	return nil
}

// wrappedText returns the wrapped text considering line break, max width and max height
func (q *TextBoxElement) wrappedText() []string {
	// determine max number of lines
	var maxLines int
	if q.Height.Value > 0 {
		maxLines = int(q.Height.Pt() / q.lineHeight())
		if maxLines == 0 {
			return nil
		}
	}

	// split by line break
	lines := strings.Split(strings.ReplaceAll(q.Text, "\r\n", "\n"), "\n")

	// iterate over lines
	var res []string
	charSpacing := q.CharSpacing.Pt()
	if q.TextScaling != 0 {
		charSpacing *= q.TextScaling / 100
	}
	spaceWidth := q.getLineWidth(" ") + charSpacing
	for _, line := range lines {
		// check if max number of lines reached
		if maxLines > 0 && len(res) >= maxLines {
			break
		}

		// if width of line does not exceed max width, we can add the entire line and continue with the next line
		if q.Width.Value <= 0 || q.getLineWidth(line) <= q.Width.Pt() {
			res = append(res, line)
			continue
		}

		// break lines
		var w float64
		var currLine string
		for _, word := range strings.Split(line, " ") {
			wordWidth := q.getLineWidth(word)
			if w+wordWidth+charSpacing+spaceWidth > q.Width.Pt() && currLine != "" {
				res = append(res, currLine)
				currLine = ""
				w = 0
			}

			if currLine != "" {
				currLine += " "
				w += charSpacing + spaceWidth
			}

			// break word?
			for wordWidth > q.Width.Pt() {
				var wordPart string
				var wordPartWidth float64
				for _, c := range word {
					charWidth := q.getLineWidth(string(c))
					if len(wordPart) > 0 && wordPartWidth+charWidth+charSpacing > q.Width.Pt()-0.1 {
						res = append(res, wordPart)
						word = word[len(wordPart):]
						wordWidth = q.getLineWidth(word)
						break
					}
					wordPart += string(c)
					wordPartWidth += charWidth + charSpacing
				}
			}

			currLine += word
			w += wordWidth + charSpacing
		}
		if currLine != "" {
			res = append(res, currLine)
		}
	}

	if maxLines > 0 && len(res) > maxLines {
		res = res[:maxLines]
	}
	return res
}

// TextHeight returns the height of the text, accounting for line breaks and max width
func (q *TextBoxElement) TextHeight() Length {
	lines := len(q.wrappedText())
	return Pt(float64(lines) * q.lineHeight())
}
