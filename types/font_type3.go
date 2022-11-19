package types

import (
	"errors"
)

// PDF Reference 1.4, Table 5.9 Entries in a Type 3 font dictionary

type Type3Font struct {
	// (Required) The type of PDF object that this dictionary describes; must be
	// Font for a font dictionary.
	// Type

	// (Required) The type of font; must be Type3 for a Type 3 font.
	// Subtype

	//	(Required in PDF 1.0; optional otherwise) See Table 5.8 on page 317.
	Name Name

	//	(Required) A rectangle (see Section 3.8.3, “Rectangles”), expressed in the
	//	glyph coordinate system, specifying the font bounding box. This is the small-
	//	est rectangle enclosing the shape that would result if all of the glyphs of the
	//	font were placed with their origins coincident and then filled.
	//	If all four elements of the rectangle are zero, no assumptions are made based
	//	on the font bounding box. If any element is nonzero, it is essential that the
	//	font bounding box be accurate; if any glyph’s marks fall outside this bound-
	//	ing box, incorrect behavior may result.
	FontBBox Rectangle

	//	(Required) An array of six numbers specifying the font matrix, mapping
	//	glyph space to text space (see Section 5.1.3, “Glyph Positioning and
	//	Metrics”). A common practice is to define glyphs in terms of a 1000-unit
	//	glyph coordinate system, in which case the font matrix is
	// [0.001 0 0 0.001 0 0].
	FontMatrix Array

	//	(Required) A dictionary in which each key is a character name and the value
	//	associated with that key is a content stream that constructs and paints the
	//	glyph for that character. The stream must include as its first operator either
	//	d0 or d1. This is followed by operators describing one or more graphics
	//	objects, which may include path, text, or image objects. See below for more
	//	details about Type 3 glyph descriptions.
	CharProcs Object

	//	(Required) An encoding dictionary whose Differences array specifies the
	//	complete character encoding for this font (see Section 5.5.5, “Character
	//	Encoding”; also see implementation note 46 in Appendix H).
	Encoding Object

	//	(Required) The first character code defined in the font’s Widths array.
	FirstChar Int

	//	(Required) The last character code defined in the font’s Widths array.
	LastChar Int

	//	(Required; indirect reference preferred) An array of (LastChar − FirstChar + 1)
	//	widths, each element being the glyph width for the character whose code is
	//	FirstChar plus the array index. For character codes outside the range FirstChar
	//	to LastChar, the width is 0. These widths are interpreted in glyph space as
	//	specified by FontMatrix (unlike the widths of a Type 1 font, which are in
	//	thousandths of a unit of text space).
	//	Note: If FontMatrix specifies a rotation, only the horizontal component of the
	//	transformed width is used. That is, the resulting displacement is always horizon-
	//	tal in text space, as is the case for all simple fonts.
	Widths Object

	//	(Optional but strongly recommended; PDF 1.2) A list of the named resources,
	//	such as fonts and images, required by the glyph descriptions in this font (see
	//	Section 3.7.2, “Resource Dictionaries”). If any glyph descriptions refer to
	//	named resources but this dictionary is absent, the names are looked up in the
	//	resource dictionary of the page on which the font is used. (See implementa-
	//	tion note 47 in Appendix H.)
	Resources Object

	//	(Optional; PDF 1.2) A stream containing a CMap file that maps character
	//	codes to Unicode values (see Section 5.9, “ToUnicode CMaps”).
	ToUnicode Reference
}

func (q Type3Font) ToRawBytes() []byte {
	d := Dictionary{
		"Type":       Name("Font"),
		"Subtype":    Name(FontSub_Type3),
		"FontBBox":   q.FontBBox,
		"FontMatrix": q.FontMatrix,
		"CharProcs":  q.CharProcs,
		"Encoding":   q.Encoding,
		"FirstChar":  q.FirstChar,
		"LastChar":   q.LastChar,
		"Widths":     q.Widths,
	}
	if q.Name != "" {
		d["Name"] = q.Name
	}
	if q.Resources != nil {
		d["Resources"] = q.Resources
	}
	if q.ToUnicode.Number != 0 {
		d["ToUnicode"] = q.ToUnicode
	}
	return d.ToRawBytes()
}

