package parser

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
)

// readFile reads a byte stream into a File object
func readFile(bts []byte) (*pdffile.File, error) {
	length := len(bts)

	// try to read xref table
	var xref pdffile.XRefTable
	startxref := bytes.LastIndex(bts, []byte("startxref"))
	if startxref >= 0 {
		startXRefObj, _, _ := readValue(bts[startxref+9:])
		startXRefVal, ok := startXRefObj.(types.Int)
		if ok && int(startXRefVal) < len(bts) {
			xref, _, _ = readXRef(bts[startXRefVal:])
		}
	}

	// parse PDF version
	var dest pdffile.File
	var firstLine []byte
	firstLine, bts = readLine(bts)
	if !bytes.HasPrefix(firstLine, []byte("%PDF-")) {
		return nil, errors.New("file does not have %PDF- prefix")
	}
	dest.Version, _ = strconv.ParseFloat(string(firstLine[5:]), 64)
	if dest.Version <= 0 {
		return nil, errors.New("file does not have valid PDF version number")
	}

	// parse objects
	objectOffsets := make(map[int]types.IndirectObject)
	var trailer types.Trailer
	for len(bts) != 0 {
		var err error
		switch {
		case bytes.HasPrefix(bts, []byte("%%EOF")):
			bts = bts[5:]

		case bytes.HasPrefix(bts, []byte("%")):
			_, bts, err = readComment(bts)
			if err != nil {
				return nil, err
			}

		case bytes.HasPrefix(bts, []byte("xref")):
			_, bts, err = readXRef(bts)
			if err != nil {
				return nil, err
			}

		case bytes.HasPrefix(bts, []byte("trailer")):
			trailer, bts, err = readTrailer(bts)
			if err != nil {
				// if we already had a trailer, it is probably a linearized pdf
				if dest.Root.Number == 0 {
					return nil, err
				}
			} else {
				dest.ID = trailer.ID
				dest.Root = trailer.Root
				dest.Info = trailer.Info
			}

		case bytes.HasPrefix(bts, []byte("startxref")):
			var startXRef types.Object
			startXRef, bts, err = readValue(bts[9:])
			if err != nil {
				return nil, err
			}
			if dest.Root.Number == 0 {
				if offset, ok := startXRef.(types.Int); ok && offset > 0 {
					if obj, ok := objectOffsets[int(offset)]; ok {
						if so, ok := obj.Data.(types.StreamObject); ok {
							if dict, ok := so.Dictionary.(types.Dictionary); ok {
								var t types.Trailer
								if err := t.Read(dict); err == nil {
									dest.ID = t.ID
									dest.Root = t.Root
									dest.Info = t.Info
								}
							}
						}
					}
				}
			}

		case isWhiteChar(bts[0]):
			bts = bts[1:]

		default:
			offset := length - len(bts)
			var obj types.IndirectObject
			obj, bts, err = readObject(bts, xref, length)
			if err != nil {
				return nil, err
			}
			dest.AddIndirectObject(obj)
			objectOffsets[offset] = obj
		}
	}

	// return without error
	return &dest, nil
}

func readTrailer(bts []byte) (types.Trailer, []byte, error) {
	bts = trimLeftWhiteChars(bts)
	if !bytes.HasPrefix(bts, []byte("trailer")) {
		return types.Trailer{}, bts, errors.New("value is not a trailer")
	}
	bts = trimLeftWhiteChars(bts[7:])

	var err error
	var trailerDict types.Dictionary
	trailerDict, bts, err = readDictionary(bts)
	if err != nil {
		return types.Trailer{}, bts, err
	}
	var trailer types.Trailer
	if err := trailer.Read(trailerDict); err != nil {
		return types.Trailer{}, bts, err
	}
	return trailer, bts, nil
}

func findObject(ref types.Reference, bts []byte, xref pdffile.XRefTable, length int) (types.Object, error) {
	for _, x := range xref {
		if x.Start > ref.Number || x.Start+x.Count-1 <= ref.Number {
			continue
		}
		entry := x.Entries[x.Start+ref.Number]
		if entry.Generation != ref.Generation {
			continue
		}
		start := int(entry.Start) - (length - len(bts))
		if start < len(bts) {
			res, _, err := readObject(bts[start:], xref, length)
			return res.Data, err
		}
	}
	return nil, errors.New("object " + strconv.Itoa(ref.Number) + "/" + strconv.Itoa(ref.Generation) + " not found")
}

