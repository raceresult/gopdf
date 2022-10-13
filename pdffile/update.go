package pdffile

import (
	"io"

	"git.rrdc.de/lib/gopdf/types"
)

type Update struct {
	Objects   []types.RawIndirectObject
	xRefTable XRefTable
	Trailer   types.Trailer
}

func (q *Update) writeTo(w io.Writer, offset int64) (int64, int64, error) {
	// objects
	var n int64
	q.xRefTable.clear()
	for _, obj := range q.Objects {
		q.xRefTable.add(obj.Number, obj.Generation, false, n+offset)

		nn, err := w.Write(obj.ToRawBytes())
		if err != nil {
			return n, 0, err
		}
		n += int64(nn)

		if obj.Number+1 > q.Trailer.Size {
			q.Trailer.Size = obj.Number + 1
		}
	}

	// xref table
	startXRef := n + offset
	nn, err := w.Write(q.xRefTable.ToRawBytes())
	if err != nil {
		return n, startXRef, err
	}
	n += int64(nn)

	// trailer
	n2, err := w.Write(q.Trailer.ToRawBytes())
	if err != nil {
		return n, startXRef, err
	}
	n += int64(n2)

	// return without error
	return n, startXRef, nil
}
