package pdf

import "github.com/raceresult/gopdf/types"

// PDF Reference 1.4, Table 5.2 Text state operators

// TextState_Tc sets the character spacing, Tc , to charSpace, which is a number expressed in un-
// scaled text space units. Character spacing is used by the Tj, TJ, and ' operators.
// Initial value: 0.
// No-Op if identical to value in current text state
func (q *Page) TextState_Tc(charSpace float64) {
	if q.graphicsState.TextState.Tc == charSpace {
		return
	}
	q.graphicsState.TextState.Tc = charSpace

	q.AddCommand("Tc", types.Int(charSpace))
}

// TextState_Tw set the word spacing, Tw , to wordSpace, which is a number expressed in unscaled
// text space units. Word spacing is used by the Tj, TJ, and ' operators.
// Initial value: 0.
// No-Op if identical to value in current text state
func (q *Page) TextState_Tw(wordSpace float64) {
	if q.graphicsState.TextState.Tw == wordSpace {
		return
	}
	q.graphicsState.TextState.Tw = wordSpace

	q.AddCommand("Tw", types.Int(wordSpace))
}

// TextState_Tz sets the horizontal scaling, Th , to (scale  ̃ 100). scale is a number specifying the
// percentage of the normal width.
// Initial value: 100 (normal width).
// No-Op if identical to value in current text state
func (q *Page) TextState_Tz(scale float64) {
	if q.graphicsState.TextState.Th == scale {
		return
	}
	q.graphicsState.TextState.Th = scale

	q.AddCommand("Tz", types.Int(scale))
}

// TextState_TL sets the text leading, Tl , to leading, which is a number expressed in unscaled text
// space units. Text leading is used only by the T*, ', and " operators.
// Initial value: 0.
// No-Op if identical to value in current text state
func (q *Page) TextState_TL(leading float64) {
	if q.graphicsState.TextState.Tl == leading {
		return
	}
	q.graphicsState.TextState.Tl = leading

	q.AddCommand("TL", types.Int(leading))
}

// TextState_Tf set the text font, Tf , to font and the text font size, Tfs , to size. font is the name of a
// font resource in the Font subdictionary of the current resource dictionary; size is
// a number representing a scale factor. There is no initial value for either font or
// size; they must be specified explicitly using Tf before any text is shown.
// No-Op if identical to value in current text state
func (q *Page) TextState_Tf(font FontHandler, fontSize float64) {
	n := q.AddFont(font)
	if q.graphicsState.TextState.Tf == n && q.graphicsState.TextState.Tfs == fontSize {
		return
	}
	q.graphicsState.TextState.Tf = n
	q.graphicsState.TextState.Tfs = fontSize

	q.currFont = font
	q.AddCommand("Tf", n, types.Number(fontSize))
}

// TextState_Tr sets the text rendering mode, Tmode , to render, which is an integer.
// Initial value: 0.
// No-Op if identical to value in current text state
func (q *Page) TextState_Tr(mode types.RenderingMode) {
	if q.graphicsState.TextState.Tmode == mode {
		return
	}
	q.graphicsState.TextState.Tmode = mode

	q.AddCommand("Tr", types.Int(mode))
}

// TextState_Ts sets the text rise, Trise , to rise, which is a number expressed in unscaled text space
// units.
// Initial value: 0
// No-Op if identical to value in current text state
func (q *Page) TextState_Ts(rise float64) {
	if q.graphicsState.TextState.Trise == rise {
		return
	}
	q.graphicsState.TextState.Trise = rise

	q.AddCommand("Ts", types.Number(rise))
}
