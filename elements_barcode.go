package gopdf

import (
	"math"
	"strings"

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

// BarcodeElement is used to add a barcode to a page
type BarcodeElement struct {
	Left, Top, Width, Height Length
	Text                     string
	Color                    Color
	Rotate                   float64
}

// Build adds the element to the content stream
func (q *BarcodeElement) Build(page *pdf.Page) error {
	if q.Color != nil {
		q.Color.Build(page, false)
	} else {
		ColorRGBBlack.Build(page, false)
	}

	page.GraphicsState_q()
	page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), float64(page.Data.MediaBox.URY)-q.Top.Pt())
	r := q.Rotate * math.Pi / 180
	page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), 0, 0)

	code := "*" + strings.ToUpper(q.Text) + "*"
	wide := q.Width.Pt() / float64(len(code)) / 5.333
	narrow := wide / 3
	gap := narrow
	xpos := 0.0
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
				page.Path_re(xpos, -q.Height.Pt(), lineWidth, q.Height.Pt())
				page.Path_f()
			}
			xpos += lineWidth
		}
		xpos += gap
	}

	page.GraphicsState_Q()
	return nil
}
