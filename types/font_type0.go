package types

// PDF Reference 1.4, Table 5.17 Entries in a Type 0 font dictionary

type Type0Font struct {
	// Type
	// Subtype
	BaseFont        Name      // Required
	Encoding        Name      // Required
	DescendantFonts Array     // Required
	ToUnicode       Reference // Optional
}

func (q Type0Font) ToRawBytes() []byte {
	d := Dictionary{
		"Type":            Name("Font"),
		"Subtype":         Name("Type0"),
		"BaseFont":        q.BaseFont,
		"Encoding":        q.Encoding,
		"DescendantFonts": q.DescendantFonts,
	}
	if q.ToUnicode.Number != 0 {
		d["ToUnicode"] = q.ToUnicode
	}
	return d.ToRawBytes()
}
