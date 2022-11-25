package types

// PDF Reference 1.4, Table 3.32 Entries in a file specification dictionary

type FileSpecification struct {
	// (Required if an EF or RF entry is present; recommended always) The type of PDF object
	// that this dictionary describes; must be Filespec for a file specification dictionary.
	Type Name

	// (Optional) The name of the file system to be used to interpret this file specification. If
	// this entry is present, all other entries in the dictionary are interpreted by the desig-
	// nated file system. PDF defines only one standard file system, URL (see Section 3.10.4,
	// “URL Specifications”); a viewer application or plug-in extension can register a differ-
	// ent one (see Appendix E). Note that this entry is independent of the F, DOS, Mac, and
	// Unix entries.
	FS Name

	// (Required if the DOS, Mac, and Unix entries are all absent) A file specification string of
	// the form described in Section 3.10.1, “File Specification Strings,” or (if the file system
	// is URL) a uniform resource locator, as described in Section 3.10.4, “URL Specifica-
	// tions.”
	F String

	// (Optional) A file specification string (see Section 3.10.1, “File Specification Strings”)
	// representing a DOS file name.
	DOS String

	// (Optional) A file specification string (see Section 3.10.1, “File Specification Strings”)
	// representing a Mac OS file name.
	Mac String

	// (Optional) A file specification string (see Section 3.10.1, “File Specification Strings”)
	// representing a UNIX file name.
	Unix String

	// (Optional) An array of two strings constituting a file identifier (see Section 9.3, “File
	// Identifiers”) that is also included in the referenced file. The use of this entry improves
	// a viewer application’s chances of finding the intended file and allows it to warn the
	// user if the file has changed since the link was made.
	ID Array

	// (Optional; PDF 1.2) A flag indicating whether the file referenced by the file specifica-
	// tion is volatile (changes frequently with time). If the value is true, viewer applications
	// should never cache a copy of the file. For example, a movie annotation referencing a
	// URL to a live video camera could set this flag to true, notifying the application that it
	// should reacquire the movie each time it is played. Default value: false.
	V Object

	// (Required if RF is present; PDF 1.3) A dictionary containing a subset of the keys F,
	// DOS, Mac, and Unix, corresponding to the entries by those names in the file specifica-
	// tion dictionary. The value of each such key is an embedded file stream (see Section
	// 3.10.3, “Embedded File Streams”) containing the corresponding file. If this entry is
	// present, the Type entry is required and the file specification dictionary must be indi-
	// rectly referenced.
	EF Dictionary

	// (Optional; PDF 1.3) A dictionary with the same structure as the EF dictionary, which
	// must also be present. Each key in the RF dictionary must also be present in the EF dic-
	// tionary. Each value is a related files array (see “Related Files Arrays” on page 125)
	// identifying files that are related to the corresponding file in the EF dictionary. If this
	// entry is present, the Type entry is required and the file specification dictionary must
	// be indirectly referenced.
	RF Dictionary
}

func (q FileSpecification) ToRawBytes() []byte {
	d := Dictionary{
		"Type": q.Type,
	}

	if q.FS != "" {
		d["FS"] = q.FS
	}
	if q.F != "" {
		d["F"] = q.F
	}
	if q.DOS != "" {
		d["DOS"] = q.DOS
	}
	if q.Mac != "" {
		d["Mac"] = q.Mac
	}
	if q.Unix != "" {
		d["Unix"] = q.Unix
	}
	if len(q.ID) != 0 {
		d["ID"] = q.ID
	}
	if q.V != nil {
		d["V"] = q.V
	}
	if q.EF != nil {
		d["EF"] = q.EF
	}
	if q.RF != nil {
		d["RF"] = q.RF
	}
	return d.ToRawBytes()
}

func (q FileSpecification) Copy(copyRef func(reference Reference) Reference) Object {
	return FileSpecification{
		Type: q.Type.Copy(copyRef).(Name),
		FS:   q.FS.Copy(copyRef).(Name),
		F:    q.F.Copy(copyRef).(String),
		DOS:  q.DOS.Copy(copyRef).(String),
		Mac:  q.Mac.Copy(copyRef).(String),
		Unix: q.Unix.Copy(copyRef).(String),
		ID:   q.ID.Copy(copyRef).(Array),
		V:    Copy(q.V, copyRef),
		EF:   q.EF.Copy(copyRef).(Dictionary),
		RF:   q.RF.Copy(copyRef).(Dictionary),
	}
}

func (q FileSpecification) Equal(obj Object) bool {
	a, ok := obj.(FileSpecification)
	if !ok {
		return false
	}
	if !Equal(q.Type, a.Type) {
		return false
	}
	if !Equal(q.FS, a.FS) {
		return false
	}
	if !Equal(q.F, a.F) {
		return false
	}
	if !Equal(q.DOS, a.DOS) {
		return false
	}
	if !Equal(q.Mac, a.Mac) {
		return false
	}
	if !Equal(q.Unix, a.Unix) {
		return false
	}
	if !Equal(q.ID, a.ID) {
		return false
	}
	if !Equal(q.V, a.V) {
		return false
	}
	if !Equal(q.EF, a.EF) {
		return false
	}
	if !Equal(q.RF, a.RF) {
		return false
	}
	return true
}
