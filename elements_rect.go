package gopdf

import "github.com/raceresult/gopdf/pdf"

// RectElement is used to add a rectangle to a page
type RectElement struct {
	X1, Y1, Width, Height Length
	LineWidth             Length
	LineColor             Color
	FillColor             Color
}

// Build adds the element to the content stream
func (q *RectElement) Build(page *pdf.Page) {
	if q.LineColor == nil {
		page.GraphicsState_w(0)
	} else {
		q.LineColor.Build(page, true)
		page.GraphicsState_w(q.LineWidth.Pt())
	}
	if q.FillColor != nil {
		q.FillColor.Build(page, false)
	}

	page.Path_re(q.X1.Pt(), float64(page.Data.MediaBox.URY)-q.Y1.Pt()-q.Height.Pt(), q.Width.Pt(), q.Height.Pt())

	if q.LineColor != nil && q.FillColor != nil {
		page.Path_B()
	} else if q.LineColor != nil {
		page.Path_S()
	} else {
		page.Path_f()
	}
}
