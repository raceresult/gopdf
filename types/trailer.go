package types

import (
	"bytes"
	"errors"
)

// PDF Reference 1.4, Table 3.12 Entries in the file trailer dictionary

type Trailer struct {
	// (Required) The total number of entries in the file’s cross-reference table, as defined
	// by the combination of the original section and all update sections. Equivalently, this
	// value is 1 greater than the highest object number used in the file.
	Size int

	// (Present only if the file has more than one cross-reference section) The byte offset from
	// the beginning of the file to the beginning of the previous cross-reference section
	Prev Int // todo: unclear: (Present only if the pdffile has more than one cross-reference section) The byte offset from the beginning of the pdffile to the beginning of the previous cross-reference section.

	// (Required; must be an indirect reference) The catalog dictionary for the PDF docu-
	// ment contained in the file (see Section 3.6.1, “Document Catalog”).
	Root Reference

	// (Required if document is encrypted; PDF 1.1) The document’s encryption dictionary
	// (see Section 3.5, “Encryption”).
	Encrypt Object

	// (Optional; must be an indirect reference) The document’s information dictionary
	// (see Section 9.2.1, “Document Information Dictionary”).
	Info Reference

	// (Optional; PDF 1.1) An array of two strings constituting a file identifier (see Section
	// 9.3, “File Identifiers”) for the file.
	ID [2]String
}

func (q *Trailer) ToRawBytes() []byte {
	var sb bytes.Buffer
	sb.WriteString("trailer\n")

	d := Dictionary{
		"Size": Int(q.Size),
		"Root": q.Root,
	}
	if q.Prev != 0 {
		d[Name("Prev")] = q.Prev
	}
	if q.Encrypt != nil {
		d[Name("Encrypt")] = q.Encrypt
	}
	if q.Info.Number > 0 {
		d[Name("Info")] = q.Info
	}
	if q.ID[0] != "" || q.ID[1] != "" {
		d[Name("ID")] = Array{q.ID[0], q.ID[1]}
	}

	sb.Write(d.ToRawBytes())
	return sb.Bytes()
}

func (q *Trailer) Read(dict Dictionary, file Resolver) error {
	// Size
	v, ok := dict.GetValue("Size", file)
	if !ok {
		return errors.New("trailer missing Size")
	}
	size, ok := v.(Int)
	if !ok {
		return errors.New("trailer field Size invalid")
	}
	q.Size = int(size)

	// Prev
	v, ok = dict.GetValue("Prev", file)
	if ok {
		prev, ok := v.(Int)
		if !ok {
			return errors.New("trailer field Prev invalid")
		}
		q.Prev = prev
	}

	// Root
	v, ok = dict["Root"]
	if !ok {
		return errors.New("trailer missing Root")
	}
	root, ok := v.(Reference)
	if !ok {
		return errors.New("trailer field Size invalid")
	}
	q.Root = root

	// Encrypt
	v, ok = dict.GetValue("Encrypt", file)
	if ok {
		q.Encrypt = v
	}

	// Info
	v, ok = dict["Info"]
	if ok {
		info, ok := v.(Reference)
		if !ok {
			return errors.New("trailer field Info invalid")
		}
		q.Info = info
	}

	// ID
	v, ok = dict.GetValue("ID", file)
	if ok {
		a, ok := v.(Array)
		if !ok {
			return errors.New("trailer field ID invalid")
		}
		if len(a) != 2 {
			return errors.New("trailer field ID invalid")
		}
		v1, ok := a[0].(String)
		if !ok {
			return errors.New("trailer field ID invalid")
		}
		v2, ok := a[0].(String)
		if !ok {
			return errors.New("trailer field ID invalid")
		}
		q.ID = [2]String{v1, v2}
	}

	// return without error
	return nil
}