func readObject(bts []byte, xref pdffile.XRefTable, length int) (types.IndirectObject, []byte, error) {
	var header [3][]byte
	for i := 0; i < 3; i++ {
		header[i], bts = readWord(bts)
	}

	if string(header[2]) != "obj" {
		return types.IndirectObject{}, bts, errors.New("error parsing object header")
	}
	id, err := strconv.Atoi(string(header[0]))
	if err != nil || id <= 0 {
		return types.IndirectObject{}, bts, errors.New("object header invalid")
	}
	gen, err := strconv.Atoi(string(header[1]))
	if err != nil || gen < 0 {
		return types.IndirectObject{}, bts, errors.New("object header invalid")
	}

	var obj types.Object
	obj, bts, err = readAny(bts)
	if err != nil {
		return types.IndirectObject{}, bts, err
	}

	bts = trimLeftWhiteChars(bts)
	if bytes.HasPrefix(bts, []byte("stream")) {
		_, bts = readLine(bts)

		switch v := obj.(type) {
		case types.Dictionary:
			streamLength, ok := v["Length"]
			if !ok {
				return types.IndirectObject{}, bts, errors.New("stream dictionary does not have length")
			}
			lengthVal, ok := streamLength.(types.Int)
			if !ok {
				lengthRef, ok := streamLength.(types.Reference)
				if !ok {
					return types.IndirectObject{}, bts, errors.New("stream dictionary Length invalid")
				}
				lengthObj, err := findObject(lengthRef, bts, xref, length)
				if err != nil {
					return types.IndirectObject{}, bts, err
				}
				lengthVal, ok = lengthObj.(types.Int)
				if !ok {
					return types.IndirectObject{}, bts, errors.New("stream dictionary Length invalid")
				}
			}

			stream := types.StreamObject{
				Dictionary: obj,
				Stream:     bts[:lengthVal],
			}
			bts = bts[lengthVal:]
			obj = stream

		default:
			return types.IndirectObject{}, bts, errors.New("stream does not have dictionary")
		}

		bts = trimLeftWhiteChars(bts)
		if !bytes.HasPrefix(bts, []byte("endstream")) {
			return types.IndirectObject{}, bts, errors.New("unterminated stream")
		}
		bts = bts[9:]
		bts = trimLeftWhiteChars(bts)
	}
	if !bytes.HasPrefix(bts, []byte("endobj")) {
		return types.IndirectObject{}, bts, errors.New("unterminated object")
	}
	bts = bts[6:]

	return types.IndirectObject{
		Number:     id,
		Generation: gen,
		Data:       obj,
	}, bts, nil
}

func readLine(bts []byte) ([]byte, []byte) {
	for i := 0; i < len(bts); i++ {
		switch bts[i] {
		case '\n':
			return bts[:i], bts[i+1:]
		case '\r':
			if i < len(bts)-1 && bts[i+1] == '\n' {
				return bts[:i], bts[i+2:]
			}
			return bts[:i], bts[i+1:]
		default:
		}
	}
	return bts, nil
}

func readWord(bts []byte) ([]byte, []byte) {
	bts = trimLeftWhiteChars(bts)
	for i := 0; i < len(bts); i++ {
		c := bts[i]
		if isWhiteChar(c) || isDelimiterChar(c) {
			return bts[:i], bts[i:]
		}
	}
	return bts, nil
}

func nextTwoWords(bts []byte) ([]byte, []byte) {
	bts = trimLeftWhiteChars(bts)

	var w1 []byte
	var start int
	for i := 0; i < len(bts); i++ {
		c := bts[i]
		if isWhiteChar(c) || isDelimiterChar(c) {
			if len(w1) == 0 {
				w1 = bts[:i]
				start = i + 1
			} else if i == start {
				start = i + 1
			} else {
				return w1, bts[start:i]
			}
		}
	}
	return w1, bts[start:]
}

func readXRef(bts []byte) (pdffile.XRefTable, []byte, error) {
	bts = trimLeftWhiteChars(bts)
	if !bytes.HasPrefix(bts, []byte("xref")) {
		return nil, bts, errors.New("value is not a n xref table")
	}
	bts = trimLeftWhiteChars(bts[4:])

	var xrefTable pdffile.XRefTable
	for {
		if len(bts) == 0 {
			return nil, nil, errors.New("unexpected end of xref table")
		}
		if bts[0] != '0' && bts[0] != '1' && bts[0] != '2' && bts[0] != '3' && bts[0] != '4' &&
			bts[0] != '5' && bts[0] != '6' && bts[0] != '7' && bts[0] != '8' && bts[0] != '9' {
			break
		}

		var line []byte
		line, bts = readLine(bts)

		ww := splitWords(string(line))
		if len(ww) != 2 {
			return nil, bts, errors.New("invalid xref table section start")
		}

		secStart, _ := strconv.Atoi(ww[0])
		secCount, _ := strconv.Atoi(ww[1])
		if secStart < 0 || secCount < 0 || secCount*20 > len(bts) {
			return nil, nil, errors.New("invalid xref table section start")
		}

		lines := bts[:20*secCount]
		bts = bts[20*secCount:]
		xrefTable = append(xrefTable, pdffile.XRefTableSection{
			Start: secStart,
			Count: secCount,
		})

		for len(lines) >= 20 {
			item := lines[:20]
			lines = lines[20:]

			start, _ := strconv.Atoi(string(item[:10]))
			gen, _ := strconv.Atoi(string(item[12:17]))
			xrefTable[len(xrefTable)-1].Entries = append(xrefTable[len(xrefTable)-1].Entries, pdffile.XRefTableEntry{
				Start:      int64(start),
				Generation: gen,
				Free:       item[18] == 'f',
			})
		}
	}
	return xrefTable, bts, nil
}

