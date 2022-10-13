package pdf

import "github.com/raceresult/gopdf/types"

// PDF Reference 1.4, Table 5.2 Text state operators

// TextState_Tc sets the character spacing, Tc , to charSpace, which is a number expressed in un-
// scaled text space units. Character spacing is used by the Tj, TJ, and ' operators.
// Initial value: 0.
// No-Op if identical to value in current text state
func (q *Page) TextState_Tc(charSpace float64) {
	if q.textState.Tc == charSpace {
		return
	}
	q.textState.Tc = charSpace

	q.AddCommand("Tc", types.Int(charSpace))
}

// TextState_Tw set the word spacing, Tw , to wordSpace, which is a number expressed in unscaled
// text space units. Word spacing is used by the Tj, TJ, and ' operators.
// Initial value: 0.
// No-Op if identical to value in current text state
func (q *Page) TextState_Tw(wordSpace float64) {
	if q.textState.Tw == wordSpace {
		return
	}
	q.textState.Tw = wordSpace

	q.AddCommand("Tw", types.Int(wordSpace))
}

// TextState_Tz sets the horizontal scaling, Th , to (scale  ̃ 100). scale is a number specifying the
// percentage of the normal width.
// Initial value: 100 (normal width).
// No-Op if identical to value in current text state
func (q *Page) TextState_Tz(scale float64) {
	if q.textState.Th == scale {
		return
	}
	q.textState.Th = scale

	q.AddCommand("Tz", types.Int(scale))
}

// TextState_TL sets the text leading, Tl , to leading, which is a number expressed in unscaled text
// space units. Text leading is used only by the T*, ', and " operators.
// Initial value: 0.
// No-Op if identical to value in current text state
func (q *Page) TextState_TL(leading float64) {
	if q.textState.Tl == leading {
		return
	}
	q.textState.Tl = leading

	q.AddCommand("TL", types.Int(leading))
}

// TextState_Tf set the text font, Tf , to font and the text font size, Tfs , to size. font is the name of a
// font resource in the Font subdictionary of the current resource dictionary; size is
// a number representing a scale factor. There is no initial value for either font or
// size; they must be specified explicitly using Tf before any text is shown.
// No-Op if identical to value in current text state
func (q *Page) TextState_Tf(font FontHandler, fontSize float64) {
	n := q.AddFont(font)
	if q.textState.Tf == n && q.textState.Tfs == fontSize {
		return
	}
	q.textState.Tf = n
	q.textState.Tfs = fontSize

	q.currFont = font
	q.AddCommand("Tf", n, types.Number(fontSize))
}

// TextState_Tr sets the text rendering mode, Tmode , to render, which is an integer.
// Initial value: 0.
// No-Op if identical to value in current text state
func (q *Page) TextState_Tr(mode types.RenderingMode) {
	if q.textState.Tmode == mode {
		return
	}
	q.textState.Tmode = mode

	q.AddCommand("Tr", types.Int(mode))
}

// TextState_Ts sets the text rise, Trise , to rise, which is a number expressed in unscaled text space
// units.
// Initial value: 0
// No-Op if identical to value in current text state
func (q *Page) TextState_Ts(rise float64) {
	if q.textState.Trise == rise {
		return
	}
	q.textState.Trise = rise

	q.AddCommand("Ts", types.Number(rise))
}
