package types

// PDF Reference 1.4, Table 5.1 Text state parameters

type TextState struct {
	// Character spacing
	Tc float64

	// Word spacing
	Tw float64

	// Horizontal scaling
	Th float64

	// Leading
	Tl float64

	// Text font
	Tf Name

	// Text font size
	Tfs float64

	// Text rendering mode
	Tmode RenderingMode

	// Text rise
	Trise float64

	// Text knockout
	Tk float64
}
