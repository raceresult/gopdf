package types

import "golang.org/x/text/encoding/unicode"

// PDF Reference 1.4, 3.8.1 Text Strings

type TextString string

func (q TextString) ToRawBytes() []byte {
	sn, _ := unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewEncoder().String(string(q))
	return []byte(sn)
}

func (q TextString) Copy(_ func(reference Reference) Reference) Object {
	return q
}
