package types

import "bytes"

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
