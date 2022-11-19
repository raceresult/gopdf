package types

import (
	"strconv"
	"strings"
)

// PDF Reference 1.4, 3.2.2 Numeric Objects

type Number float64

func (q Number) ToRawBytes() []byte {
	s := strconv.FormatFloat(float64(q), 'f', 3, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return []byte(s)
}

func (q Number) Copy(_ func(reference Reference) Reference) Object {
	return q
}

func (q Number) Equal(obj Object) bool {
	a, ok := obj.(Number)
	if !ok {
		return false
	}
	return q == a
}
