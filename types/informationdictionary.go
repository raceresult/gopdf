package types

// PDF Reference 1.4, Table 9.2 Entries in the document information dictionary

type InformationDictionary struct {
	// (Optional; PDF 1.1) The document’s title.
	Title String

	// (Optional) The name of the person who created the document.
	Author String

	// (Optional; PDF 1.1) The subject of the document.
	Subject String

	// (Optional; PDF 1.1) Keywords associated with the document.
	Keywords String

	// (Optional) If the document was converted to PDF from another format, the
	// name of the application (for example, Adobe FrameMaker®) that created the
	// original document from which it was converted.
	Creator String

	// (Optional) If the document was converted to PDF from another format, the
	// name of the application (for example, Acrobat Distiller) that converted it to
	// PDF.
	Producer String

	// (Optional) The date and time the document was created, in human-readable
	// form (see Section 3.8.2, “Dates”).
	CreationDate Date

	// (Optional; PDF 1.1) The date and time the document was most recently
	// modified, in human-readable form (see Section 3.8.2, “Dates”).
	ModDate Date

	// (Optional; PDF 1.3) A name object indicating whether the document has
	// been modified to include trapping information (see Section 9.10.5, “Trap-
	// ping Support”):
	Trapped Name
}

func (q InformationDictionary) ToRawBytes() []byte {
	d := Dictionary{}
	if q.Title != "" {
		d["Title"] = q.Title
	}
	if q.Author != "" {
		d["Author"] = q.Author
	}
	if q.Subject != "" {
		d["Subject"] = q.Subject
	}
	if q.Keywords != "" {
		d["Keywords"] = q.Keywords
	}
	if q.Creator != "" {
		d["Creator"] = q.Creator
	}
	if q.Producer != "" {
		d["Producer"] = q.Producer
	}
	if !q.CreationDate.IsZero() {
		d["CreationDate"] = q.CreationDate
	}
	if !q.ModDate.IsZero() {
		d["ModDate"] = q.ModDate
	}
	if q.Trapped != "" {
		d["Trapped"] = q.Trapped
	}
	return d.ToRawBytes()
}
