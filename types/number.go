package types

import "strconv"

// PDF Reference 1.4, 3.2.2 Numeric Objects

type Number float64

func (q Number) ToRawBytes() []byte {
	s := strconv.FormatFloat(float64(q), 'f', 3, 64)
	return []byte(s)
}
