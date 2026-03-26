package gopdf

import (
	"image"
	"math"

	"github.com/boombuler/barcode/datamatrix"
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
func (q *QRCodeElement) Build(page *pdf.Page) (string, error) {
	// set color
	// This setting will leak to later elements in the PDF. It should probably be moved after page.GraphicsState_q() is called.
	if q.Color != nil {
		q.Color.Build(page, false)
	} else {
		ColorRGBBlack.Build(page, false)
	}

	// create qr code bitmap
	qr, err := qrcode.New(q.Text, q.RecoveryLevel)
	if err != nil {
		return "", err
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
	// page.GraphicsState_cm(1, 0, 0, 1, 0, -q.Size.Pt())

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
	return "", nil
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
// for any 2D matrix code passed as image.
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
		rot := rotate * math.Pi / 180
		page.GraphicsState_cm(math.Cos(rot), math.Sin(rot), -math.Sin(rot), math.Cos(rot), left.Pt(), y)
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
	moduleSize := math.Min(size.Pt()/float64(cols), size.Pt()/float64(rows))
	xBitSize := moduleSize
	yBitSize := moduleSize

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
