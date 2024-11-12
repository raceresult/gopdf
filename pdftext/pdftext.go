package pdftext

import (
	"github.com/raceresult/gopdf/pdftext/arabic"
)

func StringModifications(s string) string {
	s = arabic.Shape(s)
	s = reverseRTLString(s)
	return s
}
