package builder

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
)

// Builder is the main object to build a PDF file
type Builder struct {
	Info            types.InformationDictionary
	ID              [2]string
	CompressStreams bool
	Version         float64

	file  *pdf.File
	pages []*Page
}

// NewBuilder creates a new Builder object
func NewBuilder() *Builder {
	return &Builder{
		file:            pdf.NewFile(),
		CompressStreams: true,
		Version:         pdffile.DefaultVersion,
	}
}

// Build builds the PDF document and returns the file as byte slice
func (q *Builder) Build() ([]byte, error) {
	q.file.Version = q.Version
	q.file.Info = q.Info
	q.file.ID = q.ID
	q.file.CompressStreams = q.CompressStreams

	for _, p := range q.pages {
		p.build(q)
	}
	return q.file.Write()
}

// NewPage adds a new page to the pdf
func (q *Builder) NewPage(size PageSize) *Page {
	p := NewPage(size)
	q.pages = append(q.pages, p)
	return p
}
