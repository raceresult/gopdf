package types

// PDF Reference 1.4, Table 5.13 Entries in a CIDFont dictionary

type CIDFont struct {
	// (Required) The type of PDF object that this dictionary describes; must be
	// Font for a CIDFont dictionary.
	// Type

	// (Required) The type of CIDFont; CIDFontType0 or CIDFontType2.
	Subtype FontSubType

	// (Required) The PostScript name of the CIDFont. For Type 0 CIDFonts, this
	// is usually the value of the CIDFontName entry in the CIDFont program. For
	// Type 2 CIDFonts, it is derived the same way as for a simple TrueType font;
	// see Section 5.5.2, “TrueType Fonts.” In either case, the name can have a sub-
	// set prefix if appropriate; see Section 5.5.3, “Font Subsets.”
	BaseFont Name

	// (Required) A dictionary containing entries that define the character collec-
	// tion of the CIDFont. See Table 5.12 on page 337.
	CIDSystemInfo Reference

	// (Required; must be an indirect reference) A font descriptor describing the
	// CIDFont’s default metrics other than its glyph widths (see Section 5.7,
	// “Font Descriptors”).
	FontDescriptor Reference

	// (Optional) The default width for glyphs in the CIDFont (see “Glyph Met-
	// rics in CIDFonts” on page 340). Default value: 1000.
	DW Object

	// (Optional) A description of the widths for the glyphs in the CIDFont. The
	// array’s elements have a variable format that can specify individual widths
	// for consecutive CIDs or one width for a range of CIDs (see “Glyph Metrics
	// in CIDFonts” on page 340). Default value: none (the DW value is used for
	// all glyphs).
	W Object

	// (Optional; applies only to CIDFonts used for vertical writing) An array of two
	// numbers specifying the default metrics for vertical writing (see “Glyph
	// Metrics in CIDFonts” on page 340). Default value: [880 −1000].
	DW2 Object

	// (Optional; applies only to CIDFonts used for vertical writing) A description of
	// the metrics for vertical writing for the glyphs in the CIDFont (see “Glyph
	// Metrics in CIDFonts” on page 340). Default value: none (the DW2 value is
	// used for all glyphs).
	W2 Object

	// (Optional; Type 2 CIDFonts only) A specification of the mapping from CIDs
	// to glyph indices. If the value is a stream, the bytes in the stream contain the
	// mapping from CIDs to glyph indices: the glyph index for a particular CID
	// value c is a 2-byte value stored in bytes 2 ◊ c and 2 ◊ c + 1, where the first
	// byte is the high-order byte. If the value of CIDToGIDMap is a name, it must
	// be Identity, indicating that the mapping between CIDs and glyph indices is
	// the identity mapping. Default value: Identity.
	// This entry may appear only in a Type 2 CIDFont whose associated True-
	// Type font program is embedded in the PDF file (see the next section).
	CIDToGIDMap Object
}

func (q CIDFont) ToRawBytes() []byte {
	d := Dictionary{
		"Type":           Name("Font"),
		"Subtype":        q.Subtype,
		"BaseFont":       q.BaseFont,
		"CIDSystemInfo":  q.CIDSystemInfo,
		"FontDescriptor": q.FontDescriptor,
	}
	if q.DW != nil {
		d["DW"] = q.DW
	}
	if q.W != nil {
		d["W"] = q.W
	}
	if q.DW2 != nil {
		d["DW2"] = q.DW2
	}
	if q.W2 != nil {
		d["W2"] = q.W2
	}
	if q.CIDToGIDMap != nil {
		d["CIDToGIDMap"] = q.CIDToGIDMap
	}
	return d.ToRawBytes()
}
