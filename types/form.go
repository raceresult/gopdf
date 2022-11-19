package types

import "bytes"

// PDF Reference 1.4, Table 4.41 Additional entries specific to a type 1 form dictionary

type Form struct {
	Dictionary StreamDictionary
	Stream     []byte

	// (Optional) The type of PDF object that this dictionary describes; if present,
	// must be XObject for a form XObject.
	// Type

	// (Required) The type of XObject that this dictionary describes; must be Form
	// for a form XObject.
	// Subtype Name

	// (Optional) A code identifying the type of form XObject that this dictionary
	// describes. The only valid value defined at the time of publication is 1. Default
	// value: 1.
	FormType Int

	// (Required in PDF 1.0; optional otherwise) The name by which this form
	// XObject is referenced in the XObject subdictionary of the current resource
	// dictionary (see Section 3.7.2, “Resource Dictionaries”).
	// Note: This entry is obsolescent and its use is no longer recommended. (See
	// implementation note 38 in Appendix H.)
	Name Name

	// (Required if PieceInfo is present; optional otherwise; PDF 1.3) The date and
	// time (see Section 3.8.2, “Dates”) when the form XObject’s contents were
	// most recently modified. If a page-piece dictionary (PieceInfo) is present, the
	// modification date is used to ascertain which of the application data diction-
	// aries it contains correspond to the current content of the form (see Section
	// 9.4, “Page-Piece Dictionaries”).
	LastModified Date

	// (Required) An array of four numbers in the form coordinate system (see
	// below), giving the coordinates of the left, bottom, right, and top edges,
	// respectively, of the form XObject’s bounding box. These boundaries are used
	// to clip the form XObject and to determine its size for caching.
	BBox Rectangle

	// (Optional) An array of six numbers specifying the form matrix, which maps
	// form space into user space (see Section 4.2.3, “Transformation Matrices”).
	// Default value: the identity matrix [1 0 0 1 0 0].
	Matrix Array

	// (Optional but strongly recommended; PDF 1.2) A dictionary specifying any
	// resources (such as fonts and images) required by the form XObject (see Sec-
	// tion 3.7, “Content Streams and Resources”).
	// In PDF 1.1 and earlier, all named resources used in the form XObject must be
	// included in the resource dictionary of each page object on which the form
	// XObject appears, whether or not they also appear in the resource dictionary
	// of the form XObject itself. It can be useful to specify these resources in the
	// form XObject’s own resource dictionary as well, in order to determine which
	// resources are used inside the form XObject. If a resource is included in both
	// dictionaries, it should have the same name in both locations.In PDF 1.2 and later versions, form XObjects can be independent of the
	// content streams in which they appear, and this is strongly recommended
	// although not required. In an independent form XObject, the resource dic-
	// tionary of the form XObject is required and contains all named resources
	// used by the form XObject. These resources are not “promoted” to the outer
	// content stream’s resource dictionary, although that stream’s resource diction-
	// ary will refer to the form XObject itself.
	Resources Object

	// (Optional; PDF 1.4) A group attributes dictionary indicating that the contents
	// of the form XObject are to be treated as a group and specifying the attributes
	// of that group (see Section 4.9.2, “Group XObjects”).
	// Note: If a Ref entry (see below) is present, the group attributes also apply to the
	// external page imported by that entry. This allows such an imported page to be
	// treated as a group without further modification.
	// Ref dictionary (Optional; PDF 1.4) A reference dictionary identifying a page to be imported
	// from another PDF file, and for which the form XObject serves as a proxy (see
	// Section 4.9.3, “Reference XObjects”).
	Group Object

	// (Optional; PDF 1.4) A metadata stream containing metadata for the form
	// XObject (see Section 9.2.2, “Metadata Streams”).
	Metadata Object

	// (Optional; PDF 1.3) A page-piece dictionary associated with the form
	// XObject (see Section 9.4, “Page-Piece Dictionaries”).
	PieceInfo Object

	// (Required if the form XObject is a structural content item; PDF 1.3) The integer
	// key of the form XObject’s entry in the structural parent tree (see “Finding
	// Structure Elements from Content Items” on page 600).
	StructParent Int

	// (Required if the form XObject contains marked-content sequences that are struc-
	// tural content items; PDF 1.3) The integer key of the form XObject’s entry in
	// the structural parent tree (see “Finding Structure Elements from Content
	// Items” on page 600).
	// Note: At most one of the entries StructParent or StructParents may be present. A
	// form XObject can be either a content item in its entirety or a container for
	// marked-content sequences that are content items, but not both.
	StructParents Int

	// (Optional; PDF 1.2) An OPI version dictionary for the form XObject (see
	//Section 9.10.6, “Open Prepress Interface (OPI)”).
	OPI Object
}

