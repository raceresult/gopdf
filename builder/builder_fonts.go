package builder

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

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
