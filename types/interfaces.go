package types

type Object interface {
	ToRawBytes() []byte
	Equal(object Object) bool
}

type copyable interface {
	Copy(copyRef func(reference Reference) Reference) Object
}

func Copy(obj Object, copyRef func(reference Reference) Reference) Object {
	if ac, ok := obj.(copyable); ok {
		return ac.Copy(copyRef)
	} else {
		return obj
	}
}

func Equal(obj1, obj2 Object) bool {
	if obj1 == nil {
		return obj2 == nil
	}
	return obj1.Equal(obj2)
}
