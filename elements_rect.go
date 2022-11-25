package gopdf

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// RectElement is used to add a rectangle to a page
type RectElement struct {
	Left, Top, Width, Height Length
	LineWidth                Length
	LineColor                Color
	FillColor                Color
	DashPattern              DashPattern
	Transparency             float64
}

// Build adds the element to the content stream
func (q *RectElement) Build(page *pdf.Page) error {
	// set colors
	if q.LineColor == nil {
		page.GraphicsState_w(0)
	} else {
		q.LineColor.Build(page, true)
		page.GraphicsState_w(q.LineWidth.Pt())
	}
	if q.FillColor != nil {
		q.FillColor.Build(page, false)
	}
	if err := q.DashPattern.Build(page); err != nil {
		return err
	}

	// transparency
	if q.Transparency > 0 && q.Transparency <= 1 {
		page.GraphicsState_q()
		defer page.GraphicsState_Q()
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - q.Transparency),
			"CA": types.Number(1 - q.Transparency),
		})
		page.GraphicsState_gs(n)
	}

	// create rect
	page.Path_re(q.Left.Pt(), float64(page.Data.MediaBox.URY)-q.Top.Pt()-q.Height.Pt(), q.Width.Pt(), q.Height.Pt())

	// fill or stroke
	if q.LineColor != nil && q.FillColor != nil {
		page.Path_B()
	} else if q.LineColor != nil {
		page.Path_S()
	} else {
		page.Path_f()
	}
	return nil
}
