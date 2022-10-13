package builder

import "github.com/raceresult/gopdf/pdf"

// NewImage adds a new image to the PDF file
func (q *Builder) NewImage(bts []byte) (*pdf.Image, error) {
	return q.file.NewImage(bts)
}
