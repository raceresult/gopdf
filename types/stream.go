package types

import (
	"bytes"
	"compress/zlib"
	"encoding/ascii85"
	"encoding/hex"
	"errors"
	"io/ioutil"

	"github.com/raceresult/gopdf/types/runlength"
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
			data = []byte(hex.EncodeToString(data))

		case Filter_ASCII85Decode:
			var bts bytes.Buffer
			w := ascii85.NewEncoder(&bts)
			if _, err := w.Write(data); err != nil {
				return StreamObject{}, err
			}
			if err := w.Close(); err != nil {
				return StreamObject{}, err
			}
			data = bts.Bytes()

		case Filter_LZWDecode:
			return StreamObject{}, errors.New("filter " + string(filter) + " not implemented")

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
			var err error
			data, err = runlength.Encode(bytes.NewReader(data))
			if err != nil {
				return StreamObject{}, err
			}
		case Filter_CCITTFaxDecode:
			return StreamObject{}, errors.New("filter " + string(filter) + " not implemented")
		case Filter_JBIG2Decode:
			return StreamObject{}, errors.New("filter " + string(filter) + " not implemented")
		case Filter_DCTDecode:
			return StreamObject{}, errors.New("filter " + string(filter) + " not implemented")
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

func (q *StreamObject) getFilters() ([]Filter, error) {
	switch d := q.Dictionary.(type) {
	case StreamDictionary:
		return d.Filter, nil

	case Dictionary:
		v, ok := d["Filter"]
		if !ok {
			return nil, nil
		}
		filter, ok := v.(Name)
		if ok {
			return []Filter{Filter(filter)}, nil
		}

		arr, ok := v.(Array)
		if !ok {
			return nil, errors.New("stream dictionary field Filter invalid")
		}
		var filters []Filter
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
		return filters, nil

	default:
		return nil, nil
	}
}

// getDecodeParams returns the FilterParams in DecodeParms
func (q *StreamObject) getDecodeParams(file Resolver) ([]FilterParameters, error) {
	// get stream dictionary
	var sd StreamDictionary
	switch d := q.Dictionary.(type) {
	case StreamDictionary:
		sd = d

	case Dictionary:
		if err := sd.Read(d, file); err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("could not get DecodeParms of stream")
	}

	// read
	switch dp := sd.DecodeParms.(type) {
	case Array:
		var fps []FilterParameters
		for _, item := range dp {
			d, ok := item.(Dictionary)
			if !ok {
				return nil, errors.New("unexpected value in DecodeParms array")
			}
			var fp FilterParameters
			if err := fp.Read(d, file); err != nil {
				return nil, err
			}
			fps = append(fps, fp)
		}
		return fps, nil

	case Dictionary:
		var fp FilterParameters
		if err := fp.Read(dp, file); err != nil {
			return nil, err
		}
		return []FilterParameters{fp}, nil

	case nil:
		return nil, nil

	default:
		return nil, errors.New("unexpected value in DecodeParms of stream")
	}
}

func (q *StreamObject) Decode(file Resolver) ([]byte, error) {
	// get filter list and decodeParams list
	filters, err := q.getFilters()
	if err != nil {
		return nil, err
	}
	decodeParms, err := q.getDecodeParams(file)
	if err != nil {
		return nil, err
	}

	// decode
	data := q.Stream
	for i, filter := range filters {
		switch filter {
		case Filter_ASCIIHexDecode:
			data, err = hex.DecodeString(string(data))
			if err != nil {
				return nil, err
			}

		case Filter_ASCII85Decode:
			data = bytes.ReplaceAll(data, []byte{10}, nil)
			data = bytes.ReplaceAll(data, []byte{13}, nil)
			data = bytes.TrimSuffix(data, []byte{'~', '>'})
			r := ascii85.NewDecoder(bytes.NewReader(data))
			data, err = ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}

		case Filter_LZWDecode:
			return nil, errors.New("filter " + string(filter) + " not implemented")

		case Filter_FlateDecode:
			// for some reason, sometimes the data is missing at least one byte, although the encoded stream has the
			// exact length as defined in the stream dictionary. Decoding the flat encoded data then fails, but
			// other PDF readers accept the file. As a workaround, we add 4 extra zero bytes.
			c := make([]byte, len(data)+4)
			copy(c, data)

			r, err := zlib.NewReader(bytes.NewReader(c))
			if err != nil {
				return nil, err
			}
			data, err = ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}

		case Filter_RunLengthDecode:
			return runlength.Decode(bytes.NewReader(data))
		case Filter_CCITTFaxDecode:
			return nil, errors.New("filter " + string(filter) + " not implemented")
		case Filter_JBIG2Decode:
			return nil, errors.New("filter " + string(filter) + " not implemented")
		case Filter_DCTDecode:
			return nil, errors.New("filter " + string(filter) + " not implemented")
		}

		if i < len(decodeParms) {
			data, err = decodePredictor(data, decodeParms[i])
			if err != nil {
				return nil, err
			}
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
