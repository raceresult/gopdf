package gopdf

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// Form is used to add a captured page from another document to a page
type Form struct {
	BBox  types.Rectangle
	Form  types.Reference
	Left  Length
	Top   Length
	Scale float64
}

// Build adds the element to the content stream
func (q *Form) Build(page *pdf.Page) error {
	if q.Scale == 0 {
		q.Scale = 1
	}
	offsetY := float64(page.Data.MediaBox.URY) - float64(q.BBox.URY)*q.Scale - q.Top.Pt()
	if q.Left.Value != 0 || offsetY != 0 || q.Scale != 1 {
		page.GraphicsState_q()
		page.GraphicsState_cm(q.Scale, 0, 0, q.Scale, q.Left.Pt(), offsetY)
	}
	page.XObject_Do(q.Form)
	if q.Left.Value != 0 || offsetY != 0 || q.Scale != 1 {
		page.GraphicsState_Q()
	}
	return nil
}

// PageSize returns the page size of the captured page
func (q *Form) PageSize() PageSize {
	return PageSize{
		Pt(float64(q.BBox.URX)),
		Pt(float64(q.BBox.URY)),
	}
}
