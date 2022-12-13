package types

import "errors"

// PDF Reference 1.6, Table 3.14 Additional entries specific to an object stream dictionary

type ObjectStreamDictionary struct {
	StreamDictionary

	// (Required) The type of PDF object that this dictionary describes; must be ObjStm
	// for an object stream.
	// Type name

	// (Required) The number of compressed objects in the stream.
	N Int

	// (Required) The byte offset (in the decoded stream) of the first compressed object.
	First Int

	// (Optional) A reference to an object stream, of which the current object stream is
	// considered an extension. Both streams are considered part of a collection of object
	// streams (see below). A given collection consists of a set of streams whose Extends
	// links form a directed acyclic graph.
	Extends Reference
}

func (q ObjectStreamDictionary) ToRawBytes() []byte {
	d := q.createDict()
	d["N"] = q.N
	d["First"] = q.First
	if q.Extends.Number > 0 {
		d["Extends"] = q.Extends
	}
	return d.ToRawBytes()
}

func (q *ObjectStreamDictionary) Read(dict Dictionary, file Resolver) error {
	v, ok := dict.GetValue("Type", file)
	if !ok {
		return errors.New("ObjectStream dictionary missing type")
	}
	s, ok := v.(Name)
	if !ok {
		return errors.New("ObjectStream dictionary type is not Name")
	}
	if s != "ObjStm" {
		return errors.New("ObjectStream dictionary type is not ObjStm")
	}

	if err := q.StreamDictionary.Read(dict, file); err != nil {
		return err
	}

	v, ok = dict.GetValue("N", file)
	if !ok {
		return errors.New("ObjectStream dictionary field N missing")
	}
	q.N, ok = v.(Int)
	if !ok {
		return errors.New("ObjectStream dictionary field N invalid")
	}

	v, ok = dict.GetValue("First", file)
	if !ok {
		return errors.New("ObjectStream dictionary field First missing")
	}
	q.First, ok = v.(Int)
	if !ok {
		return errors.New("ObjectStream dictionary field First invalid")
	}

	v, ok = dict["Extends"]
	if ok {
		q.Extends, ok = v.(Reference)
		if !ok {
			return errors.New("ObjectStream dictionary field Extends invalid")
		}
	}

	// return without error
	return nil
}

func (q ObjectStreamDictionary) Copy(copyRef func(reference Reference) Reference) Object {
	return ObjectStreamDictionary{
		StreamDictionary: q.StreamDictionary.Copy(copyRef).(StreamDictionary),
		N:                q.N.Copy(copyRef).(Int),
		First:            q.First.Copy(copyRef).(Int),
		Extends:          q.Extends.Copy(copyRef).(Reference),
	}
}

func (q ObjectStreamDictionary) Equal(obj Object) bool {
	a, ok := obj.(ObjectStreamDictionary)
	if !ok {
		return false
	}

	if !Equal(q.StreamDictionary, a.StreamDictionary) {
		return false
	}
	if !Equal(q.N, a.N) {
		return false
	}
	if !Equal(q.First, a.First) {
		return false
	}
	if !Equal(q.Extends, a.Extends) {
		return false
	}
	return true
}