func (q Form) ToRawBytes() []byte {
	d := q.Dictionary.createDict()
	d["Type"] = Name("XObject")
	d["Subtype"] = Name("Form")
	d["BBox"] = q.BBox
	if q.FormType != 0 {
		d["FormType"] = q.FormType
	}
	if q.Name != "" {
		d["Name"] = q.Name
	}
	if !q.LastModified.IsZero() {
		d["LastModified"] = q.LastModified
	}
	if len(q.Matrix) != 0 {
		d["Matrix"] = q.Matrix
	}
	if q.Resources != nil {
		d["Resources"] = q.Resources
	}
	if q.Group != nil {
		d["Group"] = q.Group
	}
	if q.Metadata != nil {
		d["Metadata"] = q.Metadata
	}
	if q.PieceInfo != nil {
		d["PieceInfo"] = q.PieceInfo
	}
	if q.StructParent != 0 {
		d["StructParent"] = q.StructParent
	}
	if q.StructParents != 0 {
		d["StructParents"] = q.StructParents
	}
	if q.OPI != nil {
		d["OPI"] = q.OPI
	}

	sb := bytes.Buffer{}
	sb.Write(d.ToRawBytes())

	sb.WriteString("stream\n")
	sb.Write(q.Stream)
	sb.WriteString("\n")
	sb.WriteString("endstream\n")

	return sb.Bytes()
}

func (q Form) Copy(copyRef func(reference Reference) Reference) Object {
	return Form{
		Dictionary:    q.Dictionary.Copy(copyRef).(StreamDictionary),
		Stream:        q.Stream,
		FormType:      q.FormType.Copy(copyRef).(Int),
		Name:          q.Name.Copy(copyRef).(Name),
		LastModified:  q.LastModified.Copy(copyRef).(Date),
		BBox:          q.BBox.Copy(copyRef).(Rectangle),
		Matrix:        q.Matrix.Copy(copyRef).(Array),
		Resources:     Copy(q.Resources, copyRef),
		Group:         Copy(q.Group, copyRef),
		Metadata:      Copy(q.Metadata, copyRef),
		PieceInfo:     Copy(q.PieceInfo, copyRef),
		StructParent:  q.StructParent.Copy(copyRef).(Int),
		StructParents: q.StructParents.Copy(copyRef).(Int),
		OPI:           Copy(q.OPI, copyRef),
	}
}

func (q Form) Equal(obj Object) bool {
	a, ok := obj.(Form)
	if !ok {
		return false
	}
	if !Equal(q.Dictionary, a.Dictionary) {
		return false
	}
	if !bytes.Equal(q.Stream, a.Stream) {
		return false
	}
	if !Equal(q.FormType, a.FormType) {
		return false
	}
	if !Equal(q.Name, a.Name) {
		return false
	}
	if !Equal(q.LastModified, a.LastModified) {
		return false
	}
	if !Equal(q.BBox, a.BBox) {
		return false
	}
	if !Equal(q.Matrix, a.Matrix) {
		return false
	}
	if !Equal(q.Resources, a.Resources) {
		return false
	}
	if !Equal(q.Group, a.Group) {
		return false
	}
	if !Equal(q.Metadata, a.Metadata) {
		return false
	}
	if !Equal(q.PieceInfo, a.PieceInfo) {
		return false
	}
	if !Equal(q.StructParent, a.StructParent) {
		return false
	}
	if !Equal(q.StructParents, a.StructParents) {
		return false
	}
	if !Equal(q.OPI, a.OPI) {
		return false
	}
	return true
}
