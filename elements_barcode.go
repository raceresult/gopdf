package gopdf

import (
	"errors"
	"math"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/raceresult/gopdf/pdf"
)

var code39Chars = map[rune]string{
	'0': "nnnwwnwnn",
	'1': "wnnwnnnnw",
	'2': "nnwwnnnnw",
	'3': "wnwwnnnnn",
	'4': "nnnwwnnnw",
	'5': "wnnwwnnnn",
	'6': "nnwwwnnnn",
	'7': "nnnwnnwnw",
	'8': "wnnwnnwnn",
	'9': "nnwwnnwnn",
	'A': "wnnnnwnnw",
	'B': "nnwnnwnnw",
	'C': "wnwnnwnnn",
	'D': "nnnnwwnnw",
	'E': "wnnnwwnnn",
	'F': "nnwnwwnnn",
	'G': "nnnnnwwnw",
	'H': "wnnnnwwnn",
	'I': "nnwnnwwnn",
	'J': "nnnnwwwnn",
	'K': "wnnnnnnww",
	'L': "nnwnnnnww",
	'M': "wnwnnnnwn",
	'N': "nnnnwnnww",
	'O': "wnnnwnnwn",
	'P': "nnwnwnnwn",
	'Q': "nnnnnnwww",
	'R': "wnnnnnwwn",
	'S': "nnwnnnwwn",
	'T': "nnnnwnwwn",
	'U': "wwnnnnnnw",
	'V': "nwwnnnnnw",
	'W': "wwwnnnnnn",
	'X': "nwnnwnnnw",
	'Y': "wwnnwnnnn",
	'Z': "nwwnwnnnn",
	'-': "nwnnnnwnw",
	'.': "wwnnnnwnn",
	' ': "nwwnnnwnn",
	'*': "nwnnwnwnn",
	'$': "nwnwnwnnn",
	'/': "nwnwnnnwn",
	'+': "nwnnnwnwn",
	'%': "nnnwnwnwn",
}

type BarcodeType int

const (
	BarcodeType39  BarcodeType = 1
	BarcodeType128 BarcodeType = 3
)

// BarcodeElement is used to add a barcode to a page
type BarcodeElement struct {
	Left, Top, Width, Height Length
	Text                     string
	Color                    Color
	Rotate                   float64
	Flipped                  bool
	Type                     BarcodeType
}

// Build adds the element to the content stream
func (q *BarcodeElement) Build(page *pdf.Page) error {
	// set color
	if q.Color != nil {
		q.Color.Build(page, false)
	} else {
		ColorRGBBlack.Build(page, false)
	}

	// set position
	page.GraphicsState_q()
	page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), float64(page.Data.MediaBox.URY)-q.Top.Pt())
	r := q.Rotate * math.Pi / 180
	page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), 0, 0)

	// encode
	var width float64
	if q.Flipped {
		width = q.Height.Pt()
	} else {
		width = q.Width.Pt()
	}
	bars, err := q.encode(width)
	if err != nil {
		return err
	}

	// build rectangles
	if q.Flipped {
		for _, bar := range bars {
			page.Path_re(0, bar[0]-q.Height.Pt(), q.Width.Pt(), bar[1])
		}
	} else {
		for _, bar := range bars {
			page.Path_re(bar[0], -q.Height.Pt(), bar[1], q.Height.Pt())
		}
	}
	page.Path_f()

	// reset graphics state
	page.GraphicsState_Q()
	return nil
}

// encode returns a list of x-pos + width representing the bars of the barcode
func (q *BarcodeElement) encode(width float64) ([][2]float64, error) {
	switch q.Type {
	case BarcodeType39:
		return q.encodeCode39(width)
	case BarcodeType128:
		return q.encodeCode128(width)
	default:
		return nil, errors.New("unknown barcode type")
	}
}

// encodeCode39 returns a list of x-pos + width representing the bars of a code39 barcode
func (q *BarcodeElement) encodeCode39(width float64) ([][2]float64, error) {
	code := "*" + strings.ToUpper(q.Text) + "*"
	wide := width / float64(len(code)) / 5.333

	var res [][2]float64
	narrow := wide / 3
	gap := narrow
	pos := 0.0
	for _, char := range code {
		seq, ok := code39Chars[char]
		if !ok {
			continue
		}
		for bar := 0; bar < 9; bar++ {
			lineWidth := wide
			if seq[bar] == 'n' {
				lineWidth = narrow
			}
			if bar%2 == 0 {
				res = append(res, [2]float64{pos, lineWidth})
			}
			pos += lineWidth
		}
		pos += gap
	}

	return res, nil
}

// encodeCode128 returns a list of x-pos + width representing the bars of a code128 barcode
func (q *BarcodeElement) encodeCode128(width float64) ([][2]float64, error) {
	code, err := code128.Encode(q.Text)
	if err != nil {
		return nil, err
	}
	w := int(width) * 100
	scaled, err := barcode.Scale(code, w, 1)
	if err != nil {
		return nil, err
	}

	var res [][2]float64
	start := -1
	for x := 0; x < w; x++ {
		v, _, _, _ := scaled.At(x, 0).RGBA()
		if start >= 0 && (x == w-1 || v == 0) {
			res = append(res, [2]float64{float64(start) / 100, float64(x-start) / 100})
			start = -1
		}
		if start < 0 && v != 0 {
			start = x
		}
	}

	return res, nil
}
