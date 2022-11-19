package types

import "errors"

// PDF Reference 1.4, Table 5.18 Entries common to all font descriptors

type FontDescriptor struct {
	// Type // Required
	FontName     Name      // Required
	Flags        Int       // Required
	FontBBox     Rectangle // Required
	ItalicAngle  Number    // Required
	Ascent       Number    // Required
	Descent      Number    // Required
	Leading      Number    // optional
	CapHeight    Number    // Required
	XHeight      Number    // optional
	StemV        Number    // Required
	StemH        Number    // optional
	AvgWidth     Number    // optional
	MaxWidth     Number    // optional
	MissingWidth Number    // optional
	FontFile     Reference // for Type 1 font
	FontFile2    Reference // for TrueType font
	FontFile3    Reference // for other font types
	CharSet      String    // optional
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
