package types

import "bytes"

type EmbeddedFile struct {
	Dictionary StreamDictionary
	Stream     []byte

	Params  Dictionary
	Subtype Name
}

func (q EmbeddedFile) ToRawBytes() []byte {
	d := q.Dictionary.createDict()
	d["Type"] = Name("EmbeddedFile")
	d["Subtype"] = q.Subtype
	d["Params"] = q.Params

	sb := bytes.Buffer{}
	sb.Write(d.ToRawBytes())

	sb.WriteString("stream\n")
	sb.Write(q.Stream)
	sb.WriteString("\n")
	sb.WriteString("endstream\n")

	return sb.Bytes()
}

func (q EmbeddedFile) Copy(copyRef func(reference Reference) Reference) Object {
	return EmbeddedFile{
		Dictionary: q.Dictionary.Copy(copyRef).(StreamDictionary),
		Stream:     q.Stream,
		Params:     q.Params.Copy(copyRef).(Dictionary),
		Subtype:    q.Subtype.Copy(copyRef).(Name),
	}
}

func (q EmbeddedFile) Equal(obj Object) bool {
	a, ok := obj.(EmbeddedFile)
	if !ok {
		return false
	}
	if !Equal(q.Dictionary, a.Dictionary) {
		return false
	}
	if !bytes.Equal(q.Stream, a.Stream) {
		return false
	}
	if !Equal(q.Params, a.Params) {
		return false
	}
	if !Equal(q.Subtype, a.Subtype) {
		return false
	}
	return true
}
