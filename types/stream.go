package types

import (
	"bytes"
	"compress/zlib"
	"errors"
)

// PDF Reference 1.4, 3.2.7 Stream Objects

type StreamObject struct {
	Dictionary Object
	Stream     []byte
}

func (q StreamObject) ToRawBytes() []byte {
	q.Stream = bytes.TrimSpace(q.Stream)

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
