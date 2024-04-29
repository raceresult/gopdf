package gopdf

import (
	"errors"

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

// getLineWidth returns the width of the given text line consider font, fontsize, text-scaling, char spacing
func (q *TextChunk) getLineWidth(line string) float64 {
	if line == "" {
		return 0
	}

	var v float64
	if fb := q.Font.FallbackFont(); fb != nil {
		rr := []rune(line)
		hasGlyph := q.Font.HasGylph(rr)
		curr := hasGlyph[0]
		start := 0
		for i := range hasGlyph {
			if i < len(hasGlyph)-1 && hasGlyph[i+1] == curr {
				continue
			}

			var w float64
			part := string(rr[start : i+1])
			if !curr {
				w = fb.GetWidth(part, q.FontSize)
			} else {
				w = q.Font.GetWidth(part, q.FontSize)
			}
			v += w
			curr = !curr
			start = i + 1
		}

	} else {
		v = q.Font.GetWidth(line, q.FontSize)
	}

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
	return Pt(q.Font.GetHeight(q.FontSize))
}

func (q *TextChunk) setFontAndColor(page *pdf.Page) (string, error) {
	// set color and rendering mode
	color := q.Color
	if color == nil && q.OutlineColor == nil {
		color = ColorRGBBlack
	}
	if q.Bold && q.OutlineWidth.Value == 0 && color != nil && q.OutlineColor == nil {
		color.Build(page, false)
		color.Build(page, true)
		page.GraphicsState_w(q.FontSize * 0.05)
		page.TextState_Tr(types.RenderingModeFillAndStroke)
	} else {
		page.GraphicsState_w(q.OutlineWidth.Pt())
		switch {
		case color != nil && q.OutlineColor != nil && q.OutlineWidth.Value != 0:
			color.Build(page, false)
			q.OutlineColor.Build(page, true)
			page.TextState_Tr(types.RenderingModeFillAndStroke)
		case color != nil:
			color.Build(page, false)
			page.TextState_Tr(types.RenderingModeFill)
		case q.OutlineColor != nil && q.OutlineWidth.Value != 0:
			page.TextState_Tr(types.RenderingModeStroke)
			q.OutlineColor.Build(page, true)
		}
	}
	warning, err := q.DashPattern.Build(page)
	if err != nil {
		return warning, err
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
	return warning, nil
}

func (q *TextChunk) draw(page *pdf.Page, left, top float64) (string, error) {
	// if no font or text given, ignore chunk
	if q.Font == nil {
		return "", errors.New("no font set")
	}

	// set format
	warning, err := q.setFontAndColor(page)
	if err != nil {
		return warning, err
	}

	// begin text and set position
	var c float64
	if q.Italic {
		c = 0.333
	}

	// draw text
	if fb := q.Font.FallbackFont(); fb != nil || false {
		rr := []rune(q.Text)
		hasGlyph := q.Font.HasGylph(rr)
		curr := hasGlyph[0]
		start := 0
		for i := range hasGlyph {
			if i < len(hasGlyph)-1 && hasGlyph[i+1] == curr {
				continue
			}

			var w float64
			part := reverseRTLString(string(rr[start : i+1]))
			if !curr {
				page.TextState_Tf(fb, q.FontSize)
				w = fb.GetWidth(part, q.FontSize)
			} else {
				if start != 0 {
					page.TextState_Tf(q.Font, q.FontSize)
				}
				w = q.Font.GetWidth(part, q.FontSize)
			}
			page.TextObjects_BT()
			page.TextPosition_Tm(1, 0, c, 1, left, top)
			page.TextShowing_Tj(part)
			page.TextObjects_ET()

			left += w
			curr = !curr
			start = i + 1
		}

	} else {
		page.TextObjects_BT()
		page.TextPosition_Tm(1, 0, c, 1, left, top)
		page.TextShowing_Tj(reverseRTLString(q.Text))
		page.TextObjects_ET()
	}

	// underline/strike-through text
	if q.Underline {
		th := q.Font.GetUnderlineThickness(q.FontSize)
		page.Path_re(left, top+q.Font.GetUnderlinePosition(q.FontSize)-th, q.getLineWidth(q.Text), th)
		page.Path_f()
	}
	if q.StrikeThrough {
		page.Path_re(
			left, top+q.Font.GetTop(q.FontSize)/3,
			q.getLineWidth(q.Text), q.Font.GetUnderlineThickness(q.FontSize),
		)
		page.Path_f()
	}

	return warning, nil
}
