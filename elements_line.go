package gopdf

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// LineElement is used to add a line to a page
type LineElement struct {
	X1, Y1, X2, Y2 Length
	LineWidth      Length
	Color          Color
	DashPattern    DashPattern
	Transparency   float64
}

// Build adds the element to the content stream
func (q *LineElement) Build(page *pdf.Page) error {
	// color
	if q.Color == nil {
		ColorRGBBlack.Build(page, true)
	} else {
		q.Color.Build(page, true)
	}
	if err := q.DashPattern.Build(page); err != nil {
		return err
	}

	// line width
	page.GraphicsState_w(q.LineWidth.Pt())

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

	// draw
	page.Path_m(q.X1.Pt(), float64(page.Data.MediaBox.URY)-q.Y1.Pt())
	page.Path_l(q.X2.Pt(), float64(page.Data.MediaBox.URY)-q.Y2.Pt())
	page.Path_S()
	return nil
}
