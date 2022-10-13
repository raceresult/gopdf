package pdffile

import (
	"io"

	"github.com/raceresult/gopdf/types"
)

const DefaultVersion = 1.4

// File is a pdf file with a list of indirect objects
type File struct {
	Version float64
	Root    types.Reference
	Info    types.Reference
	ID      [2]string
	objects []*types.IndirectObject
}

// NewFile creates a new File object
func NewFile() *File {
	return &File{
		Version: DefaultVersion,
	}
}

// AddObject adds an object to the file and returns its reference
func (q *File) AddObject(obj types.Object) types.Reference {
	var no int
	if len(q.objects) != 0 {
		no = q.objects[len(q.objects)-1].Number + 1
	} else {
		no = 1
	}

	q.objects = append(q.objects, &types.IndirectObject{
		Number:     no,
		Generation: 0,
		Data:       obj,
	})
	return types.Reference{
		Number:     no,
		Generation: 0,
	}
}

// WriteTo writes the parsed to the given writer
func (q *File) WriteTo(w io.Writer) (int64, error) {
	var v1 Update
	v1.Objects = make([]types.RawIndirectObject, 0, len(q.objects))
	for _, obj := range q.objects {
		v1.Objects = append(v1.Objects, types.RawIndirectObject{
			Number:     obj.Number,
			Generation: obj.Generation,
			Data:       obj.Data.ToRawBytes(),
		})
	}
	v1.Trailer = types.Trailer{
		Root: q.Root,
		Info: q.Info,
		ID:   [2]types.String{types.String(q.ID[0]), types.String(q.ID[1])},
	}

	// output
	r := NewRawFile(q.Version)
	r.AddUpdate(v1)
	return r.WriteTo(w)
}
