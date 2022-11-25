package gopdf

import (
	"errors"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// ImageElement is used to add an image to a page
type ImageElement struct {
	Img                      *pdf.Image
	Left, Top, Width, Height Length
	Transparency             float64
}

// Build adds the element to the content stream
func (q *ImageElement) Build(page *pdf.Page) error {
	// abort if image not set
	if q.Img == nil {
		return errors.New("image not set")
	}

	// graphics state
	page.GraphicsState_q()
	defer page.GraphicsState_Q()

	// transparency
	if q.Transparency > 0 && q.Transparency <= 1 {
		page.GraphicsState_q()
		defer page.GraphicsState_Q()
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - q.Transparency),
			"CA": types.Number(1 - q.Transparency),
		})
		page.GraphicsState_gs(n)
	}

	// draw
	page.GraphicsState_cm(q.Width.Pt(), 0, 0, q.Height.Pt(), q.Left.Pt(), float64(page.Data.MediaBox.URY)-q.Top.Pt()-q.Height.Pt())
	page.XObject_Do(q.Img.Reference)
	return nil
}
