package pdf

import "github.com/raceresult/gopdf/types"

// PDF Reference 1.4, Table 4.7 Graphics state operators

// GraphicsState_q saves the current graphics state on the graphics state stack (see “Graphics
// State Stack” on page 152).
func (q *Page) GraphicsState_q() {
	currState := *q.graphicsState
	q.graphicsStateStack = append(q.graphicsStateStack, q.graphicsState)
	q.graphicsState = &currState

	q.AddCommand("q")
}

// GraphicsState_Q restores the graphics state by removing the most recently saved state from
// the stack and making it the current state (see “Graphics State Stack” on
// page 152).
func (q *Page) GraphicsState_Q() {
	if len(q.graphicsStateStack) != 0 {
		q.graphicsState = q.graphicsStateStack[len(q.graphicsStateStack)-1]
		q.graphicsStateStack = q.graphicsStateStack[:len(q.graphicsStateStack)-1]
	}

	q.AddCommand("Q")
}

// GraphicsState_cm modifies the current transformation matrix (CTM) by concatenating the
// specified matrix (see Section 4.2.1, “Coordinate Spaces”). Although the
// operands specify a matrix, they are written as six separate numbers, not as
// an array.
func (q *Page) GraphicsState_cm(a, b, c, d, e, f float64) {
	q.AddCommand("cm", types.Number(a), types.Number(b), types.Number(c), types.Number(d), types.Number(e), types.Number(f))
}

// GraphicsState_w sets the line width in the graphics state (see “Line Width” on page 152).
func (q *Page) GraphicsState_w(lineWidth float64) {
	if q.graphicsState.LineWidth == lineWidth {
		return
	}
	q.graphicsState.LineWidth = lineWidth

	q.AddCommand("w", types.Number(lineWidth))
}

// GraphicsState_J sets the line cap style in the graphics state (see “Line Cap Style” on page 153)
func (q *Page) GraphicsState_J(lineCap int) {
	if q.graphicsState.LineCap == lineCap {
		return
	}
	q.graphicsState.LineCap = lineCap

	q.AddCommand("J", types.Int(lineCap))
}

// GraphicsState_j sets the line join style in the graphics state (see “Line Join Style” on page 153).
func (q *Page) GraphicsState_j(lineJoin int) {
	if q.graphicsState.LineJoin == lineJoin {
		return
	}
	q.graphicsState.LineJoin = lineJoin

	q.AddCommand("j", types.Int(lineJoin))
}

// GraphicsState_M sets the miter limit in the graphics state (see “Miter Limit” on page 153).
func (q *Page) GraphicsState_M(miterLimit float64) {
	if q.graphicsState.MiterLimit == miterLimit {
		return
	}
	q.graphicsState.MiterLimit = miterLimit

	q.AddCommand("M", types.Number(miterLimit))
}

// GraphicsState_d sets the line dash pattern in the graphics state (see “Line Dash Pattern” on page 155)
func (q *Page) GraphicsState_d(dashArray []float64, dashPhase float64) {
	if q.graphicsState.DashPhase == dashPhase && len(q.graphicsState.DashArray) == len(dashArray) {
		var different bool
		for i, v := range dashArray {
			if v != q.graphicsState.DashArray[i] {
				different = true
				break
			}
		}
		if !different {
			return
		}
	}
	q.graphicsState.DashPhase = dashPhase
	q.graphicsState.DashArray = dashArray

	arr := make(types.Array, 0, len(dashArray))
	for _, v := range dashArray {
		arr = append(arr, types.Number(v))
	}
	q.AddCommand("d", arr, types.Number(dashPhase))
}

// GraphicsState_ri sets the color rendering intent in the graphics state (see “Rendering Intents” on page 197).
func (q *Page) GraphicsState_ri(intent string) {
	if q.graphicsState.Intent == intent {
		return
	}
	q.graphicsState.Intent = intent

	q.AddCommand("ri", types.Name(intent))
}

// GraphicsState_i sets the flatness tolerance in the graphics state (see Section 6.5.1, “Flatness
// Tolerance”). flatness is a number in the range 0 to 100; a value of 0 speci-
// fies the output device’s default flatness tolerance.
func (q *Page) GraphicsState_i(flatness float64) {
	if q.graphicsState.Flatness == flatness {
		return
	}
	q.graphicsState.Flatness = flatness

	q.AddCommand("i", types.Number(flatness))
}

// GraphicsState_gs sets the specified parameters in the graphics state. dictName is
// the name of a graphics state parameter dictionary in the ExtGState sub-
// dictionary of the current resource dictionary (see the next section).
func (q *Page) GraphicsState_gs(dictName types.Name) {
	// set q.graphicsState to some impossible values so that new operators will not be skipped
	q.graphicsState.NonStrokingColor = graphicsStateColor{}
	q.graphicsState.StrokingColor = graphicsStateColor{}
	q.graphicsState.LineWidth = -1
	q.graphicsState.LineCap = -1
	q.graphicsState.Flatness = -1
	q.graphicsState.Intent = "-"
	q.graphicsState.LineJoin = -1
	q.graphicsState.MiterLimit = -1

	q.AddCommand("gs", dictName)
}
