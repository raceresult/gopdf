package builder

import "git.rrdc.de/lib/gopdf/pdf"

// LineElement is used to add a line to a page
type LineElement struct {
	X1, Y1, X2, Y2 Length
	LineWidth      Length
	Color          Color
}

// Build adds the element to the content stream
func (q *LineElement) Build(page *pdf.Page) {
	if q.Color == nil {
		ColorRGBBlack.Build(page, true)
	} else {
		q.Color.Build(page, true)
	}
	page.GraphicsState_w(q.LineWidth.Pt())
	page.Path_m(q.X1.Pt(), float64(page.Data.MediaBox.URY)-q.Y1.Pt())
	page.Path_l(q.X2.Pt(), float64(page.Data.MediaBox.URY)-q.Y2.Pt())
	page.Path_S()
}
