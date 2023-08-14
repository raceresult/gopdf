package gopdf

import (
	"errors"
	"math"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// ImageElement is used to add an image to a page
type ImageElement struct {
	Img                      *pdf.Image
	Left, Top, Width, Height Length
	Rotate                   float64
	Transparency             float64
}

// Build adds the element to the content stream
func (q *ImageElement) Build(page *pdf.Page) (string, error) {
	// abort if image not set
	if q.Img == nil {
		return "", errors.New("image not set")
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

	// set coordinate system
	y := float64(page.Data.MediaBox.URY) - q.Top.Pt()
	if q.Rotate == 0 {
		page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), y)
	} else {
		r := q.Rotate * math.Pi / 180
		page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), q.Left.Pt(), y)
	}

	// draw
	page.GraphicsState_cm(q.Width.Pt(), 0, 0, q.Height.Pt(), 0, -q.Height.Pt())
	page.XObject_Do(q.Img.Reference)
	return "", nil
}
