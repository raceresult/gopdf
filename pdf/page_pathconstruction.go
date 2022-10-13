package pdf

import "github.com/raceresult/gopdf/types"

// PDF Reference 1.4, Table 4.9 Path construction operators

// Path_m begins a new subpath by moving the current point to coordinates
// (x, y), omitting any connecting line segment. If the previous path
// construction operator in the current path was also m, the new m
// overrides it; no vestige of the previous m operation remains in the
// path.
func (q *Page) Path_m(x, y float64) {
	q.AddCommand("m", types.Number(x), types.Number(y))
}

// Path_l appends a straight line segment from the current point to the point
// (x, y). The new current point is (x, y).
func (q *Page) Path_l(x, y float64) {
	q.AddCommand("l", types.Number(x), types.Number(y))
}

// Path_c appends a cubic Bézier curve to the current path. The curve extends
// from the current point to the point (x3, y3), using (x1, y1) and
// (x2, y2) as the Bézier control points (see “Cubic Bézier Curves,” be-
// low). The new current point is (x3, y3).
func (q *Page) Path_c(x1, y1, x2, y2, x3, y3 float64) {
	q.AddCommand("c",
		types.Number(x1), types.Number(y1),
		types.Number(x2), types.Number(y2),
		types.Number(x3), types.Number(y3))
}

// Path_v appends a cubic Bézier curve to the current path. The curve extends
// from the current point to the point (x3, y3), using the current point
// and (x2, y2) as the Bézier control points (see “Cubic Bézier Curves,”
// below). The new current point is (x3, y3).
func (q *Page) Path_v(x2, y2, x3, y3 float64) {
	q.AddCommand("v", types.Number(x2), types.Number(y2), types.Number(x3), types.Number(y3))
}

// Path_y appends a cubic Bézier curve to the current path. The curve extends
// from the current point to the point (x3, y3), using (x1, y1) and
// (x3, y3) as the Bézier control points (see “Cubic Bézier Curves,” be-
// low). The new current point is (x3, y3).
func (q *Page) Path_y(x1, y1, x3, y3 float64) {
	q.AddCommand("y", types.Number(x1), types.Number(y1), types.Number(x3), types.Number(y3))
}

// Path_h closes the current subpath by appending a straight line segment
// from the current point to the starting point of the subpath. This
// operator terminates the current subpath; appending another seg-
// ment to the current path will begin a new subpath, even if the new
// segment begins at the endpoint reached by the h operation. If the
// current subpath is already closed, h does nothing.
func (q *Page) Path_h() {
	q.AddCommand("h")
}

// Path_re append a rectangle to the current path as a complete subpath, with
// lower-left corner (x, y) and dimensions width and height in user space.
func (q *Page) Path_re(x, y, width, height float64) {
	q.AddCommand("re", types.Number(x), types.Number(y), types.Number(width), types.Number(height))
}
