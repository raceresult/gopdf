package pdf

import "git.rrdc.de/lib/gopdf/types"

// PDF Reference 1.4, Table 5.4 Text object operators

// TextObjects_BT begins a text object, initializing the text matrix, Tm , and the text line matrix, Tlm , to
// the identity matrix. Text objects cannot be nested; a second BT cannot appear before an ET.
func (q *Page) TextObjects_BT() {
	q.AddProcSets(types.ProcSetText)
	q.AddCommand("BT")
}

// TextObjects_ET ends a text object, discarding the text matrix.
func (q *Page) TextObjects_ET() {
	q.AddCommand("ET")
}
