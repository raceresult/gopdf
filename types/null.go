package types

// PDF Reference 1.4, 3.2.8 Null Object

type Null struct{}

func (q Null) ToRawBytes() []byte {
	return []byte("null")
}

func (q Null) Copy(_ func(reference Reference) Reference) Object {
	return q
}

func (q Null) Equal(obj Object) bool {
	_, ok := obj.(Null)
	return ok
}
