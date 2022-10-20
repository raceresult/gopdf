package types

import (
	"strings"
)

// PDF Reference 1.4, 3.2.3 String Objects

type String string

var replace = strings.NewReplacer(
	"\n", "\\n",
	"\r", "\\r",
	"\t", "\\t",
	"\b", "\\b",
	"\f", "\\f",
	"(", "\\(",
	")", "\\)",
	"\\", "\\\\",
)

func (q String) ToRawBytes() []byte {
	s := replace.Replace(string(q))
	return []byte("(" + s + ")")
}

func (q String) Copy(_ func(reference Reference) Reference) Object {
	return q
}
