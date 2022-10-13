package types

import (
	"bytes"
	"embed"

	"github.com/raceresult/gopdf/types/standardfont/afm"
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
	Encoding Encoding

	// (Optional; PDF 1.2) A stream containing a CMap file that maps character
	// codes to Unicode values (see Section 5.9, “ToUnicode CMaps”).
	ToUnicode *Stream
}

func (q Font) ToRawBytes() []byte {
	d := Dictionary{
		"Type":           Name("Font"),
		"Subtype":        q.Subtype,
		"BaseFont":       q.BaseFont,
		"FirstChar":      q.FirstChar,
		"LastChar":       q.LastChar,
		"Widths":         q.Widths,
		"FontDescriptor": q.FontDescriptor,
	}
	if q.Name != "" {
		d["Name"] = q.Name
	}
	if q.Encoding != "" {
		d["Encoding"] = q.Encoding
	}
	if q.ToUnicode != nil {
		d["ToUnicode"] = q.ToUnicode
	}
	return d.ToRawBytes()
}

type StandardFont struct {
	// Type
	// Subtype   FontSubType  // Required
	BaseFont StandardFontName // Required
	Encoding Encoding         // Name or dictionary, Optional
}

func (q StandardFont) ToRawBytes() []byte {
	d := Dictionary{
		"Type":     Name("Font"),
		"Subtype":  FontSub_Type1,
		"BaseFont": q.BaseFont,
	}
	if q.Encoding != "" {
		d["Encoding"] = q.Encoding
	}
	return d.ToRawBytes()
}

var (
	//go:embed standardfont/core14
	afmFiles embed.FS
)

// Metrics returns font metrics from the embedded afm files
func (q StandardFont) Metrics() (*afm.Font, error) {
	bts, err := afmFiles.ReadFile("standardfont/core14/" + string(q.BaseFont) + ".afm")
	if err != nil {
		return nil, err
	}
	m, err := afm.Parse(bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}
	return &m, nil
}
