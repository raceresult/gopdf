package types

import (
	"bytes"
	"strconv"
)

// PDF Reference 1.4, 3.2.9 Indirect Objects

type RawIndirectObject struct {
	Number     int
	Generation int
	Data       []byte
}

func (q RawIndirectObject) ToRawBytes() []byte {
	var sb bytes.Buffer
	sb.WriteString(strconv.Itoa(q.Number) + " " + strconv.Itoa(q.Generation) + " obj\n")

	sb.Write(q.Data)

	sb.WriteString("endobj\n")
	return sb.Bytes()
}

type IndirectObject struct {
	Number     int
	Generation int
	Data       Object
}
