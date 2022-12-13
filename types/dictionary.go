package types

import "bytes"

// PDF Reference 1.4, 3.2.6 Dictionary Objects

type Dictionary map[Name]Object

func (q Dictionary) ToRawBytes() []byte {
	sb := bytes.Buffer{}
	sb.WriteString("<<\n")
	for k, v := range q {
		sb.Write(k.ToRawBytes())
		sb.WriteString(" ")
		sb.Write(v.ToRawBytes())
		sb.WriteString("\n")
	}
	sb.WriteString(">>\n")
	return sb.Bytes()
}

func (q Dictionary) Copy(copyRef func(reference Reference) Reference) Object {
	c := make(Dictionary, len(q))
	for i, a := range q {
		c[i] = Copy(a, copyRef)
	}
	return c
}

func (q Dictionary) Equal(obj Object) bool {
	a, ok := obj.(Dictionary)
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

func (q Dictionary) GetValue(n Name, file Resolver) (Object, bool) {
	v, ok := q[n]
	if !ok {
		return nil, false
	}
	v2, err := file.ResolveReference(v)
	if err != nil {
		return nil, false
	}
	return v2, true
}
