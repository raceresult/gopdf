package gopdf

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// Form is used to add a captured page from another document to a page
type Form struct {
	BBox                     types.Rectangle
	Form                     types.Reference
	Left, Top, Width, Height Length
}

// Build adds the element to the content stream
func (q *Form) Build(page *pdf.Page) error {
	h := q.Height.Pt()
	if h == 0 {
		h = float64(q.BBox.URY)
	}
	w := q.Width.Pt()
	if w == 0 {
		w = float64(q.BBox.URX)
	}

	offsetY := float64(page.Data.MediaBox.URY) - float64(q.BBox.URY)*h/float64(q.BBox.URY) - q.Top.Pt()
	if q.Left.Value != 0 || offsetY != 0 || q.Width.Value != 0 || q.Height.Value != 0 {
		page.GraphicsState_q()
		page.GraphicsState_cm(w/float64(q.BBox.URX), 0, 0, h/float64(q.BBox.URY), q.Left.Pt(), offsetY)
	}
	page.XObject_Do(q.Form)
	if q.Left.Value != 0 || offsetY != 0 || q.Width.Value != 0 || q.Height.Value != 0 {
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
