package pdffile

import (
	"errors"
	"io"
	"strconv"

	"github.com/raceresult/gopdf/types"
)

const DefaultVersion = 1.4

// File is a pdf file with a list of indirect objects
type File struct {
	Version float64
	Root    types.Reference
	Info    types.Reference
	ID      [2]types.String
	objects []types.IndirectObject
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

	q.objects = append(q.objects, types.IndirectObject{
		Number:     no,
		Generation: 0,
		Data:       obj,
	})
	return types.Reference{
		Number:     no,
		Generation: 0,
	}
}

// GetObject returns an object from the file
func (q *File) GetObject(ref types.Reference) (types.Object, error) {
	for _, obj := range q.objects {
		if obj.Number == ref.Number && obj.Generation == ref.Generation {
			return obj.Data, nil
		}
	}
	return nil, errors.New("object " + strconv.Itoa(ref.Number) + "/" + strconv.Itoa(ref.Generation) + " not found")
}

// WriteTo writes the parsed to the given writer
func (q *File) WriteTo(w io.Writer) (int64, error) {
	// check version
	if q.Version == 0 {
		return 0, errors.New("no pdf version set")
	}

	// build raw objects
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
		ID:   [2]types.String{q.ID[0], q.ID[1]},
	}

	// header (3.4.1 File Header)
	var n int64
	n1, err := w.Write([]byte("%PDF-" + strconv.FormatFloat(q.Version, 'f', -1, 64) + "\n"))
	if err != nil {
		return n, err
	}
	n += int64(n1)

	// Note: If a PDF file contains binary data, as most do (see Section 3.1, “Lexical Con-
	// ventions”), it is recommended that the header line be immediately followed by a
	// comment line containing at least four binary characters—that is, characters whose
	// codes are 128 or greater. This will ensure proper behavior of file transfer applications
	// that inspect data near the beginning of a file to determine whether to treat the file’s
	// contents as text or as binary.
	n1, err = w.Write([]byte{'%', 250, 251, 252, 253, '\n'})
	if err != nil {
		return n, err
	}
	n += int64(n1)

	// generations
	nn, newStartXRef, err := v1.writeTo(w, n)
	if err != nil {
		return n, err
	}
	n += nn
	startXRef := newStartXRef

	// footer (3.4.4 File Trailer)
	n1, err = w.Write([]byte("startxref\n" + strconv.FormatInt(startXRef, 10) + "\n" + "%%EOF\n"))
	if err != nil {
		return n, err
	}
	n += int64(n1)

	// return without error
	return n, nil
}
