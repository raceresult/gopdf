package pdffile

import (
	"errors"
	"io"
	"strconv"
)

// RawFile represents a pdf file with unparsed raw objects
type RawFile struct {
	Version float64
	updates []Update
}

// NewRawFile creates a new RawFile object
func NewRawFile(version float64) *RawFile {
	return &RawFile{
		Version: version,
	}
}

// AddUpdate adds a new update to the raw file
func (q *RawFile) AddUpdate(u ...Update) {
	q.updates = append(q.updates, u...)
}

// WriteTo writes the parsed to the given writer
func (q *RawFile) WriteTo(w io.Writer) (int64, error) {
	if q.Version == 0 {
		return 0, errors.New("no pdf version set")
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
	var startXRef int64
	for _, u := range q.updates {
		nn, newStartXRef, err := u.writeTo(w, n)
		if err != nil {
			return n, err
		}
		n += nn
		startXRef = newStartXRef
	}

	// footer (3.4.4 File Trailer)
	n1, err = w.Write([]byte("startxref\n" + strconv.FormatInt(startXRef, 10) + "\n" + "%%EOF\n"))
	if err != nil {
		return n, err
	}
	n += int64(n1)

	// return without error
	return n, nil
}
