package gopdf

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

// New creates a new Builder object
func New() *Builder {
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
	q.file.ID = [2]types.String{types.String(q.ID[0]), types.String(q.ID[1])}
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

// PageCount returns the number of pages already added
func (q *Builder) PageCount() int {
	return len(q.pages)
}

// NewImage adds a new image to the PDF file
func (q *Builder) NewImage(bts []byte) (*pdf.Image, error) {
	return q.file.NewImage(bts)
}

// NewCapturedPage adds a new captured page to the PDF file
func (q *Builder) NewCapturedPage(sourcePage types.Page, sourceFile *pdffile.File) (*CapturedPage, error) {
	cp, err := q.file.NewCapturedPage(sourcePage, sourceFile)
	if err != nil {
		return nil, err
	}
	return &CapturedPage{
		CapturedPage: cp,
	}, nil
}

// NewStandardFont adds a new standard font (expected to be available in all PDF consuming systems) to the pdf
func (q *Builder) NewStandardFont(name types.StandardFontName, encoding types.Encoding) (*pdf.StandardFont, error) {
	return q.file.NewStandardFont(name, encoding)
}

// NewTrueTypeFont adds a new TrueType font to the pdf
func (q *Builder) NewTrueTypeFont(ttf []byte, encoding types.Encoding, embed bool) (*pdf.TrueTypeFont, error) {
	return q.file.NewTrueTypeFont(ttf, encoding, embed)
}

// NewCompositeFont adds a font as composite font to the pdf, i.e. with Unicode support
func (q *Builder) NewCompositeFont(ttf []byte) (*pdf.CompositeFont, error) {
	return q.file.NewCompositeFont(ttf)
}
