package gopdf

import (
	"math"
	"strings"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// TextElement draws a text, may have line breaks
type TextElement struct {
	Text          string
	Left, Top     Length
	Font          pdf.FontHandler
	FontSize      float64
	Color         Color
	RenderMode    types.RenderingMode
	OutlineColor  Color
	OutlineWidth  Length
	TextAlign     HorizontalAlign
	LineHeight    float64
	Bold          bool
	Italic        bool
	Underline     bool
	StrikeThrough bool
	CharSpacing   Length
	TextScaling   float64
	Rotate        float64
}

// Build adds the element to the content stream
func (q *TextElement) Build(page *pdf.Page) error {
	// if no font or text given, ignore element
	if q.Font == nil || q.Text == "" {
		return nil
	}

	// set coordinate system
	page.GraphicsState_q()
	y := float64(page.Data.MediaBox.URY) - q.Top.Pt()
	if q.Rotate == 0 {
		page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), y)
	} else {
		r := q.Rotate * math.Pi / 180
		page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), q.Left.Pt(), y)
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

	// begin text and set font
	page.TextObjects_BT()
	page.TextState_Tf(q.Font, q.FontSize)

	// calculate some values needed below
	lineHeight := q.lineHeight()
	top := 0.0
	var c float64
	if q.Italic {
		c = 0.333
	}

	// iterate over lines
	for _, line := range strings.Split(q.Text, "\n") {
		width := q.getLineWidth(line)
		left := 0.0
		switch q.TextAlign {
		case HorizontalAlignCenter:
			left -= width / 2
		case HorizontalAlignRight:
			left -= width
		}

		page.TextPosition_Tm(1, 0, c, 1, left, top)
		page.TextShowing_Tj(line)

		// underline/strike-through text
		if q.Underline {
			page.Path_re(
				left, top+q.Font.GetUnderlinePosition(q.FontSize),
				width, q.Font.GetUnderlineThickness(q.FontSize),
			)
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

// TextHeight returns the height of the text, accounting for line breaks
func (q *TextElement) TextHeight() Length {
	lines := strings.Count(q.Text, "\n") + 1
	return Pt(float64(lines) * q.lineHeight())
}

// FontHeight returns the height of the font (bounding box y min to max)
func (q *TextElement) FontHeight() Length {
	return Pt(q.Font.GetTop(q.FontSize) - q.Font.GetBottom(q.FontSize))
}

// getLineWidth returns the width of the given text line consider font, fontsize, text-scaling, char spacing
func (q *TextElement) getLineWidth(line string) float64 {
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

func (q *TextElement) lineHeight() float64 {
	if q.LineHeight != 0 {
		return q.LineHeight
	}
	return q.FontHeight().Pt()
}
