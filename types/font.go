package types

import (
	"errors"
)

// PDF Reference 1.4, Table 5.8 Entries in a font dictionary

type Font struct {
	// (Required) The type of PDF object that this dictionary describes; must be
	// Font for a font dictionary.
	// Type

	// (Required) The type of font; must be Type1 for a Type 1 font.
	Subtype FontSubType

	// (Required in PDF 1.0; optional otherwise) The name by which this font is ref-
	// erenced in the Font subdictionary of the current resource dictionary.
	// Note: This entry is obsolescent and its use is no longer recommended. (See
	// implementation note 42 in Appendix H.)
	Name Name

	// (Required) The PostScript name of the font. For Type 1 fonts, this is usually
	// the value of the FontName entry in the font program; for more information,
	// see Section 5.2 of the PostScript Language Reference, Third Edition. The Post-
	// Script name of the font can be used to find the font’s definition in the viewer
	// application or its environment. It is also the name that will be used when
	// printing to a PostScript output device.
	BaseFont Name

	// (Required except for the standard 14 fonts) The first character code defined in
	// the font’s Widths array.
	FirstChar Int

	// (Required except for the standard 14 fonts) The last character code defined in
	// the font’s Widths array.
	LastChar Int

	// (Required except for the standard 14 fonts; indirect reference preferred) An array
	// of (LastChar − FirstChar + 1) widths, each element being the glyph width for
	// the character whose code is FirstChar plus the array index. For character
	// codes outside the range FirstChar to LastChar, the value of MissingWidth from
	// the FontDescriptor entry for this font is used. The glyph widths are measured
	// in units in which 1000 units corresponds to 1 unit in text space. These widths
	// must be consistent with the actual widths given in the font program itself.
	// (See implementation note 43 in Appendix H.) For more information on
	// glyph widths and other glyph metrics, see Section 5.1.3, “Glyph Positioning
	// and Metrics.”
	Widths Array

	// (Required except for the standard 14 fonts; must be an indirect reference) A font
	//descriptor describing the font’s metrics other than its glyph widths (see Sec-
	//tion 5.7, “Font Descriptors”).
	FontDescriptor Reference

	// (Optional) A specification of the font’s character encoding, if different from
	// its built-in encoding. The value of Encoding may be either the name of a pre-
	// defined encoding (MacRomanEncoding, MacExpertEncoding, or WinAnsi-
	// Encoding, as described in Appendix D) or an encoding dictionary that
	// specifies differences from the font’s built-in encoding or from a specified pre-
	// defined encoding (see Section 5.5.5, “Character Encoding”).
	Encoding Object

	// (Optional; PDF 1.2) A stream containing a CMap file that maps character
	// codes to Unicode values (see Section 5.9, “ToUnicode CMaps”).
	ToUnicode Object
}

func (q Font) ToRawBytes() []byte {
	d := Dictionary{
		"Type":     Name("Font"),
		"Subtype":  q.Subtype,
		"BaseFont": q.BaseFont,
	}
	if q.FirstChar != 0 || q.LastChar != 0 {
		d["FirstChar"] = q.FirstChar
		d["LastChar"] = q.LastChar
	}
	if len(q.Widths) != 0 {
		d["Widths"] = q.Widths
	}
	if q.FontDescriptor.Number > 0 {
		d["FontDescriptor"] = q.FontDescriptor
	}
	if q.Name != "" {
		d["Name"] = q.Name
	}
	if q.Encoding != nil {
		d["Encoding"] = q.Encoding
	}
	if q.ToUnicode != nil {
		d["ToUnicode"] = q.ToUnicode
	}
	return d.ToRawBytes()
}

func (q *Font) Read(dict Dictionary) error {
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
	q.Subtype, ok = v.(FontSubType)
	if !ok {
		n, ok := v.(Name)
		if !ok {
			return errors.New("font field Pages invalid")
		}
		q.Subtype = FontSubType(n)
	}

	// Name
	v, ok = dict["Name"]
	if ok {
		q.Name, ok = v.(Name)
		if !ok {
			return errors.New("font field Name invalid")
		}
	}

	// BaseFont
	v, ok = dict["BaseFont"]
	if !ok {
		return errors.New("font field BaseFont missing")
	}
	q.BaseFont, ok = v.(Name)
	if !ok {
		return errors.New("font field Pages invalid")
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
	v, ok = dict["Widths"]
	if !ok {
		return errors.New("font field Widths missing")
	}
	q.Widths, ok = v.(Array)
	if !ok {
		return errors.New("font field Widths invalid")
	}

	// FontDescriptor
	v, ok = dict["FontDescriptor"]
	if !ok {
		return errors.New("font field FontDescriptor missing")
	}
	q.FontDescriptor, ok = v.(Reference)
	if !ok {
		return errors.New("font field FontDescriptor invalid")
	}

	// Encoding
	v, ok = dict["Encoding"]
	if ok {
		q.Encoding = v
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

func (q Font) Copy(copyRef func(reference Reference) Reference) Object {
	return Font{
		Subtype:        q.Subtype.Copy(copyRef).(FontSubType),
		Name:           q.Name.Copy(copyRef).(Name),
		BaseFont:       q.BaseFont.Copy(copyRef).(Name),
		FirstChar:      q.FirstChar.Copy(copyRef).(Int),
		LastChar:       q.LastChar.Copy(copyRef).(Int),
		Widths:         q.Widths.Copy(copyRef).(Array),
		FontDescriptor: q.FontDescriptor.Copy(copyRef).(Reference),
		Encoding:       Copy(q.Encoding, copyRef),
		ToUnicode:      Copy(q.ToUnicode, copyRef),
	}
}
