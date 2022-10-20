package types

import "time"

// PDF Reference 1.4, 3.8.2 Dates

type Date time.Time

func (q Date) IsZero() bool {
	return time.Time(q).IsZero()
}

func (q Date) ToRawBytes() []byte {
	var s String
	s = String("D:" + time.Time(q).Format("20060102150405-07'00'"))
	return s.ToRawBytes()
}

func (q Date) Copy(_ func(reference Reference) Reference) Object {
	return q
}
