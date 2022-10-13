package pdf

import "git.rrdc.de/lib/gopdf/types"

// PDF Reference 1.4, Table  4.24 Shading operator

// Shading_sh paints the shape and color shading described by a shading dictionary, sub-
// ject to the current clipping path. The current color in the graphics state is neither
// used nor altered. The effect is different from that of painting a path using a shading
// pattern as the current color.
// name is the name of a shading dictionary resource in the Shading subdictionary of
// the current resource dictionary (see Section 3.7.2, “Resource Dictionaries”). All co-
// ordinates in the shading dictionary are interpreted relative to the current user
// space. (By contrast, when a shading dictionary is used in a type 2 pattern, the
// coordinates are expressed in pattern space.) All colors are interpreted in the color
// space identified by the shading dictionary’s ColorSpace entry (see Table 4.25 on
// page 234). The Background entry, if present, is ignored.
// This operator should be applied only to bounded or geometrically defined shad-
// ings. If applied to an unbounded shading, it will paint the shading’s gradient fill
// across the entire clipping region, which may be time-consuming.
func (q *Page) Shading_sh(name types.Name) {
	q.AddCommand("sh", name)
}
