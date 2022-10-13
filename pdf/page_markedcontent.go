package pdf

import "github.com/raceresult/gopdf/types"

// PDF Reference 1.4, Table 9.8 Marked-content operators

// MarkedContent_MP designates a marked-content point. tag is a name object indicating the role or
// significance of the point.
func (q *Page) MarkedContent_MP(tag types.Name) {
	q.AddCommand("MP", tag)
}

// MarkedContent_DP designates a marked-content point with an associated property list. tag is a
// name object indicating the role or significance of the point; properties is
// either an inline dictionary containing the property list or a name object
// associated with it in the Properties subdictionary of the current resource
// dictionary (see Section 9.5.1, “Property Lists”).
func (q *Page) MarkedContent_DP(tag, properties types.Name) {
	q.AddCommand("DP", tag, properties)
}

// MarkedContent_BMC begins a marked-content sequence terminated by a balancing EMC operator.
// tag is a name object indicating the role or significance of the sequence.
func (q *Page) MarkedContent_BMC(tag types.Name) {
	q.AddCommand("BMC", tag)
}

// MarkedContent_BDC begins a marked-content sequence with an associated property list, terminated
// by a balancing EMC operator. tag is a name object indicating the role or signif-
// icance of the sequence; properties is either an inline dictionary containing the
// property list or a name object associated with it in the Properties subdiction-
// ary of the current resource dictionary (see Section 9.5.1, “Property Lists”).
func (q *Page) MarkedContent_BDC(tag, properties types.Name) {
	q.AddCommand("BDC", tag, properties)
}

// MarkedContent_EMC ends a marked-content sequence begun by a BMC or BDC operator.
func (q *Page) MarkedContent_EMC() {
	q.AddCommand("EMC")
}
