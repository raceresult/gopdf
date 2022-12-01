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
	if len(q.Data) != 0 && !isWhiteChar(q.Data[len(q.Data)-1]) {
		sb.WriteByte('\n')
	}

	sb.WriteString("endobj\n")
	return sb.Bytes()
}

type IndirectObject struct {
	Number     int
	Generation int
	Data       Object
}

func isWhiteChar(c byte) bool {
	switch c {
	case 0, 9, 10, 12, 13, 32:
		return true
	default:
		return false
	}
}
