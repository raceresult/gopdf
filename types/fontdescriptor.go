package types

import "errors"

// PDF Reference 1.4, Table 5.18 Entries common to all font descriptors

type FontDescriptor struct {
	// Type name (Required) The type of PDF object that this dictionary describes; must be
	// FontDescriptor for a font descriptor.
	// Type // Required

	// (Required) The PostScript name of the font. This should be the same as the
	// value of BaseFont in the font or CIDFont dictionary that refers to this font
	// descriptor.
	FontName Name

	// (Required) A collection of flags defining various characteristics of the font
	// (see Section 5.7.1, “Font Descriptor Flags”).
	Flags Int

	// (Required) A rectangle (see Section 3.8.3, “Rectangles”), expressed in the
	// glyph coordinate system, specifying the font bounding box. This is the small-
	// est rectangle enclosing the shape that would result if all of the glyphs of the
	// font were placed with their origins coincident and then filled.
	FontBBox Rectangle

	// (Required) The angle, expressed in degrees counterclockwise from the verti-
	// cal, of the dominant vertical strokes of the font. (For example, the 9-o’clock
	// position is 90 degrees, and the 3-o’clock position is –90 degrees.) The value is
	// negative for fonts that slope to the right, as almost all italic fonts do.
	ItalicAngle Number

	// The maximum height above the baseline reached by glyphs in this
	// font, excluding the height of glyphs for accented characters.
	Ascent Number

	// The maximum depth below the baseline reached by glyphs in this
	// font. The value is a negative number.
	Descent Number

	// The desired spacing between baselines of consecutive lines of text.
	// Default value: 0.
	Leading Number

	// The vertical coordinate of the top of flat capital letters, measured
	// from the baseline.
	CapHeight Number

	// (Optional) The font’s x height: the vertical coordinate of the top of flat non-
	// ascending lowercase letters (like the letter x), measured from the baseline.
	// Default value: 0.
	XHeight Number

	// (Required) The thickness, measured horizontally, of the dominant vertical
	// stems of glyphs in the font.
	StemV Number

	// (Optional) The thickness, measured invertically, of the dominant horizontal
	// stems of glyphs in the font. Default value: 0.
	StemH Number

	// (Optional) The average width of glyphs in the font. Default value: 0.
	AvgWidth Number

	// (Optional) The maximum width of glyphs in the font. Default value: 0.
	MaxWidth Number

	// (Optional) The width to use for character codes whose widths are not speci-
	// fied in a font dictionary’s Widths array. This has a predictable effect only if all
	// such codes map to glyphs whose actual widths are the same as the Missing-
	// Width value. Default value: 0.
	MissingWidth Number

	// (Optional) A stream containing a Type 1 font program (see Section 5.8,
	// “Embedded Font Programs”).
	FontFile Reference

	// (Optional; PDF 1.1) A stream containing a TrueType font program (see Sec-
	// tion 5.8, “Embedded Font Programs”).
	FontFile2 Reference

	// (Optional; PDF 1.2) A stream containing a font program other than Type 1 or
	// TrueType. The format of the font program is specified by the Subtype entry
	// in the stream dictionary (see Section 5.8, “Embedded Font Programs,” and
	// implementation note 49 in Appendix H).
	// At most, only one of the FontFile, FontFile2, and FontFile3 entries may be
	// present.
	FontFile3 Reference

	// (Optional; meaningful only in Type 1 fonts; PDF 1.1) A string listing the char-
	// acter names defined in a font subset. The names in this string must be in PDF
	// syntax—that is, each name preceded by a slash (/). The names can appear in
	// any order. The name .notdef should be omitted; it is assumed to exist in the
	// font subset. If this entry is absent, the only indication of a font subset is the
	// subset tag in the FontName entry (see Section 5.5.3, “Font Subsets”).
	CharSet String
}

func (q FontDescriptor) ToRawBytes() []byte {
	d := Dictionary{
		"Type":        Name("FontDescriptor"),
		"FontName":    q.FontName,
		"Flags":       q.Flags,
		"FontBBox":    q.FontBBox,
		"ItalicAngle": q.ItalicAngle,
		"Ascent":      q.Ascent,
		"Descent":     q.Descent,
		"CapHeight":   q.CapHeight,
		"StemV":       q.StemV,
	}

	if q.Leading != 0 {
		d["Leading"] = q.Leading
	}
	if q.XHeight != 0 {
		d["XHeight"] = q.XHeight
	}
	if q.StemH != 0 {
		d["StemH"] = q.StemH
	}
	if q.AvgWidth != 0 {
		d["AvgWidth"] = q.AvgWidth
	}
	if q.MaxWidth != 0 {
		d["MaxWidth"] = q.MaxWidth
	}
	if q.MissingWidth != 0 {
		d["MissingWidth"] = q.MissingWidth
	}
	if q.FontFile.Number > 0 {
		d["FontFile"] = q.FontFile
	}
	if q.FontFile2.Number > 0 {
		d["FontFile2"] = q.FontFile2
	}
	if q.FontFile3.Number > 0 {
		d["FontFile3"] = q.FontFile3
	}
	if q.CharSet != "" {
		d["CharSet"] = q.CharSet
	}

	return d.ToRawBytes()
}

