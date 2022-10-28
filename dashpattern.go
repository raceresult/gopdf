package gopdf

import "github.com/raceresult/gopdf/pdf"

// DashPattern represents a dash pattern
type DashPattern struct {
	Phase Length
	Array []Length
}

// NewDashPattern creates a new DashPattern object
func NewDashPattern(phase Length, array ...Length) DashPattern {
	return DashPattern{
		Phase: phase,
		Array: array,
	}
}

// Build sets the dash pattern in the graphics state of the given page
func (q DashPattern) Build(page *pdf.Page) {
	arr := make([]float64, 0, len(q.Array))
	for _, v := range q.Array {
		arr = append(arr, v.Pt())
	}
	page.GraphicsState_d(arr, q.Phase.Pt())
}
