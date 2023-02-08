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
	return Pt(q.Font.GetHeight(q.FontSize))
}

func (q *TextChunk) setFontAndColor(page *pdf.Page) error {
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
		case color != nil && q.OutlineColor != nil:
			color.Build(page, false)
			q.OutlineColor.Build(page, true)
			page.TextState_Tr(types.RenderingModeFillAndStroke)
		case color != nil:
			color.Build(page, false)
			page.TextState_Tr(types.RenderingModeFill)
		case q.OutlineColor != nil:
			page.TextState_Tr(types.RenderingModeStroke)
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
	return nil
}

func (q *TextChunk) draw(page *pdf.Page, left, top float64) error {
	// if no font or text given, ignore chunk
	if q.Font == nil {
		return errors.New("no font set")
	}

	// set format
	if err := q.setFontAndColor(page); err != nil {
		return err
	}

	// begin text and set position
	page.TextObjects_BT()
	var c float64
	if q.Italic {
		c = 0.333
	}
	page.TextPosition_Tm(1, 0, c, 1, left, top)

	// draw text
	page.TextShowing_Tj(reverseRTLString(q.Text))
	page.TextObjects_ET()

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

	return nil
}
