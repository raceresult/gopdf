package pdf

import "github.com/raceresult/gopdf/types"

// graphicsState keeps track internally of the graphics state when adding PDF operators, so that commands that would not
// change the current state can be skipped
type graphicsState struct {
	StrokingColor    graphicsStateColor
	NonStrokingColor graphicsStateColor
	LineWidth        float64
	LineCap          int
	LineJoin         int
	MiterLimit       float64
	Intent           string
	Flatness         float64
	DashPhase        float64
	DashArray        []float64
}

// graphicsStateColor handles colors within the graphicsState
type graphicsStateColor struct {
	Name    types.ColorSpaceFamily
	Values  []float64
	AddName types.Name
}

// Equal checks if the graphicsStateColor is equal to the given values
func (q graphicsStateColor) Equal(name types.ColorSpaceFamily, addName types.Name, values ...float64) bool {
	if name != q.Name {
		return false
	}
	if addName != q.AddName {
		return false
	}
	if len(q.Values) != len(values) {
		return false
	}
	for i, v := range q.Values {
		if v != values[i] {
			return false
		}
	}
	return true
}

// SetIfNotEqual updates the graphics state if it differs from the given values
func (q graphicsStateColor) SetIfNotEqual(name types.ColorSpaceFamily, addName types.Name, values ...float64) bool {
	if q.Equal(name, addName, values...) {
		return false
	}
	q.Name = name
	q.AddName = addName
	q.Values = values
	return true
}
