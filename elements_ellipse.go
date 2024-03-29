package gopdf

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// EllipseElement is used to add an ellipse to a page
type EllipseElement struct {
	Left, Top, Width, Height Length
	LineWidth                Length
	LineColor                Color
	FillColor                Color
	DashPattern              DashPattern
	Transparency             float64
}

// Build adds the element to the content stream
func (q *EllipseElement) Build(page *pdf.Page) (string, error) {
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

	// set dash pattern
	if warning, err := q.DashPattern.Build(page); err != nil {
		return warning, err
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

	// draw
	x := q.Left.Pt() + q.Width.Pt()/2
	y := float64(page.Data.MediaBox.URY) - q.Top.Pt() - q.Height.Pt()/2
	page.Path_m(x, y+q.Height.Pt()/2)
	page.Path_c(x+.5523*q.Width.Pt()/2, y+q.Height.Pt()/2, x+q.Width.Pt()/2, y+.5523*q.Height.Pt()/2, x+q.Width.Pt()/2, y)
	page.Path_c(x+q.Width.Pt()/2, y-.5523*q.Height.Pt()/2, x+.5523*q.Width.Pt()/2, y-q.Height.Pt()/2, x, y-q.Height.Pt()/2)
	page.Path_c(x-.5523*q.Width.Pt()/2, y-q.Height.Pt()/2, x-q.Width.Pt()/2, y-.5523*q.Height.Pt()/2, x-q.Width.Pt()/2, y)
	page.Path_c(x-q.Width.Pt()/2, y+.5523*q.Height.Pt()/2, x-.5523*q.Width.Pt()/2, y+q.Height.Pt()/2, x, y+q.Height.Pt()/2)
	if q.LineColor != nil && q.FillColor != nil {
		page.Path_B()
	} else if q.LineColor != nil {
		page.Path_S()
	} else {
		page.Path_f()
	}
	return "", nil
}
