package pdftext

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/pdftext/arabic"
)

func StringModifications(s string, font pdf.FontHandler) string {
	s = arabic.Shape(s, font)
	s = reverseRTLString(s)
	return s
}
