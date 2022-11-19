package types

import "errors"

// PDF Reference 1.4, 3.4 Entries common to all stream dictionaries

type StreamDictionary struct {
	// (Required) The number of bytes from the beginning of the line fol-
	// lowing the keyword stream to the last byte just before the keyword
	// endstream. (There may be an additional EOL marker, preceding
	// endstream, that is not included in the count and is not logically part
	// of the stream data.) See “Stream Extent,” above, for further discus-
	// sion.
	Length int

	// (Optional) The name of a filter to be applied in processing the stream
	// data found between the keywords stream and endstream, or an array
	// of such names. Multiple filters should be specified in the order in
	// which they are to be applied.
	Filter []Filter // name or array

	// (Optional) A parameter dictionary, or an array of such dictionaries,
	// used by the filters specified by Filter. If there is only one filter and that
	// filter has parameters, DecodeParms must be set to the filter’s parame-
	// ter dictionary unless all the filter’s parameters have their default
	// values, in which case the DecodeParms entry may be omitted. If there
	// are multiple filters and any of the filters has parameters set to non-
	// default values, DecodeParms must be an array with one entry for
	// each filter: either the parameter dictionary for that filter, or the null
	// object if that filter has no parameters (or if all of its parameters have
	// their default values). If none of the filters have parameters, or if all
	// their parameters have default values, the DecodeParms entry may be
	// omitted. (See implementation note 7 in Appendix H.)
	DecodeParms Object // dictionary or array

	// (Optional; PDF 1.2) The file containing the stream data. If this entry
	// is present, the bytes between stream and endstream are ignored, the
	// filters are specified by FFilter rather than Filter, and the filter parame-
	// ters are specified by FDecodeParms rather than DecodeParms. How-
	// ever, the Length entry should still specify the number of those bytes.
	// (Usually there are no bytes and Length is 0.)
	F Object // file specification

	// (Optional; PDF 1.2) The name of a filter to be applied in processing
	// the data found in the stream’s external file, or an array of such names.
	// The same rules apply as for Filter.
	FFilter []Filter // name or array

	// (Optional; PDF 1.2) A parameter dictionary, or an array of such dic-
	// tionaries, used by the filters specified by FFilter. The same rules apply
	// as for DecodeParms.
	FDecodeParms Object // dictionary or array
}

func (q StreamDictionary) createDict() Dictionary {
	d := Dictionary{
		"Length": Int(q.Length),
	}
	switch len(q.Filter) {
	case 0:
	case 1:
		d["Filter"] = q.Filter[0]
	default:
		var arr Array
		for _, filter := range q.Filter {
			arr = append(arr, filter)
		}
		d["Filter"] = arr
	}

	if q.F != nil {
		d["F"] = q.F
	}
	if q.DecodeParms != nil {
		d["DecodeParms"] = q.DecodeParms
	}
	if q.FDecodeParms != nil {
		d["FDecodeParms"] = q.FDecodeParms
	}
	switch len(q.FFilter) {
	case 0:
	case 1:
		d["FFilter"] = q.Filter[0]
	default:
		var arr Array
		for _, filter := range q.Filter {
			arr = append(arr, filter)
		}
		d["FFilter"] = arr
	}

	return d
}

func (q StreamDictionary) ToRawBytes() []byte {
	return q.createDict().ToRawBytes()
}

func (q StreamDictionary) Copy(copyRef func(reference Reference) Reference) Object {
	return StreamDictionary{
		Length:       q.Length,
		Filter:       q.Filter,
		DecodeParms:  Copy(q.DecodeParms, copyRef),
		F:            Copy(q.F, copyRef),
		FFilter:      q.FFilter,
		FDecodeParms: Copy(q.FDecodeParms, copyRef),
	}
}

func (q *StreamDictionary) Read(dict Dictionary) error {
	// Size
	v, ok := dict["Length"]
	if !ok {
		return errors.New("stream dictionary missing Length")
	}
	length, ok := v.(Int)
	if !ok {
		return errors.New("stream dictionary field Length invalid")
	}
	q.Length = int(length)

	// Filter
	v, ok = dict["Filter"]
	if ok {
		filter, ok := v.(Name)
		if ok {
			q.Filter = []Filter{Filter(filter)}
		} else {
			arr, ok := v.(Array)
			if !ok {
				return errors.New("stream dictionary field Filter invalid")
			}
			for _, v := range arr {
				fv, ok := v.(Filter)
				if !ok {
					fvn, ok := v.(Name)
					if !ok {
						return errors.New("stream dictionary field Filter invalid")
					}
					fv = Filter(fvn)
				}
				q.Filter = append(q.Filter, fv)
			}
		}
	}

	// Root
	v, ok = dict["DecodeParms"]
	if ok {
		q.DecodeParms = v
	}

	// Encrypt
	v, ok = dict["F"]
	if ok {
		q.F = v
	}

	// FFilter
	v, ok = dict["FFilter"]
	if ok {
		filter, ok := v.(Name)
		if ok {
			q.FFilter = []Filter{Filter(filter)}
		} else {
			filterArr, ok := v.(Array)
			if !ok {
				return errors.New("stream dictionary field FFilter invalid")
			}
			for _, v := range filterArr {
				fv, ok := v.(Filter)
				if !ok {
					fvn, ok := v.(Name)
					if !ok {
						return errors.New("stream dictionary field FFilter invalid")
					}
					fv = Filter(fvn)
				}
				q.FFilter = append(q.FFilter, Filter(fv))
			}
		}
	}

	// Root
	v, ok = dict["FDecodeParms"]
	if ok {
		q.FDecodeParms = v
	}

	// return without error
	return nil
}

func (q StreamDictionary) Equal(obj Object) bool {
	a, ok := obj.(StreamDictionary)
	if !ok {
		return false
	}
	if q.Length != a.Length {
		return false
	}

	if len(q.Filter) != len(a.Filter) {
		return false
	}
	for i := range q.FFilter {
		if !Equal(q.Filter[i], a.Filter[i]) {
			return false
		}
	}

	if !Equal(q.DecodeParms, a.DecodeParms) {
		return false
	}
	if !Equal(q.F, a.F) {
		return false
	}

	if len(q.FFilter) != len(a.FFilter) {
		return false
	}
	for i := range q.FFilter {
		if !Equal(q.FFilter[i], a.FFilter[i]) {
			return false
		}
	}

	if !Equal(q.FDecodeParms, a.FDecodeParms) {
		return false
	}
	return true
}
