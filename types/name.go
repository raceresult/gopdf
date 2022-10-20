package types

import "strconv"

// PDF Reference 1.4, 3.2.4 Name Objects

type Name string

func (q Name) ToRawBytes() []byte {
	s := "/"
	for _, r := range q {
		if r < 42 {
			h := strconv.FormatInt(int64(r), 16)
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
