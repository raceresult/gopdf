package types

import "strconv"

// PDF Reference 1.4, 3.2.9 Indirect Objects

type Reference struct {
	Number     int
	Generation int
}

func (q Reference) ToRawBytes() []byte {
	return []byte(strconv.Itoa(q.Number) + " " + strconv.Itoa(q.Generation) + " R")
}

func (q Reference) Copy(copyRef func(reference Reference) Reference) Object {
	if q.Number == 0 {
		return q
	}
	return copyRef(q)
}

func (q Reference) Equal(obj Object) bool {
	a, ok := obj.(Reference)
	if !ok {
		return false
	}
	if q.Number != a.Number {
		return false
	}
	if q.Generation != a.Generation {
		return false
	}
	return true
}
