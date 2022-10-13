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
