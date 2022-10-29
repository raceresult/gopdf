package gopdf

import (
	"math"
	"strings"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// TextElement draws a text, may have line breaks
type TextElement struct {
	Text         string
	X, Y         Length
	Font         pdf.FontHandler
	FontSize     float64
	Color        Color
	RenderMode   types.RenderingMode
	OutlineColor Color
	OutlineWidth Length
	TextAlign    TextAlign
	LineHeight   float64
	Italic       bool
	Bold         bool
	Underline    bool
	Rotate       float64
}

// Build adds the element to the content stream
func (q *TextElement) Build(page *pdf.Page) error {
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

	// calculate some values needed below
	lineHeight := q.lineHeight()
	top := float64(page.Data.MediaBox.URY) - q.Y.Pt()
	var c float64
	if q.Italic {
		c = 0.333
	}

	// iterate over lines
	for _, line := range strings.Split(q.Text, "\n") {
		width := q.Font.GetWidth(line, q.FontSize)
		left := q.X.Pt()
		switch q.TextAlign {
		case TextAlignCenter:
			left -= width / 2
		case TextAlignRight:
			left -= width
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
			page.Path_re(
				left, top+q.Font.GetUnderlinePosition(q.FontSize),
				width, q.Font.GetUnderlineThickness(q.FontSize),
			)
			page.Path_f()
		}

		top += lineHeight
	}
	page.TextObjects_ET()
	return nil
}

// TextHeight returns the height of the text, accounting for line breaks
func (q *TextElement) TextHeight() Length {
	lines := strings.Count(q.Text, "\n") + 1
	return Pt(float64(lines) * q.lineHeight())
}

func (q *TextElement) lineHeight() float64 {
	if q.LineHeight != 0 {
		return q.LineHeight
	}
	return q.FontSize * 1.2
}
