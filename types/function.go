package types

type Function struct{}

func (q Function) ToRawBytes() []byte {
	// todo
	return nil
}

func (q Function) Copy(_ func(reference Reference) Reference) Object {
	return Function{}
}

func (q Function) Equal(obj Object) bool {
	_, ok := obj.(Function)
	if !ok {
		return false
	}
	return true
}
