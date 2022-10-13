package types

// PDF Reference 1.4, Table 3.21 Entries in a resource dictionary

type ResourceDictionary struct {
	// (Optional) A dictionary mapping resource names to graphics state parameter
	// dictionaries (see Section 4.3.4, “Graphics State Parameter Dictionaries”).
	ExtGState Dictionary

	// (Optional) A dictionary mapping each resource name to either the name of a
	// device-dependent color space or an array describing a color space (see Sec-
	// tion 4.5, “Color Spaces”).
	ColorSpace Dictionary

	// (Optional) A dictionary mapping resource names to pattern objects (see Sec-
	// tion 4.6, “Patterns”).
	Pattern Dictionary

	// (Optional; PDF 1.3) A dictionary mapping resource names to shading dic-
	// tionaries (see “Shading Dictionaries” on page 233).
	Shading Dictionary

	// (Optional) A dictionary mapping resource names to external objects (see Sec-
	// tion 4.7, “External Objects”)
	XObject Dictionary

	// (Optional) A dictionary mapping resource names to font dictionaries (see
	// Chapter 5).
	Font Dictionary

	// (Optional) An array of predefined procedure set names (see Section 9.1,
	// “Procedure Sets”).
	ProcSet Array

	// (Optional; PDF 1.2) A dictionary mapping resource names to property list
	// dictionaries for marked content (see Section 9.5.1, “Property Lists”).
	Properties Dictionary
}

func (q ResourceDictionary) ToRawBytes() []byte {
	d := Dictionary{}

	if q.ExtGState != nil {
		d["ExtGState"] = q.ExtGState
	}
	if q.ColorSpace != nil {
		d["ColorSpace"] = q.ColorSpace
	}
	if q.Pattern != nil {
		d["Pattern"] = q.Pattern
	}
	if q.Shading != nil {
		d["Shading"] = q.Shading
	}
	if q.XObject != nil {
		d["XObject"] = q.XObject
	}
	if q.Font != nil {
		d["Font"] = q.Font
	}
	if len(q.ProcSet) != 0 {
		d["ProcSet"] = q.ProcSet
	}
	if q.Properties != nil {
		d["Properties"] = q.Properties
	}

	return d.ToRawBytes()
}