func readAny(bts []byte) (types.Object, []byte, error) {
	for {
		switch {
		case len(bts) == 0:
			return nil, bts, errors.New("eof")

		case bytes.HasPrefix(bts, []byte("<<")):
			return readDictionary(bts)

		case bytes.HasPrefix(bts, []byte("/")):
			n, rest, err := readName(bts)
			return n, rest, err

		case bytes.HasPrefix(bts, []byte("(")):
			s, rest, err := readString(bts)
			if err != nil {
				return nil, rest, err
			}
			if strings.HasPrefix(string(s), "D:") {
				t, err := time.Parse("20060102150405-07'00'", string(s)[2:])
				if err == nil {
					return types.Date(t), rest, err
				}
			}
			return s, rest, err

		case bytes.HasPrefix(bts, []byte("<")):
			s, rest, err := readHexString(bts)
			return s, rest, err

		case bytes.HasPrefix(bts, []byte("%")):
			_, bts, err := readComment(bts)
			if err != nil {
				return nil, bts, err
			}

		case bytes.HasPrefix(bts, []byte("[")):
			a, bts, err := readArray(bts)
			return a, bts, err

		case isWhiteChar(bts[0]):
			bts = bts[1:]

		default:
			return readValue(bts)
		}
	}
}

func readArray(bts []byte) (types.Array, []byte, error) {
	bts = trimLeftWhiteChars(bts)
	if !bytes.HasPrefix(bts, []byte("[")) {
		return nil, bts, errors.New("value is not an Array")
	}
	bts = bts[1:]

	var dest types.Array
	for {
		bts = trimLeftWhiteChars(bts)
		if len(bts) == 0 {
			return nil, bts, errors.New("unterminated Array")
		}
		if bts[0] == ']' {
			return dest, bts[1:], nil
		}

		var err error
		var v types.Object
		v, bts, err = readAny(bts)
		if err != nil {
			return nil, bts, err
		}

		dest = append(dest, v)
	}
}

func readComment(bts []byte) (string, []byte, error) {
	bts = trimLeftWhiteChars(bts)
	if !bytes.HasPrefix(bts, []byte("%")) {
		return "", bts, errors.New("value is not a comment")
	}
	bts = bts[1:]

	var s []byte
	for i := 0; i < len(bts); i++ {
		c := bts[i]
		switch c {
		case '\n':
			return string(s), bts[i+1:], nil
		case '\r':
			if i < len(bts)-1 && bts[i+1] == '\n' {
				return string(s), bts[i+2:], nil
			}
			return string(s), bts[i+1:], nil
		default:
			s = append(s, c)
		}
	}
	return string(s), nil, nil
}

func readString(bts []byte) (types.String, []byte, error) {
	bts = trimLeftWhiteChars(bts)
	if !bytes.HasPrefix(bts, []byte("(")) {
		return "", bts, errors.New("value is not a String")
	}
	bts = bts[1:]

	var s []byte
	var openParenthesis int
	for i := 0; i < len(bts); i++ {
		c := bts[i]
		switch c {
		case '(':
			openParenthesis++
			s = append(s, c)
		case ')':
			if openParenthesis > 0 {
				openParenthesis--
				s = append(s, c)
				continue
			}
			return types.String(s), bts[i+1:], nil
		case '\\':
			if i >= len(bts)-1 {
				return "", bts, errors.New("invalid String escaping")
			}
			i++
			switch bts[i] {
			case 'n':
				s = append(s, '\n')
			case 'r':
				s = append(s, '\r')
			case 't':
				s = append(s, '\t')
			case 'b':
				s = append(s, '\b')
			case 'f':
				s = append(s, '\f')
			case '(':
				s = append(s, '(')
			case ')':
				s = append(s, ')')
			case '\\':
				s = append(s, '\\')
			case '\n':
				s = append(s, '\n')
			case '\r':
				s = append(s, '\n')
				if i < len(bts)-1 && bts[i+1] == '\r' {
					i++
				}
			default:
				if i >= len(bts)-2 {
					return "", bts, errors.New("invalid String escaping")
				}
				charCode, err := strconv.ParseInt(string(bts[i:i+3]), 8, 64)
				if err != nil || charCode < 0 {
					return "", bts, errors.New("invalid String escaping")
				}
				s = append(s, byte(charCode))
			}

		default:
			s = append(s, c)
		}
	}
	return "", bts, errors.New("unterminated String")
}

