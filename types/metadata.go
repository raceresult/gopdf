package types

import (
	"bytes"
	"errors"
)

type MetaData struct {
	Dictionary StreamDictionary
	Stream     []byte
	Subtype    Name
}

func (q MetaData) ToRawBytes() []byte {
	d := q.Dictionary.createDict()
	d["Type"] = Name("Metadata")
	if q.Subtype != "" {
		d["Subtype"] = q.Subtype
	}

	sb := bytes.Buffer{}
	sb.Write(d.ToRawBytes())

	sb.WriteString("stream\n")
	sb.Write(q.Stream)
	sb.WriteString("\n")
	sb.WriteString("endstream\n")

	return sb.Bytes()
}

func (q *MetaData) Read(dict Dictionary, file Resolver) error {
	// Type
	v, ok := dict.GetValue("Type", file)
	if !ok {
		return errors.New("Metadata missing Type")
	}
	dtype, ok := v.(Name)
	if !ok {
		return errors.New("Metadata field Type invalid")
	}
	if dtype != "Metadata" {
		return errors.New("unexpected value in Metadata field Type")
	}

	// Subtype
	v, ok = dict.GetValue("Subtype", file)
	vt, ok := v.(Name)
	if !ok {
		return errors.New("\n field SubType invalid")
	}
	q.Subtype = vt

	// return without error
	return nil
}

func (q MetaData) Copy(copyRef func(reference Reference) Reference) Object {
	return MetaData{
		Dictionary: q.Dictionary.Copy(copyRef).(StreamDictionary),
		Stream:     q.Stream,
		Subtype:    q.Subtype,
	}
}

func (q MetaData) Equal(obj Object) bool {
	a, ok := obj.(MetaData)
	if !ok {
		return false
	}
	if !Equal(q.Dictionary, a.Dictionary) {
		return false
	}
	if !bytes.Equal(q.Stream, a.Stream) {
		return false
	}
	if !Equal(q.Subtype, a.Subtype) {
		return false
	}
	return true
}
