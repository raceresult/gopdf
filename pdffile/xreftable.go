package pdffile

import (
	"bytes"
	"fmt"
	"strconv"
)

// PDF Reference 1.4, 3.4.3 Cross-Reference Table

type XRefTable []XRefTableSection

func (q *XRefTable) clear() {
	*q = []XRefTableSection{
		{
			Start: 0,
			Count: 1,
			Entries: []XRefTableEntry{{
				Start:      0,
				Generation: 65535,
				Free:       true,
			}},
		},
	}

}

// XRefTableSection is a section of a xref-table
type XRefTableSection struct {
	Start   int
	Count   int
	Entries []XRefTableEntry
}

func (q *XRefTableSection) ToRawBytes() []byte {
	if len(q.Entries) == 0 {
		return nil
	}

	var sb bytes.Buffer
	sb.WriteString(strconv.Itoa(q.Start) + " " + strconv.Itoa(q.Count) + "\n")
	for _, entry := range q.Entries {
		sb.Write(entry.ToRawBytes())
	}
	return sb.Bytes()
}

type XRefTableEntry struct {
	Start      int64
	Generation int
	Free       bool
}

func (q *XRefTableEntry) ToRawBytes() []byte {
	var sb bytes.Buffer
	sb.WriteString(fmt.Sprintf("%010d", q.Start) + " " + fmt.Sprintf("%05d", q.Generation))
	if q.Free {
		sb.WriteString(" f \n")
	} else {
		sb.WriteString(" n \n")
	}
	return sb.Bytes()
}

func (q *XRefTable) add(number int, generation int, free bool, pos int64) {
	var last *XRefTableSection
	if len(*q) != 0 {
		last = &(*q)[len(*q)-1]
	}
	if last == nil || last.Start+last.Count != number {
		*q = append(*q, XRefTableSection{Start: number})
		last = &(*q)[len(*q)-1]
	}
	last.Count++
	last.Entries = append(last.Entries, XRefTableEntry{
		Start:      pos,
		Generation: generation,
		Free:       free,
	})
}

func (q *XRefTable) ToRawBytes() []byte {
	var sb bytes.Buffer

	sb.WriteString("xref\n")
	for _, section := range *q {
		sb.Write(section.ToRawBytes())
	}
	return sb.Bytes()
}
