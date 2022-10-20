package types

// PDF Reference 1.4, 3.2.1 Boolean Objects

type Boolean bool

func (q Boolean) ToRawBytes() []byte {
	if q {
		return []byte("true")
	}
	return []byte("false")
}

func (q Boolean) Copy(_ func(reference Reference) Reference) Object {
	return q
}
