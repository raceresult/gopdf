package gopdf

import (
	"image"
	"math"

	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/qr"
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// QRCodeElement is used to add a QRCode to a page
type QRCodeElement struct {
	Left, Top, Size Length
	Text            string
	Color           Color
	RecoveryLevel   qr.ErrorCorrectionLevel
	Rotate          float64
	Transparency    float64
}

// Build adds the element to the content stream
// The caller needs to ensure 4 modules of quiet-zone on all sides of the symbol.
// A previous version of this library used another QR-Code-Library: "github.com/skip2/go-qrcode"
// The Error-Correction value-mappings are as follows:
//
//	qrcode.Low -> qr.L 		(7%)
//	qrcode.Medium -> qr.M 	(15%)
//	qrcode.High -> qr.Q 	(25%)
//	qrcode.Highest -> qr.H 	(30%)
func (q *QRCodeElement) Build(page *pdf.Page) (string, error) {
	bc, err := qr.Encode(q.Text, q.RecoveryLevel, qr.Auto)
	if err != nil {
		return "", err
	}

	page.GraphicsState_q()
	defer page.GraphicsState_Q()

	err = build2DCode(
		page,
		q.Left, q.Top, q.Size,
		q.Rotate, q.Transparency, q.Color,
		bc,
	)

	return "", err
}

// DataMatrixElement is used to add a Data Matrix code to a page
// Error correction is always ~25% for Data Matrix codes.
// The caller is responsible for leaving a one-module quiet zone on all sides of the symbol.
type DataMatrixElement struct {
	Left, Top, Size Length
	Text            string
	Color           Color
	Rotate          float64
	Transparency    float64
}

func (q *DataMatrixElement) Build(page *pdf.Page) (string, error) {
	bc, err := datamatrix.Encode(q.Text)
	if err != nil {
		return "", err
	}

	page.GraphicsState_q()
	defer page.GraphicsState_Q()

	err = build2DCode(
		page,
		q.Left, q.Top, q.Size,
		q.Rotate, q.Transparency, q.Color,
		bc,
	)

	return "", err
}

// build2DCode handles coordinate transform, transparency, and module drawing
// for any 2D matrix code. isSet reports whether the module at (col, row) is filled.
func build2DCode(
	page *pdf.Page,
	left, top, size Length,
	rotate, transparency float64,
	clr Color,
	bc image.Image,
) error {
	// set color
	if clr != nil {
		clr.Build(page, false)
	} else {
		ColorRGBBlack.Build(page, false)
	}

	// coordinate system
	bounds := bc.Bounds()
	cols := bounds.Max.X
	rows := bounds.Max.Y

	y := float64(page.Data.MediaBox.URY) - top.Pt()
	if rotate == 0 {
		page.GraphicsState_cm(1, 0, 0, 1, left.Pt(), y)
	} else {
		r := rotate * math.Pi / 180
		page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), left.Pt(), y)
	}

	// transparency
	if transparency > 0 && transparency <= 1 {
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - transparency),
			"CA": types.Number(1 - transparency),
		})
		page.GraphicsState_gs(n)
	}

	// draw modules — x/y step sizes are independent to support rectangular symbols
	xBitSize := size.Pt() / float64(cols)
	yBitSize := size.Pt() / float64(rows)

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			r, g, b, _ := bc.At(col, row).RGBA()
			if r == 0 && g == 0 && b == 0 {
				page.Path_re(
					float64(col)*xBitSize,
					-float64(row+1)*yBitSize,
					xBitSize,
					yBitSize,
				)
			}
		}
	}
	page.Path_f()

	return nil
}
