package types

// PDF Reference 1.4, 3.2.8 Null Object

type Null struct{}

func (q Null) ToRawBytes() []byte {
	return []byte("null")
}
