package gopdf

import "github.com/raceresult/gopdf/pdf"

// CapturedPage is used to add a captured page from another document to a page
type CapturedPage struct {
	CapturedPage *pdf.CapturedPage
	Left         Length
	Top          Length
}

// Build adds the element to the content stream
func (q *CapturedPage) Build(page *pdf.Page) {
	if q.CapturedPage == nil {
		return
	}
	offsetY := float64(page.Data.MediaBox.URY) - float64(q.CapturedPage.Source.MediaBox.URY) - q.Top.Pt()
	if q.Left.Value != 0 || offsetY != 0 {
		page.GraphicsState_q()
		page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), offsetY)
	}
	page.AddCapturedPage(q.CapturedPage)
	if q.Left.Value != 0 || offsetY != 0 {
		page.GraphicsState_Q()
	}
}

// PageSize returns the page size of the captured page
func (q *CapturedPage) PageSize() PageSize {
	return PageSize{
		Pt(float64(q.CapturedPage.Source.MediaBox.URX)),
		Pt(float64(q.CapturedPage.Source.MediaBox.URY)),
	}
}
