package gopdf

// LinearMeasure represents a linear measure like mm or inch.
type LinearMeasure int

const (
	UnitPt LinearMeasure = iota
	UnitMM
	UnitInch
)

// Length represents a length using a certain measure
type Length struct {
	Value float64
	Unit  LinearMeasure
}

// MM creates a Length measured in mm
func MM(mm float64) Length {
	return Length{
		Value: mm,
		Unit:  UnitMM,
	}
}

// Inch creates a Length measured in inch
func Inch(inch float64) Length {
	return Length{
		Value: inch,
		Unit:  UnitInch,
	}
}

// Pt creates a Length measured in pt
func Pt(pt float64) Length {
	return Length{
		Value: pt,
		Unit:  UnitPt,
	}
}

// Pt converts the value to pt
func (q Length) Pt() float64 {
	switch q.Unit {
	case UnitMM:
		return q.Value * 2.83464567
	case UnitInch:
		return q.Value * 72
	default:
		return q.Value
	}
}

// Mm converts the value to millimeters
func (q Length) Mm() float64 {
	switch q.Unit {
	case UnitPt:
		return q.Value * 25.4 / 72
	case UnitInch:
		return q.Value * 25.4
	default:
		return q.Value
	}
}
