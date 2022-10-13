package types

// PDF Reference 1.4, Table 5.3 Text rendering modes

type RenderingMode int

const (
	RenderingModeFill                 RenderingMode = 0
	RenderingModeStroke               RenderingMode = 1
	RenderingModeFillAndStroke        RenderingMode = 2
	RenderingModeNeither              RenderingMode = 3
	RenderingModeFillAndClip          RenderingMode = 4
	RenderingModeStrokeAndClip        RenderingMode = 5
	RenderingModeFillAndStrokeAndClip RenderingMode = 6
	RenderingModeClip                 RenderingMode = 7
)
