package gopdf

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// ImageBoxElement is used to add an image box to a page
type ImageBoxElement struct {
	Left, Top, Width, Height Length
	Img                      *pdf.Image
	VerticalAlign            VerticalAlign
	HorizontalAlign          HorizontalAlign
	Transparency             float64
}

// Build adds the element to the content stream
func (q *ImageBoxElement) Build(page *pdf.Page) error {
	if q.Img == nil || q.Img.Image.Height == 0 || q.Img.Image.Width == 0 {
		return nil
	}
	page.GraphicsState_q()

	// Transparency
	if q.Transparency > 0 && q.Transparency <= 1 {
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - q.Transparency),
			"CA": types.Number(1 - q.Transparency),
		})
		page.GraphicsState_gs(n)
	}

	if q.Width.Pt()/q.Height.Pt() > float64(q.Img.Image.Width)/float64(q.Img.Image.Height) {
		w := q.Height.Pt() * float64(q.Img.Image.Width) / float64(q.Img.Image.Height)
		left := q.Left.Pt()
		switch q.HorizontalAlign {
		case HorizontalAlignCenter:
			left += (q.Width.Pt() - w) / 2
		case HorizontalAlignRight:
			left += q.Width.Pt() - w
		}
		page.GraphicsState_cm(w, 0, 0, q.Height.Pt(), left, float64(page.Data.MediaBox.URY)-q.Top.Pt()-q.Height.Pt())

	} else {
		h := q.Width.Pt() * float64(q.Img.Image.Height) / float64(q.Img.Image.Width)
		top := float64(page.Data.MediaBox.URY) - q.Top.Pt() - h
		switch q.VerticalAlign {
		case VerticalAlignMiddle:
			top -= (q.Height.Pt() - h) / 2
		case VerticalAlignBottom:
			top -= q.Height.Pt() - h
		default:
		}
		page.GraphicsState_cm(q.Width.Pt(), 0, 0, h, q.Left.Pt(), top)
	}

	page.XObject_Do(q.Img.Reference)
	page.GraphicsState_Q()
	return nil
}
