package gopdf

import "github.com/raceresult/gopdf/pdf"

// ImageBoxElement is used to add an image box to a page
type ImageBoxElement struct {
	Img             *pdf.Image
	Width           Length
	Height          Length
	Left            Length
	Top             Length
	VerticalAlign   VerticalAlign
	HorizontalAlign TextAlign
}

// Build adds the element to the content stream
func (q *ImageBoxElement) Build(page *pdf.Page) {
	if q.Img == nil || q.Img.Image.Height == 0 || q.Img.Image.Width == 0 {
		return
	}
	page.GraphicsState_q()

	if q.Width.Pt()/q.Height.Pt() > float64(q.Img.Image.Width)/float64(q.Img.Image.Height) {
		w := q.Height.Pt() * float64(q.Img.Image.Width) / float64(q.Img.Image.Height)
		left := q.Left.Pt()
		switch q.HorizontalAlign {
		case TextAlignCenter:
			left += (q.Width.Pt() - w) / 2
		case TextAlignRight:
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
}
