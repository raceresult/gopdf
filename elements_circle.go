package gopdf

import "github.com/raceresult/gopdf/pdf"

// CircleElement is used to add a circle to a page
type CircleElement struct {
	X, Y, Radius Length
	LineWidth    Length
	LineColor    Color
	FillColor    Color
	DashPattern  DashPattern
}

// Build adds the element to the content stream
func (q *CircleElement) Build(page *pdf.Page) error {
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

	y := float64(page.Data.MediaBox.URY) - q.Y.Pt()
	page.Path_m(q.X.Pt(), y+q.Radius.Pt())
	page.Path_c(q.X.Pt()+.5523*q.Radius.Pt(), y+q.Radius.Pt(), q.X.Pt()+q.Radius.Pt(), y+.5523*q.Radius.Pt(), q.X.Pt()+q.Radius.Pt(), y)
	page.Path_c(q.X.Pt()+q.Radius.Pt(), y-.5523*q.Radius.Pt(), q.X.Pt()+.5523*q.Radius.Pt(), y-q.Radius.Pt(), q.X.Pt(), y-q.Radius.Pt())
	page.Path_c(q.X.Pt()-.5523*q.Radius.Pt(), y-q.Radius.Pt(), q.X.Pt()-q.Radius.Pt(), y-.5523*q.Radius.Pt(), q.X.Pt()-q.Radius.Pt(), y)
	page.Path_c(q.X.Pt()-q.Radius.Pt(), y+.5523*q.Radius.Pt(), q.X.Pt()-.5523*q.Radius.Pt(), y+q.Radius.Pt(), q.X.Pt(), y+q.Radius.Pt())

	if q.LineColor != nil && q.FillColor != nil {
		page.Path_B()
	} else if q.LineColor != nil {
		page.Path_S()
	} else {
		page.Path_f()
	}
	return nil
}
