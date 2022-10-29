package gopdf

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/skip2/go-qrcode"
)

// QRCodeElement is used to add a qrcode to a page
type QRCodeElement struct {
	Left, Top, Size Length
	Text            string
	Color           Color
	RecoveryLevel   qrcode.RecoveryLevel
}

// Build adds the element to the content stream
func (q *QRCodeElement) Build(page *pdf.Page) error {
	if q.Color != nil {
		q.Color.Build(page, false)
	} else {
		ColorRGBBlack.Build(page, false)
	}

	page.GraphicsState_q()
	page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), float64(page.Data.MediaBox.URY)-q.Top.Pt())

	// create qr code
	qr, err := qrcode.New(q.Text, q.RecoveryLevel)
	if err != nil {
		return err
	}
	bits := qr.Bitmap()

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

	page.GraphicsState_Q()
	return nil
}
