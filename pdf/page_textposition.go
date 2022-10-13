package pdf

import "git.rrdc.de/lib/gopdf/types"

// PDF Reference 1.4, Table 5.5 Text-positioning operators

// TextPosition_Td Move to the start of the next line, offset from the start of the current line by
// (tx , ty ). tx and ty are numbers expressed in unscaled text space units.
func (q *Page) TextPosition_Td(tx, ty float64) {
	q.AddCommand("Td", types.Number(tx), types.Number(ty))
}

// TextPosition_TD moves to the start of the next line, offset from the start of the current line by
// (t x , t y ). As a side effect, this operator sets the leading parameter in the text state.
// This operator has the same effect as the code
// −t y TL
// tx ty Td
func (q *Page) TextPosition_TD(tx, ty float64) {
	q.AddCommand("TD", types.Number(tx), types.Number(ty))
}

// TextPosition_Tm set the text matrix, Tm , and the text line matrix, Tlm
func (q *Page) TextPosition_Tm(a, b, c, d, e, f float64) {
	q.AddCommand("Tm", types.Number(a), types.Number(b), types.Number(c), types.Number(d), types.Number(e), types.Number(f))
}

// TextPosition_Tstar moves to the start of the next line. This operator has the same effect as the code
// 0 Tl Td
// where Tl is the current leading parameter in the text state.
func (q *Page) TextPosition_Tstar() {
	q.AddCommand("T*")
}