func (q *FontDescriptor) Read(dict Dictionary) error {
	// Type
	v, ok := dict["Type"]
	if !ok {
		return errors.New("fontdescriptor missing Type")
	}
	dtype, ok := v.(Name)
	if !ok {
		return errors.New("fontdescriptor field Type invalid")
	}
	if dtype != "FontDescriptor" {
		return errors.New("unexpected value in fontdescriptor field Type")
	}

	// FontName
	v, ok = dict["FontName"]
	if !ok {
		return errors.New("fontdescriptor missing field FontName")
	}
	q.FontName, ok = v.(Name)
	if !ok {
		return errors.New("fontdescriptor field FontName invalid")
	}

	// Flags
	v, ok = dict["Flags"]
	if !ok {
		return errors.New("fontdescriptor missing field Flags")
	}
	q.Flags, ok = v.(Int)
	if !ok {
		return errors.New("fontdescriptor field Flags invalid")
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

	// ItalicAngle
	v, ok = dict["ItalicAngle"]
	if !ok {
		return errors.New("fontdescriptor missing field ItalicAngle")
	}
	q.ItalicAngle, ok = v.(Number)
	if !ok {
		vi, ok := v.(Int)
		if !ok {
			return errors.New("fontdescriptor field ItalicAngle invalid")
		}
		q.ItalicAngle = Number(vi)
	}

	// Ascent
	v, ok = dict["Ascent"]
	if !ok {
		return errors.New("fontdescriptor missing field Ascent")
	}
	q.Ascent, ok = v.(Number)
	if !ok {
		vi, ok := v.(Int)
		if !ok {
			return errors.New("fontdescriptor field Ascent invalid")
		}
		q.Ascent = Number(vi)
	}

	// Descent
	v, ok = dict["Descent"]
	if !ok {
		return errors.New("fontdescriptor missing field Descent")
	}
	q.Descent, ok = v.(Number)
	if !ok {
		vi, ok := v.(Int)
		if !ok {
			return errors.New("fontdescriptor field Descent invalid")
		}
		q.Descent = Number(vi)
	}

	// Leading
	v, ok = dict["Leading"]
	if ok {
		q.Leading, ok = v.(Number)
		if !ok {
			vi, ok := v.(Int)
			if !ok {
				return errors.New("fontdescriptor field Leading invalid")
			}
			q.Leading = Number(vi)
		}
	}

	// CapHeight
	v, ok = dict["CapHeight"]
	if !ok {
		return errors.New("fontdescriptor missing field CapHeight")
	}
	q.CapHeight, ok = v.(Number)
	if !ok {
		vi, ok := v.(Int)
		if !ok {
			return errors.New("fontdescriptor field CapHeight invalid")
		}
		q.CapHeight = Number(vi)
	}

	// XHeight
	v, ok = dict["XHeight"]
	if ok {
		q.XHeight, ok = v.(Number)
		if !ok {
			vi, ok := v.(Int)
			if !ok {
				return errors.New("fontdescriptor field XHeight invalid")
			}
			q.XHeight = Number(vi)
		}
	}

	// StemV
	v, ok = dict["StemV"]
	if !ok {
		return errors.New("fontdescriptor missing field StemV")
	}
	q.StemV, ok = v.(Number)
	if !ok {
		vi, ok := v.(Int)
		if !ok {
			return errors.New("fontdescriptor field StemV invalid")
		}
		q.StemV = Number(vi)
	}

	// StemH
	v, ok = dict["StemH"]
	if ok {
		q.StemH, ok = v.(Number)
		if !ok {
			vi, ok := v.(Int)
			if !ok {
				return errors.New("fontdescriptor field StemH invalid")
			}
			q.StemH = Number(vi)
		}
	}

	// AvgWidth
	v, ok = dict["AvgWidth"]
	if ok {
		q.AvgWidth, ok = v.(Number)
		if !ok {
			vi, ok := v.(Int)
			if !ok {
				return errors.New("fontdescriptor field AvgWidth invalid")
			}
			q.AvgWidth = Number(vi)
		}
	}

	// MaxWidth
	v, ok = dict["MaxWidth"]
	if ok {
		q.MaxWidth, ok = v.(Number)
		if !ok {
			vi, ok := v.(Int)
			if !ok {
				return errors.New("fontdescriptor field MaxWidth invalid")
			}
			q.MaxWidth = Number(vi)
		}
	}

	// MissingWidth
	v, ok = dict["MissingWidth"]
	if ok {
		q.MissingWidth, ok = v.(Number)
		if !ok {
			vi, ok := v.(Int)
			if !ok {
				return errors.New("fontdescriptor field MissingWidth invalid")
			}
			q.MissingWidth = Number(vi)
		}
	}

	// FontFile
	v, ok = dict["FontFile"]
	if ok {
		q.FontFile, ok = v.(Reference)
		if !ok {
			return errors.New("fontdescriptor field FontFile invalid")
		}
	}

	// FontFile2
	v, ok = dict["FontFile2"]
	if ok {
		q.FontFile2, ok = v.(Reference)
		if !ok {
			return errors.New("fontdescriptor field FontFile2 invalid")
		}
	}

	// FontFile3
	v, ok = dict["FontFile3"]
	if ok {
		q.FontFile3, ok = v.(Reference)
		if !ok {
			return errors.New("fontdescriptor field FontFile3 invalid")
		}
	}

	// CharSet
	v, ok = dict["CharSet"]
	if ok {
		q.CharSet, ok = v.(String)
		if !ok {
			return errors.New("fontdescriptor field CharSet invalid")
		}
	}

	// return without error
	return nil
}

func (q FontDescriptor) Copy(copyRef func(reference Reference) Reference) Object {
	return FontDescriptor{
		FontName:     q.FontName.Copy(copyRef).(Name),
		Flags:        q.Flags.Copy(copyRef).(Int),
		FontBBox:     q.FontBBox.Copy(copyRef).(Rectangle),
		ItalicAngle:  q.ItalicAngle.Copy(copyRef).(Number),
		Ascent:       q.Ascent.Copy(copyRef).(Number),
		Descent:      q.Descent.Copy(copyRef).(Number),
		Leading:      q.Leading.Copy(copyRef).(Number),
		CapHeight:    q.CapHeight.Copy(copyRef).(Number),
		XHeight:      q.XHeight.Copy(copyRef).(Number),
		StemV:        q.StemV.Copy(copyRef).(Number),
		StemH:        q.StemH.Copy(copyRef).(Number),
		AvgWidth:     q.AvgWidth.Copy(copyRef).(Number),
		MaxWidth:     q.MaxWidth.Copy(copyRef).(Number),
		MissingWidth: q.MissingWidth.Copy(copyRef).(Number),
		FontFile:     q.FontFile.Copy(copyRef).(Reference),
		FontFile2:    q.FontFile2.Copy(copyRef).(Reference),
		FontFile3:    q.FontFile3.Copy(copyRef).(Reference),
		CharSet:      q.CharSet.Copy(copyRef).(String),
	}
}

func (q FontDescriptor) Equal(obj Object) bool {
	a, ok := obj.(FontDescriptor)
	if !ok {
		return false
	}
	if !Equal(q.FontName, a.FontName) {
		return false
	}
	if !Equal(q.Flags, a.Flags) {
		return false
	}
	if !Equal(q.FontBBox, a.FontBBox) {
		return false
	}
	if !Equal(q.ItalicAngle, a.ItalicAngle) {
		return false
	}
	if !Equal(q.Ascent, a.Ascent) {
		return false
	}
	if !Equal(q.Descent, a.Descent) {
		return false
	}
	if !Equal(q.Leading, a.Leading) {
		return false
	}
	if !Equal(q.CapHeight, a.CapHeight) {
		return false
	}
	if !Equal(q.XHeight, a.XHeight) {
		return false
	}
	if !Equal(q.StemV, a.StemV) {
		return false
	}
	if !Equal(q.StemH, a.StemH) {
		return false
	}
	if !Equal(q.AvgWidth, a.AvgWidth) {
		return false
	}
	if !Equal(q.MaxWidth, a.MaxWidth) {
		return false
	}
	if !Equal(q.MissingWidth, a.MissingWidth) {
		return false
	}
	if !Equal(q.FontFile, a.FontFile) {
		return false
	}
	if !Equal(q.FontFile2, a.FontFile2) {
		return false
	}
	if !Equal(q.FontFile3, a.FontFile3) {
		return false
	}
	if !Equal(q.CharSet, a.CharSet) {
		return false
	}
	return true
}