func (q *Type3Font) Read(dict Dictionary) error {
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
	vt, ok := v.(FontSubType)
	if !ok {
		vn, ok := v.(Name)
		if !ok {
			return errors.New("font field Subtype invalid")
		}
		vt = FontSubType(vn)
	}
	if vt != FontSub_Type3 {
		return errors.New("font field Subtype invalid")
	}

	// Name
	v, ok = dict["Name"]
	if ok {
		q.Name, ok = v.(Name)
		if !ok {
			return errors.New("font field Name invalid")
		}
	}

	// FontBBox
	v, ok = dict["FontBBox"]
	if !ok {
		return errors.New("font field FontBBox missing")
	}
	q.FontBBox, ok = v.(Rectangle)
	if !ok {
		va, ok := v.(Array)
		if !ok {
			return errors.New("font field FontBBox invalid")
		}
		if err := q.FontBBox.Read(va); err != nil {
			return err
		}
	}

	// FontMatrix
	v, ok = dict["FontMatrix"]
	if !ok {
		return errors.New("font field FontMatrix missing")
	}
	q.FontMatrix, ok = v.(Array)
	if !ok {
		return errors.New("font field FontMatrix invalid")
	}

	// CharProcs
	q.CharProcs, ok = dict["CharProcs"]
	if !ok {
		return errors.New("font field CharProcs missing")
	}

	// Encoding
	q.Encoding, ok = dict["Encoding"]
	if !ok {
		return errors.New("font field Encoding missing")
	}

	// FirstChar
	v, ok = dict["FirstChar"]
	if !ok {
		return errors.New("font field FirstChar missing")
	}
	q.FirstChar, ok = v.(Int)
	if !ok {
		return errors.New("font field FirstChar invalid")
	}

	// LastChar
	v, ok = dict["LastChar"]
	if !ok {
		return errors.New("font field LastChar missing")
	}
	q.LastChar, ok = v.(Int)
	if !ok {
		return errors.New("font field LastChar invalid")
	}

	// Widths
	q.Widths, ok = dict["Widths"]
	if !ok {
		return errors.New("font field Widths missing")
	}

	// Resources
	v, ok = dict["Resources"]
	if ok {
		q.Resources = v
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

func (q Type3Font) Copy(copyRef func(reference Reference) Reference) Object {
	return Type3Font{
		Name:       q.Name.Copy(copyRef).(Name),
		FontBBox:   q.FontBBox.Copy(copyRef).(Rectangle),
		FontMatrix: q.FontMatrix.Copy(copyRef).(Array),
		CharProcs:  Copy(q.CharProcs, copyRef),
		Encoding:   Copy(q.Encoding, copyRef),
		FirstChar:  q.FirstChar.Copy(copyRef).(Int),
		LastChar:   q.LastChar.Copy(copyRef).(Int),
		Widths:     Copy(q.Widths, copyRef),
		Resources:  Copy(q.Resources, copyRef),
		ToUnicode:  q.ToUnicode.Copy(copyRef).(Reference),
	}
}

func (q Type3Font) Equal(obj Object) bool {
	a, ok := obj.(Type3Font)
	if !ok {
		return false
	}
	if !Equal(q.Name, a.Name) {
		return false
	}
	if !Equal(q.FontBBox, a.FontBBox) {
		return false
	}
	if !Equal(q.FontMatrix, a.FontMatrix) {
		return false
	}
	if !Equal(q.CharProcs, a.CharProcs) {
		return false
	}
	if !Equal(q.Encoding, a.Encoding) {
		return false
	}
	if !Equal(q.FirstChar, a.FirstChar) {
		return false
	}
	if !Equal(q.LastChar, a.LastChar) {
		return false
	}
	if !Equal(q.Widths, a.Widths) {
		return false
	}
	if !Equal(q.Resources, a.Resources) {
		return false
	}
	if !Equal(q.ToUnicode, a.ToUnicode) {
		return false
	}
	return true
}
