package types

import (
	"bytes"
)

// PDF Reference 1.4, 3.2.5 Array Objects

type Array []Object

func (q Array) ToRawBytes() []byte {
	var sb bytes.Buffer
	sb.WriteString("[")
	for i, v := range q {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.Write(v.ToRawBytes())
	}
	sb.WriteString("]")
	return sb.Bytes()
}

func (q Array) Copy(copyRef func(reference Reference) Reference) Object {
	c := make(Array, len(q))
	for i, a := range q {
		c[i] = Copy(a, copyRef)
	}
	return c
}

func (q Array) Equal(obj Object) bool {
	a, ok := obj.(Array)
	if !ok {
		return false
	}
	if len(q) != len(a) {
		return false
	}
	for i := range q {
		if !Equal(q[i], a[i]) {
			return false
		}
	}
	return true
}
