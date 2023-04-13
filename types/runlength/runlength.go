package runlength

import (
	"bytes"
	"io"
)

func Decode(bufReader *bytes.Reader) ([]byte, error) {
	inb := []byte{}
	for {
		b, err := bufReader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b > 128 {
			v, err := bufReader.ReadByte()
			if err != nil {
				return nil, err
			}
			for i := 0; i < 257-int(b); i++ {
				inb = append(inb, v)
			}
		} else if b < 128 {
			for i := 0; i < int(b)+1; i++ {
				v, err := bufReader.ReadByte()
				if err != nil {
					return nil, err
				}
				inb = append(inb, v)
			}
		} else {
			break
		}
	}

	return inb, nil
}

func Encode(bufReader *bytes.Reader) ([]byte, error) {
	inb := []byte{}
	literal := []byte{}

	b0, err := bufReader.ReadByte()
	if err == io.EOF {
		return []byte{}, nil
	} else if err != nil {
		return nil, err
	}
	runLen := 1

	for {
		b, err := bufReader.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if b == b0 {
			if len(literal) > 0 {
				literal = literal[:len(literal)-1]
				if len(literal) > 0 {
					inb = append(inb, byte(len(literal)-1))
					inb = append(inb, literal...)
				}
				runLen = 1
				literal = []byte{}
			}
			runLen++
			if runLen >= 127 {
				inb = append(inb, byte(257-runLen), b0)
				runLen = 0
			}

		} else {
			if runLen > 0 {
				if runLen == 1 {
					literal = []byte{b0}
				} else {
					inb = append(inb, byte(257-runLen), b0)
				}

				runLen = 0
			}
			literal = append(literal, b)
			if len(literal) >= 127 {
				inb = append(inb, byte(len(literal)-1))
				inb = append(inb, literal...)
				literal = []byte{}
			}
		}
		b0 = b
	}

	if len(literal) > 0 {
		inb = append(inb, byte(len(literal)-1))
		inb = append(inb, literal...)
	} else if runLen > 0 {
		inb = append(inb, byte(257-runLen), b0)
	}
	inb = append(inb, 128)
	return inb, nil
}
