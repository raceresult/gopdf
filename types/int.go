package types

import "strconv"

// PDF Reference 1.4, 3.2.2 Numeric Objects

type Int int

func (q Int) ToRawBytes() []byte {
	return []byte(strconv.Itoa(int(q)))
}
