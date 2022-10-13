package types

// PDF Reference 1.4, 3.2.1 Boolean Objects

type Boolean bool

func (q Boolean) ToRawBytes() []byte {
	if q {
		return []byte("true")
	}
	return []byte("false")
}
