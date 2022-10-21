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
	if q.Left.Value != 0 || q.Top.Value != 0 {
		page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), -q.Top.Pt())
	}
	page.AddCapturedPage(q.CapturedPage)
}
