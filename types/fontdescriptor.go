package types

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
