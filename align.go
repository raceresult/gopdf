package gopdf

// Note: PDF does not understand text align. For example, if you want right-aligned lines of text,
// you need to place (left-aligned) texts on the appropriate x positions.
// That is why these types/constants are in the builder package - only this package understands
// the concept of alignment.

type HorizontalAlign int

const (
	HorizontalAlignLeft   HorizontalAlign = 0
	HorizontalAlignCenter HorizontalAlign = 1
	HorizontalAlignRight  HorizontalAlign = 2
)

type VerticalAlign int

const (
	VerticalAlignTop    VerticalAlign = 0
	VerticalAlignMiddle VerticalAlign = 1
	VerticalAlignBottom VerticalAlign = 2
)
