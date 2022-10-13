package builder

import (
	"git.rrdc.de/lib/gopdf/pdf"
)

// ImageElement is used to add an image to a page
type ImageElement struct {
	Img    *pdf.Image
	Width  Length
	Height Length
	Left   Length
	Top    Length
}

// Build adds the element to the content stream
func (q *ImageElement) Build(page *pdf.Page) {
	if q.Img == nil {
		return
	}
	page.GraphicsState_q()
	page.GraphicsState_cm(q.Width.Pt(), 0, 0, q.Height.Pt(), q.Left.Pt(), float64(page.Data.MediaBox.URY)-q.Top.Pt()-q.Height.Pt())
	page.XObject_Do(q.Img.Reference)
	page.GraphicsState_Q()
}
