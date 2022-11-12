package gopdf

import (
	"math"
	"strings"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

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

	// set color
	color := q.Color
	if color == nil {
		color = ColorRGBBlack
	}
	color.Build(page, false)
	color.Build(page, true)

	// text scaling / char spacing
	page.TextState_Tc(q.CharSpacing.Pt())
	if q.TextScaling == 0 {
		page.TextState_Tz(100)
	} else {
		page.TextState_Tz(q.TextScaling)
	}

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

	// begin text and set font
	page.TextObjects_BT()
	page.TextState_Tf(q.Font, q.FontSize)

	wrapped := q.wrappedText()
	lineHeight := q.lineHeight()
	var c float64
	if q.Italic {
		c = 0.333
	}

	page.GraphicsState_q()
	r := q.Rotate * math.Pi / 180
	page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), q.Left.Pt(), float64(page.Data.MediaBox.URY)-q.Top.Pt())

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
		width := q.Font.GetWidth(line, q.FontSize)
		switch q.TextAlign {
		case HorizontalAlignCenter:
			left += (q.Width.Pt() - width) / 2
		case HorizontalAlignRight:
			left += q.Width.Pt() - width
		}
		page.TextPosition_Tm(1, 0, c, 1, left, top)
		page.TextShowing_Tj(line)

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

	page.GraphicsState_Q()
	page.TextObjects_ET()
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
	spaceWidth := q.Font.GetWidth(" ", q.FontSize)
	for _, line := range lines {
		// check if max number of lines reached
		if maxLines > 0 && len(res) >= maxLines {
			break
		}

		// if width of line does not exceed max width, we can add the entire line and continue with the next line
		if q.Width.Value <= 0 || q.Font.GetWidth(line, q.FontSize) <= q.Width.Pt() {
			res = append(res, line)
			continue
		}

		// break lines
		var w float64
		var currLine string
		for _, word := range strings.Split(line, " ") {
			wordWidth := q.Font.GetWidth(word, q.FontSize)
			if wordWidth+w+spaceWidth > q.Width.Pt() && currLine != "" {
				res = append(res, currLine)
				currLine = ""
				w = 0
			}

			if currLine != "" {
				currLine += " "
				w += spaceWidth
			}

			// break word?
			for wordWidth > q.Width.Pt() {
				var wordPart string
				var wordPartWidth float64
				for _, c := range word {
					charWidth := q.Font.GetWidth(string(c), q.FontSize)
					if len(wordPart) > 0 && wordPartWidth+charWidth > q.Width.Pt()-0.1 {
						res = append(res, wordPart)
						word = word[len(wordPart):]
						wordWidth = q.Font.GetWidth(word, q.FontSize)
						break
					}
					wordPart += string(c)
					wordPartWidth += charWidth
				}
			}

			currLine += word
			w += wordWidth
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
