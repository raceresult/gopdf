package gopdf

import (
	"errors"
	"math"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// ImageBoxElement is used to add an image box to a page
type ImageBoxElement struct {
	Left, Top, Width, Height Length
	Img                      *pdf.Image
	VerticalAlign            VerticalAlign
	HorizontalAlign          HorizontalAlign
	Rotate                   float64
	Transparency             float64
}

// Build adds the element to the content stream
func (q *ImageBoxElement) Build(page *pdf.Page) (string, error) {
	// abort if image not set
	if q.Img == nil {
		return "", errors.New("image not set")
	}
	// ignore if size 0
	if q.Img.Image.Height == 0 || q.Img.Image.Width == 0 {
		return "", nil
	}

	// graphics state
	page.GraphicsState_q()
	defer page.GraphicsState_Q()

	// transparency
	if q.Transparency > 0 && q.Transparency <= 1 {
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - q.Transparency),
			"CA": types.Number(1 - q.Transparency),
		})
		page.GraphicsState_gs(n)
	}

	if q.Width.Pt()/q.Height.Pt() > float64(q.Img.Image.Width)/float64(q.Img.Image.Height) {
		w := q.Height.Pt() * float64(q.Img.Image.Width) / float64(q.Img.Image.Height)
		var left float64
		switch q.HorizontalAlign {
		case HorizontalAlignCenter:
			left += (q.Width.Pt() - w) / 2
		case HorizontalAlignRight:
			left += q.Width.Pt() - w
		}

		// set coordinate system
		y := float64(page.Data.MediaBox.URY) - q.Top.Pt()
		if q.Rotate == 0 {
			page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), y)
		} else {
			r := q.Rotate * math.Pi / 180
			page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), q.Left.Pt(), y)
		}

		page.GraphicsState_cm(w, 0, 0, q.Height.Pt(), left, -q.Height.Pt())

	} else {
		h := q.Width.Pt() * float64(q.Img.Image.Height) / float64(q.Img.Image.Width)
		y := float64(page.Data.MediaBox.URY) - q.Top.Pt()
		top := -h
		switch q.VerticalAlign {
		case VerticalAlignMiddle:
			top -= (q.Height.Pt() - h) / 2
		case VerticalAlignBottom:
			top -= q.Height.Pt() - h
		default:
		}

		// set coordinate system
		if q.Rotate == 0 {
			page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), y)
		} else {
			r := q.Rotate * math.Pi / 180
			page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), q.Left.Pt(), y)
		}

		page.GraphicsState_cm(q.Width.Pt(), 0, 0, h, 0, top)
	}

	page.XObject_Do(q.Img.Reference)
	return "", nil
}
