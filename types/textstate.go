package types

// PDF Reference 1.4, Table 5.1 Text state parameters

type TextState struct {
	Tc    float64       // Character spacing
	Tw    float64       // Word spacing
	Th    float64       // Horizontal scaling
	Tl    float64       // Leading
	Tf    Name          // Text font
	Tfs   float64       // Text font size
	Tmode RenderingMode // Text rendering mode
	Trise float64       // Text rise
	Tk    float64       // Text knockout
}

func NewTextState() *TextState {
	return &TextState{
		Th: 100,
	}
}
