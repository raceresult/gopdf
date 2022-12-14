package types

import "strconv"

// PDF Reference 1.4, 3.2.2 Numeric Objects

type Int int

func (q Int) ToRawBytes() []byte {
	return []byte(strconv.Itoa(int(q)))
}

func (q Int) Copy(_ func(reference Reference) Reference) Object {
	return q
}

func (q Int) Equal(obj Object) bool {
	a, ok := obj.(Int)
	if !ok {
		return false
	}
	return q == a
}
