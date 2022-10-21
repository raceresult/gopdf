package types

import (
	"bytes"
	"embed"
	"errors"

	"github.com/raceresult/gopdf/types/standardfont/afm"
)

type StandardFont struct {
	// Type
	// Subtype   FontSubType  // Required
	BaseFont StandardFontName // Required
	Encoding Encoding         // Name or dictionary, Optional
}

func (q StandardFont) ToRawBytes() []byte {
	d := Dictionary{
		"Type":     Name("Font"),
		"Subtype":  FontSub_Type1,
		"BaseFont": q.BaseFont,
	}
	if q.Encoding != "" {
		d["Encoding"] = q.Encoding
	}
	return d.ToRawBytes()
}

func (q *StandardFont) Read(dict Dictionary) error {
	// Type
	v, ok := dict["Type"]
	if !ok {
		return errors.New("standard font missing Type")
	}
	dtype, ok := v.(Name)
	if !ok {
		return errors.New("standard font field Type invalid")
	}
	if dtype != "Font" {
		return errors.New("unexpected value in standard font field Type")
	}

	// Subtype
	v, ok = dict["Subtype"]
	if !ok {
		return errors.New("standard font field Subtype missing")
	}
	st, ok := v.(FontSubType)
	if !ok {
		n, ok := v.(Name)
		if !ok {
			return errors.New("standard font field Subtype invalid")
		}
		st = FontSubType(n)
		if st != FontSub_Type1 {
			return errors.New("standard font field Subtype invalid")
		}
	}

	// BaseFont
	v, ok = dict["BaseFont"]
	if !ok {
		return errors.New("standard font field BaseFont missing")
	}
	q.BaseFont, ok = v.(StandardFontName)
	if !ok {
		n, ok := v.(Name)
		if !ok {
			return errors.New("standard font field Pages invalid")
		}
		q.BaseFont = StandardFontName(n)
	}

	// Encoding
	q.Encoding, ok = dict["Encoding"].(Encoding)
	if ok {
		n, ok := v.(Name)
		if !ok {
			return errors.New("standard font field Encoding invalid")
		}
		q.Encoding = Encoding(n)
	}

	// return without error
	return nil
}

var (
	//go:embed standardfont/core14
	afmFiles embed.FS
)

func (q StandardFont) Copy(copyRef func(reference Reference) Reference) Object {
	return StandardFont{
		BaseFont: q.BaseFont.Copy(copyRef).(StandardFontName),
		Encoding: q.Encoding.Copy(copyRef).(Encoding),
	}
}

// Metrics returns font metrics from the embedded afm files
func (q StandardFont) Metrics() (*afm.Font, error) {
	bts, err := afmFiles.ReadFile("standardfont/core14/" + string(q.BaseFont) + ".afm")
	if err != nil {
		return nil, err
	}
	m, err := afm.Parse(bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}
	return &m, nil
}