func readHexString(bts []byte) (types.String, []byte, error) {
	bts = trimLeftWhiteChars(bts)
	if !bytes.HasPrefix(bts, []byte("<")) {
		return "", bts, errors.New("value is not a Hex String")
	}
	bts = bts[1:]

	var s []byte
	for i := 0; i < len(bts); i += 2 {
		if bts[i] == '>' {
			return types.String(s), bts[i+1:], nil
		}

		charCode, err := strconv.ParseInt(string(bts[i:i+2]), 16, 64)
		if err != nil || charCode < 0 {
			return "", bts, errors.New("invalid Hex String")
		}
		s = append(s, byte(charCode))
	}
	return "", bts, errors.New("unterminated Hex String")
}

func readValue(bts []byte) (types.Object, []byte, error) {
	var w []byte
	w, bts = readWord(bts)
	if string(w) == "true" {
		return types.Boolean(true), bts, nil
	}
	if string(w) == "false" {
		return types.Boolean(false), bts, nil
	}
	if string(w) == "null" {
		return types.Null{}, bts, nil
	}
	if v, err := strconv.Atoi(string(w)); err == nil {
		w1, w2 := nextTwoWords(bts)
		if len(w1) != 0 && len(w2) == 1 && w2[0] == 'R' {
			gen, err := strconv.Atoi(string(w1))
			if err == nil && v >= 0 {
				w1, bts = readWord(bts)
				w2, bts = readWord(bts)

				return types.Reference{
					Number:     v,
					Generation: gen,
				}, bts, nil
			}
		}

		return types.Int(v), bts, nil
	}
	if v, err := strconv.ParseFloat(string(w), 64); err == nil {
		return types.Number(v), bts, nil
	}

	return nil, bts, errors.New("unknown value type")
}

func readDictionary(bts []byte) (types.Dictionary, []byte, error) {
	bts = trimLeftWhiteChars(bts)
	if !bytes.HasPrefix(bts, []byte("<<")) {
		return nil, bts, errors.New("value is not a Dictionary")
	}
	bts = bts[2:]

	dest := make(types.Dictionary)
	for {
		bts = trimLeftWhiteChars(bts)
		if bytes.HasPrefix(bts, []byte(">>")) {
			return dest, bts[2:], nil
		}
		if len(bts) == 0 {
			return nil, bts, errors.New("dictionary end not found")
		}

		var err error
		var key types.Name
		key, bts, err = readName(bts)
		if err != nil {
			return nil, bts, err
		}

		var value types.Object
		value, bts, err = readAny(bts)
		if err != nil {
			return nil, bts, err
		}

		dest[key] = value
	}
}

func readName(bts []byte) (types.Name, []byte, error) {
	bts = trimLeftWhiteChars(bts)
	if !bytes.HasPrefix(bts, []byte("/")) {
		return "", bts, errors.New("value is not a Name")
	}
	bts = bts[1:]

	var name []byte
	for i := 0; i < len(bts); i++ {
		c := bts[i]
		if isWhiteChar(c) || isDelimiterChar(c) {
			bts = bts[i:]
			break
		}
		if c == '#' && i < len(bts)-2 {
			charCode, err := strconv.ParseInt(string(bts[i+1:i+3]), 16, 64)
			if err != nil {
				return "", bts, errors.New("name is invalid hex code")
			}
			name = append(name, byte(charCode))
			i += 2
			continue
		}

		name = append(name, c)
	}
	if len(name) == 0 {
		return "", bts, errors.New("invalid name")
	}
	return types.Name(name), bts, nil
}

func isWhiteChar(c byte) bool {
	switch c {
	case 0, 9, 10, 12, 13, 32:
		return true
	default:
		return false
	}
}

func isDelimiterChar(c byte) bool {
	switch c {
	case '(', ')', '<', '>', '[', ']', '{', '}', '/', '%':
		return true
	default:
		return false
	}
}

func trimLeftWhiteChars(bts []byte) []byte {
	return bytes.TrimLeftFunc(bts, func(r rune) bool {
		if r > 32 {
			return false
		}
		return isWhiteChar(byte(r))
	})
}

func splitWords(s string) []string {
	var words []string
	for s != "" {
		var found bool
		for i, c := range s {
			if isWhiteChar(byte(c)) {
				if i != 0 {
					words = append(words, s[:i])
				}
				s = s[i+1:]
				found = true
				break
			}
		}
		if !found {
			words = append(words, s)
			break
		}
	}
	return words
}
