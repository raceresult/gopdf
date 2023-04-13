package gopdf

import (
	"math"

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
	Rotate                   float64
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

	page.GraphicsState_q()
	defer page.GraphicsState_Q()

	// transparency
	if q.Transparency > 0 && q.Transparency <= 1 {
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - q.Transparency),
			"CA": types.Number(1 - q.Transparency),
		})
		page.GraphicsState_gs(n)
	}

	// set coordinate system
	y := float64(page.Data.MediaBox.URY) - q.Top.Pt()
	if q.Rotate == 0 {
		page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), y)
	} else {
		r := q.Rotate * math.Pi / 180
		page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), q.Left.Pt(), y)
	}

	// create rect
	page.Path_re(0, -q.Height.Pt(), q.Width.Pt(), q.Height.Pt())

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
