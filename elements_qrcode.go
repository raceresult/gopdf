package gopdf

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
	"github.com/skip2/go-qrcode"
)

// QRCodeElement is used to add a QRCode to a page
type QRCodeElement struct {
	Left, Top, Size Length
	Text            string
	Color           Color
	RecoveryLevel   qrcode.RecoveryLevel
	Transparency    float64
}

// Build adds the element to the content stream
func (q *QRCodeElement) Build(page *pdf.Page) error {
	// set color
	if q.Color != nil {
		q.Color.Build(page, false)
	} else {
		ColorRGBBlack.Build(page, false)
	}

	// create qr code bitmap
	qr, err := qrcode.New(q.Text, q.RecoveryLevel)
	if err != nil {
		return err
	}
	bits := qr.Bitmap()

	// graphics state
	page.GraphicsState_q()
	defer page.GraphicsState_Q()

	// set position
	page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), float64(page.Data.MediaBox.URY)-q.Top.Pt())

	// transparency
	if q.Transparency > 0 && q.Transparency <= 1 {
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - q.Transparency),
			"CA": types.Number(1 - q.Transparency),
		})
		page.GraphicsState_gs(n)
	}

	// draw
	bitSize := q.Size.Pt() / float64(len(bits))
	for x, row := range bits {
		for y, value := range row {
			if value {
				page.Path_re(float64(x)*bitSize, -float64(y)*bitSize, bitSize, bitSize)
			}
		}
	}
	page.Path_f()
	return nil
}
