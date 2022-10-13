package types

import "golang.org/x/text/encoding/charmap"

// PDF Reference 1.4, Table D.1 Latin-text encodings

type Encoding Name

const (
	EncodingStandard  Encoding = "StandardEncoding"
	EncodingMacRoman  Encoding = "MacRomanEncoding"
	EncodingWinAnsi   Encoding = "WinAnsiEncoding"
	EncodingPDFDoc    Encoding = "PDFDocEncoding"
	EncodingMacExpert Encoding = "MacExpertEncoding"
)

func (q Encoding) ToRawBytes() []byte {
	return Name(q).ToRawBytes()
}

func (q Encoding) Encode(s string) string {
	switch q {
	case EncodingWinAnsi:
		sn, _ := charmap.Windows1252.NewEncoder().String(s)
		return sn
	case EncodingMacExpert:
		// todo
		return s
	case EncodingMacRoman:
		x, _ := charmap.Macintosh.NewEncoder().String(s)
		return x
	case EncodingPDFDoc:
		// todo
		return s
	default:
		return s
	}
}
