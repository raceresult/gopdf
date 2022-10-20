package types

import (
	"errors"
)

// PDF Reference 1.4, Table 5.17 Entries in a Type 0 font dictionary

type Type0Font struct {
	// (Required) The type of PDF object that this dictionary describes; must be
	// Font for a font dictionary.
	// Type

	// (Required) The type of font; must be Type0 for a Type 0 font.
	// Subtype

	// (Required) The PostScript name of the font. In principle, this is an arbitrary
	// name, since there is no font program associated directly with a Type 0 font
	// dictionary. The conventions described here ensure maximum compatibility
	// with existing Acrobat products.
	// If the descendant is a Type 0 CIDFont, this name should be the concatenation
	// of the CIDFont’s BaseFont name, a hyphen, and the CMap name given in the
	// Encoding entry (or the CMapName entry in the CMap program itself). If the
	// descendant is a Type 2 CIDFont, this name should be the same as the
	// CIDFont’s BaseFont name.
	BaseFont Name

	// (Required) The name of a predefined CMap, or a stream containing a CMap
	// stream program, that maps character codes to font numbers and CIDs. If the descen-
	// dant is a Type 2 CIDFont whose associated TrueType font program is not em-
	// bedded in the PDF file, the Encoding entry must be a predefined CMap name
	// (see “Glyph Selection in CIDFonts” on page 339).
	Encoding Name

	// (Required) An array specifying one or more fonts or CIDFonts that are
	// descendants of this composite font. This array is indexed by the font number
	// that is obtained by mapping a character code through the CMap specified in
	// the Encoding entry.
	// Note: In all PDF versions up to and including PDF 1.4, DescendantFonts must
	// be a one-element array containing a CIDFont dictionary.
	DescendantFonts Array

	// (Optional) A stream containing a CMap file that maps character codes to
	// Unicode values (see Section 5.9, “ToUnicode CMaps”).
	ToUnicode Reference
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

func (q *Type0Font) Read(dict Dictionary) error {
	// Type
	v, ok := dict["Type"]
	if !ok {
		return errors.New("font missing Type")
	}
	dtype, ok := v.(Name)
	if !ok {
		return errors.New("font field Type invalid")
	}
	if dtype != "Font" {
		return errors.New("unexpected value in font field Type")
	}

	// Subtype
	v, ok = dict["Subtype"]
	if !ok {
		return errors.New("font field Subtype missing")
	}
	vt, ok := v.(Name)
	if !ok {
		return errors.New("font field Subtype invalid")
	}
	if vt != "Type0" {
		return errors.New("font field Subtype invalid")
	}

	// BaseFont
	v, ok = dict["BaseFont"]
	if !ok {
		return errors.New("font field BaseFont missing")
	}
	q.BaseFont, ok = v.(Name)
	if !ok {
		return errors.New("font field BaseFont invalid")
	}

	// Encoding
	v, ok = dict["Encoding"]
	if !ok {
		return errors.New("font field Encoding missing")
	}
	q.Encoding, ok = v.(Name)
	if !ok {
		return errors.New("font field Encoding invalid")
	}

	// DescendantFonts
	v, ok = dict["DescendantFonts"]
	if !ok {
		return errors.New("font field DescendantFonts missing")
	}
	q.DescendantFonts, ok = v.(Array)
	if !ok {
		return errors.New("font field DescendantFonts invalid")
	}

	// ToUnicode
	v, ok = dict["ToUnicode"]
	if ok {
		q.ToUnicode, ok = v.(Reference)
		if !ok {
			return errors.New("font field ToUnicode invalid")
		}
	}

	// return without error
	return nil
}

func (q Type0Font) Copy(copyRef func(reference Reference) Reference) Object {
	return Type0Font{
		BaseFont:        q.BaseFont.Copy(copyRef).(Name),
		Encoding:        q.Encoding.Copy(copyRef).(Name),
		DescendantFonts: q.DescendantFonts.Copy(copyRef).(Array),
		ToUnicode:       q.ToUnicode.Copy(copyRef).(Reference),
	}
}
