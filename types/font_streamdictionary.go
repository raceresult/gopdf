package types

import "bytes"

// additional entries according to table 5.23

type StreamFont struct {
	Dictionary Dictionary
	Stream     []byte

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

func (q StreamFont) ToRawBytes() []byte {
	d := make(Dictionary)
	for k, v := range q.Dictionary {
		d[k] = v
	}
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

func (q StreamFont) Copy(copyRef func(reference Reference) Reference) Object {
	return StreamFont{
		Stream:     q.Stream,
		Dictionary: q.Dictionary.Copy(copyRef).(Dictionary),
		Length1:    q.Length1.Copy(copyRef).(Int),
		Length2:    q.Length2.Copy(copyRef).(Int),
		Length3:    q.Length3.Copy(copyRef).(Int),
		Subtype:    q.Subtype.Copy(copyRef).(Name),
		Metadata:   Copy(q.Metadata, copyRef),
	}
}

func (q StreamFont) Equal(obj Object) bool {
	a, ok := obj.(StreamFont)
	if !ok {
		return false
	}
	if !bytes.Equal(q.Stream, a.Stream) {
		return false
	}
	if !Equal(q.Dictionary, a.Dictionary) {
		return false
	}
	if !Equal(q.Length1, a.Length1) {
		return false
	}
	if !Equal(q.Length2, a.Length2) {
		return false
	}
	if !Equal(q.Length3, a.Length3) {
		return false
	}
	if !Equal(q.Subtype, a.Subtype) {
		return false
	}
	if !Equal(q.Metadata, a.Metadata) {
		return false
	}
	return true
}
