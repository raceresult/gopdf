package pdf

import "github.com/raceresult/gopdf/types"

// PDF Reference 1.4, Table 5.10 Type 3 font operators

// Type3Font_d0 sets width information for the glyph and declare that the glyph descrip-
// tion specifies both its shape and its color. (Note that this operator name
// ends in the digit 0.) wx specifies the horizontal displacement in the glyph
// coordinate system; it must be consistent with the corresponding width
// in the font’s Widths array. w y must be 0 (see Section 5.1.3, “Glyph Posi-
// tioning and Metrics”).
// This operator is permitted only in a content stream appearing in a
// Type 3 font’s CharProcs dictionary. It is typically used only if the glyph
// description executes operators to set the color explicitly.
func (q *Page) Type3Font_d0(wx, wy float64) {
	q.AddCommand("d0", types.Number(wx), types.Number(wy))
}

// Type3Font_d1 sets width and bounding box information for the glyph and declare that
// the glyph description specifies only shape, not color. (Note that this
// operator name ends in the digit 1.) wx specifies the horizontal displace-
// ment in the glyph coordinate system; it must be consistent with the
// corresponding width in the font’s Widths array. w y must be 0 (see Sec-
// tion 5.1.3, “Glyph Positioning and Metrics”).
// llx and ll y are the coordinates of the lower-left corner, and ur x and ur y the
// upper-right corner, of the glyph bounding box. The glyph bounding box
// is the smallest rectangle, oriented with the axes of the glyph coordinate
// system, that completely encloses all marks placed on the page as a result
// of executing the glyph’s description. The declared bounding box must be
// correct—in other words, sufficiently large to enclose the entire glyph. If
// any marks fall outside this bounding box, the result is unpredictable.
// A glyph description that begins with the d1 operator should not execute
// any operators that set the color (or other color-related parameters) in
// the graphics state; any use of such operators will be ignored. The glyph
// description is executed solely to determine the glyph’s shape; its color is
// determined by the graphics state in effect each time this glyph is painted
// by a text-showing operator. For the same reason, the glyph description
// may not include an image; however, an image mask is acceptable, since it
// merely defines a region of the page to be painted with the current color.
// This operator is permitted only in a content stream appearing in a
// Type 3 font’s CharProcs dictionary.
func (q *Page) Type3Font_d1(wx, wy, llx, lly, urx, ury float64) {
	q.AddCommand("d1",
		types.Number(wx), types.Number(wy),
		types.Number(llx), types.Number(lly),
		types.Number(urx), types.Number(ury))
}
