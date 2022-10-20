package types

import (
	"bytes"
	"embed"

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
