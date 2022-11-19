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
	RenderMode    types.RenderingMode // todo: remove - can be done by outlinecolor+color
	OutlineColor  Color
	OutlineWidth  Length
	DashPattern   DashPattern
	TextAlign     HorizontalAlign
	LineHeight    float64
	Bold          bool
	Italic        bool
	Underline     bool
	StrikeThrough bool
	CharSpacing   Length
	TextScaling   float64
	Rotate        float64
	Transparency  float64
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
	if err := q.DashPattern.Build(page); err != nil {
		return err
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

	// Transparency
	if q.Transparency > 0 && q.Transparency <= 1 {
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - q.Transparency),
			"CA": types.Number(1 - q.Transparency),
		})
		page.GraphicsState_gs(n)
	}

	// begin text
	page.TextObjects_BT()

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
		page.TextShowing_Tj(reverseRTLString(line))

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

// reverseRTLString reverse the parts of the given string which are part of a right-to-left language
// only if the first character of the string is right-to-left
func reverseRTLString(s string) string {
	if s == "" {
		return s
	}

	isRTL := func(r rune) bool {
		return r >= 0x591 && r <= 0x6EF
	}

	rr := []rune(s)
	if !isRTL(rr[0]) {
		return s
	}

	b := true
	l := 0
	arr := make([]rune, 0, len(rr))
	for i, r := range rr {
		rtl := isRTL(r)
		if rtl && !b {
			arr = append(arr, rr[l:i]...)
			l = i
			b = !b
		} else if r >= 0x30 && r <= 0x39 && b || r >= 0x41 && r <= 0x5A && b || r >= 0x61 && r <= 0x7A && b {
			for x := i - 1; x >= l; x-- {
				arr = append(arr, rr[x])
			}
			l = i
			b = !b
		} else if i == len(rr)-1 {
			if rtl {
				for x := len(rr) - 1; x >= l; x-- {
					arr = append(arr, rr[x])
				}
			} else {
				arr = append(arr, rr[l:]...)
			}
		}
	}
	return string(arr)
}
