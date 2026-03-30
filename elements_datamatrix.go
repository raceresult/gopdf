package gopdf

import (
	"image"
	"math"

	"github.com/boombuler/barcode/datamatrix"
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

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
	cols := bounds.Dx()
	rows := bounds.Dy()

	y := float64(page.Data.MediaBox.URY) - top.Pt()
	if rotate == 0 {
		page.GraphicsState_cm(1, 0, 0, 1, left.Pt(), y)
	} else {
		// Translate to the centre of the symbol on the page
		cx := left.Pt() + size.Pt()/2
		cy := y - size.Pt()/2
		rot := rotate * math.Pi / 180
		page.GraphicsState_cm(math.Cos(rot), math.Sin(rot), -math.Sin(rot), math.Cos(rot), cx, cy)
		// Shift origin back to top-left so module drawing coordinates are unchanged
		page.GraphicsState_cm(1, 0, 0, 1, -size.Pt()/2, size.Pt()/2)
	}

	// transparency
	if transparency > 0 && transparency <= 1 {
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - transparency),
			"CA": types.Number(1 - transparency),
		})
		page.GraphicsState_gs(n)
	}

	// draw modules — the boombuler library only generates square data-matrix codes (even if rectangular is allowed in the spec)
	moduleSize := math.Min(size.Pt()/float64(cols), size.Pt()/float64(rows))

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			// The image is not guaranteed to start at (0,0) so we add the start offset.
			r, g, b, _ := bc.At(bounds.Min.X+col, bounds.Min.Y+row).RGBA()
			if r == 0 && g == 0 && b == 0 {
				page.Path_re(
					float64(col)*moduleSize,
					-float64(row+1)*moduleSize,
					moduleSize,
					moduleSize,
				)
			}
		}
	}
	page.Path_f()

	return nil
}
