package types

import (
	"errors"
)

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
	CIDSystemInfo Object

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

func (q *CIDFont) Read(dict Dictionary) error {
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
	if vt != Name(FontSub_CIDFontType0) && vt != Name(FontSub_CIDFontType2) {
		return errors.New("font field Subtype invalid")
	}
	q.Subtype = FontSubType(vt)

	// BaseFont
	v, ok = dict["BaseFont"]
	if !ok {
		return errors.New("font field BaseFont missing")
	}
	q.BaseFont, ok = v.(Name)
	if !ok {
		return errors.New("font field BaseFont invalid")
	}

	// CIDSystemInfo
	q.CIDSystemInfo, ok = dict["CIDSystemInfo"]
	if !ok {
		return errors.New("font field CIDSystemInfo missing")
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

	// DW
	v, ok = dict["DW"]
	if ok {
		q.DW = v
	}

	// W
	v, ok = dict["W"]
	if ok {
		q.W = v
	}

	// DW2
	v, ok = dict["DW2"]
	if ok {
		q.DW2 = v
	}

	// W2
	v, ok = dict["W2"]
	if ok {
		q.W2 = v
	}

	// CIDToGIDMap
	v, ok = dict["CIDToGIDMap"]
	if ok {
		q.CIDToGIDMap = v
	}

	// return without error
	return nil
}

func (q CIDFont) Copy(copyRef func(reference Reference) Reference) Object {
	return CIDFont{
		Subtype:        q.Subtype.Copy(copyRef).(FontSubType),
		BaseFont:       q.BaseFont.Copy(copyRef).(Name),
		CIDSystemInfo:  Copy(q.CIDSystemInfo, copyRef),
		FontDescriptor: q.FontDescriptor.Copy(copyRef).(Reference),
		DW:             Copy(q.DW, copyRef),
		W:              Copy(q.W, copyRef),
		DW2:            Copy(q.DW2, copyRef),
		W2:             Copy(q.W2, copyRef),
		CIDToGIDMap:    Copy(q.CIDToGIDMap, copyRef),
	}
}

func (q CIDFont) Equal(obj Object) bool {
	a, ok := obj.(CIDFont)
	if !ok {
		return false
	}
	if !Equal(q.Subtype, a.Subtype) {
		return false
	}
	if !Equal(q.BaseFont, a.BaseFont) {
		return false
	}
	if !Equal(q.CIDSystemInfo, a.CIDSystemInfo) {
		return false
	}
	if !Equal(q.FontDescriptor, a.FontDescriptor) {
		return false
	}
	if !Equal(q.DW, a.DW) {
		return false
	}
	if !Equal(q.W, a.W) {
		return false
	}
	if !Equal(q.DW2, a.DW2) {
		return false
	}
	if !Equal(q.W2, a.W2) {
		return false
	}
	if !Equal(q.CIDToGIDMap, a.CIDToGIDMap) {
		return false
	}
	return true
}
