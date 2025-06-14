package types

import (
	"strconv"
	"strings"
)

// PDF Reference 1.4, 3.2.4 Name Objects

type Name string

func (q Name) ToRawBytes() []byte {
	s := "/"
	for _, r := range q {
		if r < 42 || r == '/' {
			h := strings.ToUpper(strconv.FormatInt(int64(r), 16))
			if len(h) == 1 {
				h = "0" + h
			}
			s += "#" + h
		} else {
			s += string(r)
		}
	}

	return []byte(s)
}

func (q Name) Copy(_ func(reference Reference) Reference) Object {
	return q
}

func (q Name) Equal(obj Object) bool {
	a, ok := obj.(Name)
	if !ok {
		return false
	}
	return q == a
}
