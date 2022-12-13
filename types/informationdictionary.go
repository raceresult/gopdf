package types

import "errors"

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

func (q *InformationDictionary) Read(dict Dictionary, file Resolver) error {
	v, ok := dict.GetValue("Title", file)
	if ok {
		q.Title, ok = v.(String)
		if !ok {
			return errors.New("information dictionary field Title invalid")
		}
	}

	v, ok = dict.GetValue("Author", file)
	if ok {
		q.Author, ok = v.(String)
		if !ok {
			return errors.New("information dictionary field Author invalid")
		}
	}

	v, ok = dict.GetValue("Subject", file)
	if ok {
		q.Subject, ok = v.(String)
		if !ok {
			return errors.New("information dictionary field Subject invalid")
		}
	}

	v, ok = dict.GetValue("Keywords", file)
	if ok {
		q.Keywords, ok = v.(String)
		if !ok {
			return errors.New("information dictionary field Keywords invalid")
		}
	}

	v, ok = dict.GetValue("Creator", file)
	if ok {
		q.Creator, ok = v.(String)
		if !ok {
			return errors.New("information dictionary field Creator invalid")
		}
	}

	v, ok = dict.GetValue("Producer", file)
	if ok {
		q.Producer, ok = v.(String)
		if !ok {
			return errors.New("information dictionary field Producer invalid")
		}
	}

	v, ok = dict.GetValue("CreationDate", file)
	if ok {
		q.CreationDate, ok = v.(Date)
		if !ok {
			return errors.New("information dictionary field CreationDate invalid")
		}
	}

	v, ok = dict.GetValue("ModDate", file)
	if ok {
		q.ModDate, ok = v.(Date)
		if !ok {
			return errors.New("information dictionary field ModDate invalid")
		}
	}

	v, ok = dict.GetValue("Trapped", file)
	if ok {
		q.Trapped, ok = v.(Name)
		if !ok {
			return errors.New("information dictionary field Trapped invalid")
		}
	}

	// return without error
	return nil
}

func (q InformationDictionary) Copy(copyRef func(reference Reference) Reference) Object {
	return InformationDictionary{
		Title:        q.Title.Copy(copyRef).(String),
		Author:       q.Author.Copy(copyRef).(String),
		Subject:      q.Subject.Copy(copyRef).(String),
		Keywords:     q.Keywords.Copy(copyRef).(String),
		Creator:      q.Creator.Copy(copyRef).(String),
		Producer:     q.Producer.Copy(copyRef).(String),
		CreationDate: q.CreationDate.Copy(copyRef).(Date),
		ModDate:      q.ModDate.Copy(copyRef).(Date),
		Trapped:      q.Trapped.Copy(copyRef).(Name),
	}
}

func (q InformationDictionary) Equal(obj Object) bool {
	a, ok := obj.(InformationDictionary)
	if !ok {
		return false
	}
	if !Equal(q.Title, a.Title) {
		return false
	}
	if !Equal(q.Author, a.Author) {
		return false
	}
	if !Equal(q.Subject, a.Subject) {
		return false
	}
	if !Equal(q.Keywords, a.Keywords) {
		return false
	}
	if !Equal(q.Creator, a.Creator) {
		return false
	}
	if !Equal(q.Producer, a.Producer) {
		return false
	}
	if !Equal(q.CreationDate, a.CreationDate) {
		return false
	}
	if !Equal(q.ModDate, a.ModDate) {
		return false
	}
	if !Equal(q.Trapped, a.Trapped) {
		return false
	}
	return true
}
