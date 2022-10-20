package types

// additional entries according to table 5.23

type StreamDictionaryFont struct {
	Stream StreamObject

	// (Required for Type 1 and TrueType fonts) The length in bytes of the clear-text portion
	// of the Type 1 font program (see below), or the entire TrueType font program, after it
	// has been decoded using the filters specified by the stream’s Filter entry, if any.
	Length1 Int

	// (Required for Type 1 fonts) The length in bytes of the encrypted portion of the Type 1
	// font program (see below) after it has been decoded using the filters specified by the
	// stream’s Filter entry.
	Length2 Int

	// (Required for Type 1 fonts) The length in bytes of the fixed-content portion of the
	// Type 1 font program (see below), after it has been decoded using the filters specified
	// by the stream’s Filter entry. If Length3 is 0, it indicates that the 512 zeros and clearto-
	// mark have not been included in the FontFile font program and must be added.
	Length3 Int

	// (Required if referenced from FontFile3; PDF 1.2) A name specifying the format of the
	// embedded font program. The name must be Type1C for Type 1 compact fonts or CID-
	// FontType0C for Type 0 compact CIDFonts. When additional font formats are added
	// to PDF, more values will be defined for Subtype.
	Subtype Name

	// (Optional; PDF 1.4) A metadata stream containing metadata for the embedded font
	// program (see Section 9.2.2, “Metadata Streams”).
	Metadata Object
}

func (q StreamDictionaryFont) ToRawBytes() []byte {
	d := q.Stream.Dictionary.createDict()
	if q.Length1 != 0 {
		d["Length1"] = q.Length1
	}
	if q.Length2 != 0 {
		d["Length2"] = q.Length2
	}
	if q.Length3 != 0 {
		d["Length3"] = q.Length3
	}
	if q.Subtype != "" {
		d["Subtype"] = q.Subtype
	}
	if q.Metadata != nil {
		d["Metadata"] = q.Metadata
	}
	return d.ToRawBytes()
}

func (q StreamDictionaryFont) Copy(copyRef func(reference Reference) Reference) Object {
	return StreamDictionaryFont{
		Stream:   q.Stream.Copy(copyRef).(StreamObject),
		Length1:  q.Length1.Copy(copyRef).(Int),
		Length2:  q.Length2.Copy(copyRef).(Int),
		Length3:  q.Length3.Copy(copyRef).(Int),
		Subtype:  q.Subtype.Copy(copyRef).(Name),
		Metadata: Copy(q.Metadata, copyRef),
	}
}
