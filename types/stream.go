package types

import (
	"bytes"
	"compress/zlib"
	"errors"
)

// PDF Reference 1.4, 3.4 Entries common to all stream dictionaries

type Stream struct {
	// (Required) The number of bytes from the beginning of the line fol-
	// lowing the keyword stream to the last byte just before the keyword
	// endstream. (There may be an additional EOL marker, preceding
	// endstream, that is not included in the count and is not logically part
	// of the stream data.) See “Stream Extent,” above, for further discus-
	// sion.
	//Length       int

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

	Data []byte // dictionary or array
}

func NewStream(data []byte, filters ...Filter) (*Stream, error) {
	for _, filter := range filters {
		switch filter {
		case Filter_ASCIIHexDecode:
			return nil, errors.New("filter not implemented")
		case Filter_ASCII85Decode:
			return nil, errors.New("filter not implemented")
		case Filter_LZWDecode:
			return nil, errors.New("filter not implemented")
		case Filter_FlateDecode:
			var bts bytes.Buffer
			w := zlib.NewWriter(&bts)
			if _, err := w.Write(data); err != nil {
				return nil, err
			}
			if err := w.Close(); err != nil {
				return nil, err
			}
			data = bts.Bytes()

		case Filter_RunLengthDecode:
			return nil, errors.New("filter not implemented")
		case Filter_CCITTFaxDecode:
			return nil, errors.New("filter not implemented")
		case Filter_JBIG2Decode:
			return nil, errors.New("filter not implemented")
		case Filter_DCTDecode:
			return nil, errors.New("filter not implemented")
		}
	}
	return &Stream{
		Data:   data,
		Filter: filters,
	}, nil
}

func (q Stream) createDict() Dictionary {
	d := Dictionary{
		"Length": Int(len(q.Data)),
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

func (q Stream) ToRawBytes() []byte {
	q.Data = bytes.TrimSpace(q.Data)

	sb := bytes.Buffer{}
	sb.Write(q.createDict().ToRawBytes())

	sb.WriteString("stream\n")
	sb.Write(q.Data)
	sb.WriteString("\n")
	sb.WriteString("endstream\n")

	return sb.Bytes()
}

type StreamFont struct {
	Stream

	// additional entries according to table 5.23

	// (Required for Type 1 and TrueType fonts) The length in bytes of the clear-text portion
	// of the Type 1 font program (see below), or the entire TrueType font program, after it
	// has been decoded using the filters specified by the stream’s Filter entry, if any.
	Length1 Int

	// (Required for Type 1 fonts) The length in bytes of the encrypted portion of the Type 1
	// font program (see below) after it has been decoded using the filters specified by the
	// stream’s Filter entry.
	Length2 Int

	// (Required for Type 1 fonts) The length in bytes of the fixed-content portion of the
	// Type 1 font program (see below), after it has been decoded using the filters specified
	// by the stream’s Filter entry. If Length3 is 0, it indicates that the 512 zeros and clearto-
	// mark have not been included in the FontFile font program and must be added.
	Length3 Int

	// (Required if referenced from FontFile3; PDF 1.2) A name specifying the format of the
	// embedded font program. The name must be Type1C for Type 1 compact fonts or CID-
	// FontType0C for Type 0 compact CIDFonts. When additional font formats are added
	// to PDF, more values will be defined for Subtype.
	Subtype Name

	// (Optional; PDF 1.4) A metadata stream containing metadata for the embedded font
	// program (see Section 9.2.2, “Metadata Streams”).
	Metadata Object
}

func (q StreamFont) ToRawBytes() []byte {
	q.Data = bytes.TrimSpace(q.Data)

	sb := bytes.Buffer{}
	d := q.createDict()
	if q.Length1 != 0 {
		d["Length1"] = q.Length1
	}
	if q.Length2 != 0 {
		d["Length2"] = q.Length2
	}
	if q.Length3 != 0 {
		d["Length3"] = q.Length3
	}
	if q.Subtype != "" {
		d["Subtype"] = q.Subtype
	}
	if q.Metadata != nil {
		d["Metadata"] = q.Metadata
	}
	sb.Write(d.ToRawBytes())

	sb.WriteString("stream\n")
	sb.Write(q.Data)
	sb.WriteString("\n")
	sb.WriteString("endstream\n")

	return sb.Bytes()
}
