package types

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io/ioutil"
)

// PDF Reference 1.4, 3.2.7 Stream Objects

type StreamObject struct {
	Dictionary Object
	Stream     []byte
}

func (q StreamObject) ToRawBytes() []byte {
	sb := bytes.Buffer{}
	sb.Write(q.Dictionary.ToRawBytes())

	sb.WriteString("stream\n")
	sb.Write(q.Stream)
	sb.WriteString("\n")
	sb.WriteString("endstream\n")

	return sb.Bytes()
}

func (q StreamObject) Copy(copyRef func(reference Reference) Reference) Object {
	return StreamObject{
		Dictionary: Copy(q.Dictionary, copyRef),
		Stream:     q.Stream,
	}
}

func NewStream(data []byte, filters ...Filter) (StreamObject, error) {
	for _, filter := range filters {
		switch filter {
		case Filter_ASCIIHexDecode:
			return StreamObject{}, errors.New("filter not implemented")
		case Filter_ASCII85Decode:
			return StreamObject{}, errors.New("filter not implemented")
		case Filter_LZWDecode:
			return StreamObject{}, errors.New("filter not implemented")
		case Filter_FlateDecode:
			var bts bytes.Buffer
			w := zlib.NewWriter(&bts)
			if _, err := w.Write(data); err != nil {
				return StreamObject{}, err
			}
			if err := w.Close(); err != nil {
				return StreamObject{}, err
			}
			data = bts.Bytes()

		case Filter_RunLengthDecode:
			return StreamObject{}, errors.New("filter not implemented")
		case Filter_CCITTFaxDecode:
			return StreamObject{}, errors.New("filter not implemented")
		case Filter_JBIG2Decode:
			return StreamObject{}, errors.New("filter not implemented")
		case Filter_DCTDecode:
			return StreamObject{}, errors.New("filter not implemented")
		default:
			return StreamObject{}, errors.New("unknown filter " + string(filter))
		}
	}
	return StreamObject{
		Dictionary: StreamDictionary{
			Filter: filters,
			Length: len(data),
		},
		Stream: data,
	}, nil
}

func (q *StreamObject) Decode() ([]byte, error) {
	data := q.Stream
	var filters []Filter
	if d, ok := q.Dictionary.(StreamDictionary); ok {
		filters = d.Filter
	} else if d, ok := q.Dictionary.(Dictionary); ok {
		v, ok := d["Filter"]
		if ok {
			filter, ok := v.(Name)
			if ok {
				filters = []Filter{Filter(filter)}
			} else {
				arr, ok := v.(Array)
				if !ok {
					return nil, errors.New("stream dictionary field Filter invalid")
				}
				for _, v := range arr {
					fv, ok := v.(Filter)
					if !ok {
						fvn, ok := v.(Name)
						if !ok {
							return nil, errors.New("stream dictionary field Filter invalid")
						}
						fv = Filter(fvn)
					}
					filters = append(filters, fv)
				}
			}
		}
	}

	for _, filter := range filters {
		switch filter {
		case Filter_ASCIIHexDecode:
			return nil, errors.New("filter not implemented")
		case Filter_ASCII85Decode:
			return nil, errors.New("filter not implemented")
		case Filter_LZWDecode:
			return nil, errors.New("filter not implemented")
		case Filter_FlateDecode:
			r, err := zlib.NewReader(bytes.NewReader(data))
			if err != nil {
				return nil, err
			}
			data, err = ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}

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
	return data, nil
}

func (q StreamObject) Equal(obj Object) bool {
	a, ok := obj.(StreamObject)
	if !ok {
		return false
	}
	if !Equal(q.Dictionary, a.Dictionary) {
		return false
	}
	if !bytes.Equal(q.Stream, a.Stream) {
		return false
	}
	return true
}
