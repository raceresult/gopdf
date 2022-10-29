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
	top := float64(page.Data.MediaBox.URY) - q.Y.Pt() - q.Font.GetAscent(q.FontSize)
	var c float64
	if q.Italic {
		c = 0.333
	}
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
		left := q.X.Pt()
		width := q.Font.GetWidth(line, q.FontSize)
		switch q.TextAlign {
		case TextAlignCenter:
			left += (q.Width.Pt() - width) / 2
		case TextAlignRight:
			left += q.Width.Pt() - width
		}
		if q.Rotate != 0 {
			r := q.Rotate * math.Pi / 180
			page.TextPosition_Tm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), left, top)
		} else {
			page.TextPosition_Tm(1, 0, c, 1, left, top)
		}
		page.TextShowing_Tj(line)

		// underline text
		if q.Underline {
			page.Path_re(left, top+q.Font.GetUnderlinePosition(q.FontSize), width, q.Font.GetUnderlineThickness(q.FontSize))
			page.Path_f()
		}

		top -= lineHeight
	}
	page.TextObjects_ET()
	return nil
}

// wrappedText returns the wrapped text considering line break, max width and max height
func (q *TextBoxElement) wrappedText() []string {
	// determine max number of lines
	var maxLines int
	if q.Height.Value > 0 {
		maxLines = int(q.Height.Pt() / q.lineHeight())
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
			if wordWidth+w > q.Width.Pt() && currLine != "" {
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
					if len(wordPart) > 0 && wordPartWidth+charWidth > q.Width.Pt() {
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
	return res
}

// TextHeight returns the height of the text, accounting for line breaks and max width
func (q *TextBoxElement) TextHeight() Length {
	lines := len(q.wrappedText())
	return Pt(float64(lines) * q.lineHeight())
}

func (q *TextBoxElement) lineHeight() float64 {
	if q.LineHeight != 0 {
		return q.LineHeight
	}
	return q.FontSize * 1.2
}
