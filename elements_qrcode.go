package gopdf

import (
	"math"

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
	Rotate          float64
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
	qr.DisableBorder = true
	bits := qr.Bitmap()

	// graphics state
	page.GraphicsState_q()
	defer page.GraphicsState_Q()

	// set coordinate system
	y := float64(page.Data.MediaBox.URY) - q.Top.Pt()
	if q.Rotate == 0 {
		page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), y)
	} else {
		r := q.Rotate * math.Pi / 180
		page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), q.Left.Pt(), y)
	}

	// set position
	//page.GraphicsState_cm(1, 0, 0, 1, 0, -q.Size.Pt())

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
				page.Path_re(float64(x)*bitSize, -float64(y+1)*bitSize, bitSize, bitSize)
			}
		}
	}
	page.Path_f()
	return nil
}
