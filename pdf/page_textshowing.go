package pdf

import "git.rrdc.de/lib/gopdf/types"

// PDF Reference 1.4, Table 5.6 Text-showing operators

// TextShowing_Tj shows a text string.
func (q *Page) TextShowing_Tj(s string) {
	if q.currFont != nil {
		s = q.currFont.Encode(s)
	}
	q.AddCommand("Tj", types.String(s))
}

// TextShowing_Pos moves to the next line and show a text string. This operator has the same effect as
// the code
// T*
// string T
func (q *Page) TextShowing_Pos(s string) {
	if q.currFont != nil {
		s = q.currFont.Encode(s)
	}
	q.AddCommand("'", types.String(s))
}

// TextShowing_Quot moves to the next line and show a text string, using a w as the word spacing and a c
// as the character spacing (setting the corresponding parameters in the text state).
// a w and a c are numbers expressed in unscaled text space units. This operator has
// the same effect as the code
// aw Tw
// ac Tc
// string '
func (q *Page) TextShowing_Quot(ac, aw float64, s string) {
	if q.currFont != nil {
		s = q.currFont.Encode(s)
	}
	q.AddCommand("\"", types.Number(ac), types.Number(aw), types.String(s))
}

// TextShowing_TJ shows one or more text strings, allowing individual glyph positioning (see imple-
// mentation note 40 in Appendix H). Each element of array can be a string or a
// number. If the element is a string, this operator shows the string. If it is a num-
// ber, the operator adjusts the text position by that amount; that is, it translates
// the text matrix, Tm . The number is expressed in thousandths of a unit of text
// space (see Section 5.3.3, “Text Space Details,” and implementation note 41 in
// Appendix H). This amount is subtracted from the current horizontal or vertical
// coordinate, depending on the writing mode. In the default coordinate system, a
// positive adjustment has the effect of moving the next glyph painted either to the
// left or down by the given amount. Figure 5.11 shows an example of the effect of
// passing offsets to TJ.
func (q *Page) TextShowing_TJ(array types.Array) {
	if q.currFont != nil {
		ne := make(types.Array, 0, len(array))
		for _, v := range array {
			if s, ok := v.(types.String); ok {
				ne = append(ne, types.String(q.currFont.Encode(string(s))))
			} else {
				ne = append(ne, v)
			}
		}
		array = ne
	}

	q.AddCommand("TJ", array...)
}
