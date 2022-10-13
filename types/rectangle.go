package types

// PDF Reference 1.4, 3.8.3 Rectangles

type Rectangle struct {
	LLX Number
	LLY Number
	URX Number
	URY Number
}

func (q Rectangle) ToRawBytes() []byte {
	a := Array{q.LLX, q.LLY, q.URX, q.URY}
	return a.ToRawBytes()
}
