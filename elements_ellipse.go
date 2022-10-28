package gopdf

import "github.com/raceresult/gopdf/pdf"

// EllipseElement is used to add an ellipse to a page
type EllipseElement struct {
	X, Y, Width, Height Length

	LineWidth   Length
	LineColor   Color
	FillColor   Color
	DashPattern DashPattern
}

// Build adds the element to the content stream
func (q *EllipseElement) Build(page *pdf.Page) {
	if q.LineColor == nil {
		page.GraphicsState_w(0)
	} else {
		q.LineColor.Build(page, true)
		page.GraphicsState_w(q.LineWidth.Pt())
	}
	if q.FillColor != nil {
		q.FillColor.Build(page, false)
	}
	q.DashPattern.Build(page)

	y := float64(page.Data.MediaBox.URY) - q.Y.Pt() - q.Height.Pt()/2
	page.Path_m(q.X.Pt(), y+q.Height.Pt()/2)
	page.Path_c(q.X.Pt()+.5523*q.Width.Pt()/2, y+q.Height.Pt()/2, q.X.Pt()+q.Width.Pt()/2, y+.5523*q.Height.Pt()/2, q.X.Pt()+q.Width.Pt()/2, y)
	page.Path_c(q.X.Pt()+q.Width.Pt()/2, y-.5523*q.Height.Pt()/2, q.X.Pt()+.5523*q.Width.Pt()/2, y-q.Height.Pt()/2, q.X.Pt(), y-q.Height.Pt()/2)
	page.Path_c(q.X.Pt()-.5523*q.Width.Pt()/2, y-q.Height.Pt()/2, q.X.Pt()-q.Width.Pt()/2, y-.5523*q.Height.Pt()/2, q.X.Pt()-q.Width.Pt()/2, y)
	page.Path_c(q.X.Pt()-q.Width.Pt()/2, y+.5523*q.Height.Pt()/2, q.X.Pt()-.5523*q.Width.Pt()/2, y+q.Height.Pt()/2, q.X.Pt(), y+q.Height.Pt()/2)

	if q.LineColor != nil && q.FillColor != nil {
		page.Path_B()
	} else if q.LineColor != nil {
		page.Path_S()
	} else {
		page.Path_f()
	}
}
