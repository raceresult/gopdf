package extractor

import "github.com/raceresult/gopdf/pdffile"

type Extractor struct {
	file *pdffile.File
}

func NewExtractor(f *pdffile.File) *Extractor {
	return &Extractor{
		file: f,
	}
}
