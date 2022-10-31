package types

// PDF Reference 1.4, Table 3.21 Entries in a resource dictionary

type ResourceDictionary struct {
	// (Optional) A dictionary mapping resource names to graphics state parameter
	// dictionaries (see Section 4.3.4, “Graphics State Parameter Dictionaries”).
	ExtGState Object

	// (Optional) A dictionary mapping each resource name to either the name of a
	// device-dependent color space or an array describing a color space (see Sec-
	// tion 4.5, “Color Spaces”).
	ColorSpace Object

	// (Optional) A dictionary mapping resource names to pattern objects (see Sec-
	// tion 4.6, “Patterns”).
	Pattern Object

	// (Optional; PDF 1.3) A dictionary mapping resource names to shading dic-
	// tionaries (see “Shading Dictionaries” on page 233).
	Shading Object

	// (Optional) A dictionary mapping resource names to external objects (see Sec-
	// tion 4.7, “External Objects”)
	XObject Object

	// (Optional) A dictionary mapping resource names to font dictionaries (see
	// Chapter 5).
	Font Object

	// (Optional) An array of predefined procedure set names (see Section 9.1,
	// “Procedure Sets”).
	ProcSet Object

	// (Optional; PDF 1.2) A dictionary mapping resource names to property list
	// dictionaries for marked content (see Section 9.5.1, “Property Lists”).
	Properties Object
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
	if q.ProcSet != nil {
		d["ProcSet"] = q.ProcSet
	}
	if q.Properties != nil {
		d["Properties"] = q.Properties
	}

	return d.ToRawBytes()
}

func (q *ResourceDictionary) Read(dict Dictionary) error {
	// ExtGState
	v, ok := dict["ExtGState"]
	if ok {
		q.ExtGState = v
	}

	// ColorSpace
	v, ok = dict["ColorSpace"]
	if ok {
		q.ColorSpace = v
	}

	// Pattern
	v, ok = dict["Pattern"]
	if ok {
		q.Pattern = v
	}

	// Shading
	v, ok = dict["Shading"]
	if ok {
		q.Shading = v
	}

	// XObject
	v, ok = dict["XObject"]
	if ok {
		q.XObject = v
	}

	// Font
	v, ok = dict["Font"]
	if ok {
		q.Font = v
	}

	// ProcSet
	v, ok = dict["ProcSet"]
	if ok {
		q.ProcSet = v
	}

	// Properties
	v, ok = dict["Properties"]
	if ok {
		q.Properties = v
	}

	// return without errors
	return nil
}

func (q ResourceDictionary) Copy(copyRef func(reference Reference) Reference) Object {
	return ResourceDictionary{
		ExtGState:  Copy(q.ExtGState, copyRef),
		ColorSpace: Copy(q.ColorSpace, copyRef),
		Pattern:    Copy(q.Pattern, copyRef),
		Shading:    Copy(q.Shading, copyRef),
		XObject:    Copy(q.XObject, copyRef),
		Font:       Copy(q.Font, copyRef),
		ProcSet:    Copy(q.ProcSet, copyRef),
		Properties: Copy(q.Properties, copyRef),
	}
}
