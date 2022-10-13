package pdf

import "github.com/raceresult/gopdf/types"

// PDF Reference 1.4, Table 4.10 Path-painting operators

// Path_S strokes the path.
func (q *Page) Path_S() {
	q.AddProcSets(types.ProcSetPDF)
	q.AddCommand("S")
}

// Path_s closes and stroke the path. This operator has the same effect as the sequence h S.
func (q *Page) Path_s() {
	q.AddProcSets(types.ProcSetPDF)
	q.AddCommand("s")
}

// Path_f fills the path, using the nonzero winding number rule to determine the region to fill
// (see “Nonzero Winding Number Rule” on page 169).
func (q *Page) Path_f() {
	q.AddProcSets(types.ProcSetPDF)
	q.AddCommand("f")
}

// Path_F : Equivalent to f; included only for compatibility. Although applications that read
// PDF files must be able to accept this operator, those that generate PDF files should
// use f instead
func (q *Page) Path_F() {
	q.AddProcSets(types.ProcSetPDF)
	q.AddCommand("F")
}

// Path_fstar fills the path, using the even-odd rule to determine the region to fill (see “Even-Odd
// Rule” on page 170).
func (q *Page) Path_fstar() {
	q.AddCommand("f*")
}

// Path_B fill and then strokes the path, using the nonzero winding number rule to determine
// the region to fill. This produces the same result as constructing two identical path
// objects, painting the first with f and the second with S. Note, however, that the fill-
// ing and stroking portions of the operation consult different values of several graph-
// ics state parameters, such as the current color. See also “Special Path-Painting
// Considerations” on page 462.
func (q *Page) Path_B() {
	q.AddProcSets(types.ProcSetPDF)
	q.AddCommand("B")
}

// Path_Bstar fills and then strokes the path, using the even-odd rule to determine the region to fill.
// This operator produces the same result as B, except that the path is filled as if with
// f* instead of f. See also “Special Path-Painting Considerations” on page 462.
func (q *Page) Path_Bstar() {
	q.AddProcSets(types.ProcSetPDF)
	q.AddCommand("B*")
}

// Path_b closes, fills, and then strokes the path, using the nonzero winding number rule to de-
// termine the region to fill. This operator has the same effect as the sequence h B. See
// also “Special Path-Painting Considerations” on page 462.
func (q *Page) Path_b() {
	q.AddProcSets(types.ProcSetPDF)
	q.AddCommand("b")
}

// Path_bstar closes, fills, and then strokes the path, using the even-odd rule to determine the re-
// gion to fill. This operator has the same effect as the sequence h B*. See also “Special
// Path-Painting Considerations” on page 462.
func (q *Page) Path_bstar() {
	q.AddProcSets(types.ProcSetPDF)
	q.AddCommand("b*")
}

// Path_n ends the path object without filling or stroking it. This operator is a “path-painting
// no-op,” used primarily for the side effect of changing the current clipping path (see
// Section 4.4.3, “Clipping Path Operators”).
func (q *Page) Path_n() {
	q.AddProcSets(types.ProcSetPDF)
	q.AddCommand("n")
}
