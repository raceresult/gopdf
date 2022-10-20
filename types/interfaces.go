package types

type Object interface {
	ToRawBytes() []byte
}

type PDFFile interface {
	AddObject(obj Object) Reference
	GetObject(ref Reference) (Object, error)
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
