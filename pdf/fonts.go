package pdf

import (
	"sync"

	"github.com/raceresult/gopdf/pdf/unitype"
	"github.com/raceresult/gopdf/types"
	"github.com/raceresult/gopdf/types/standardfont/afm"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/encoding/unicode"
)

type FontHandler interface {
	finish() error
	Reference() types.Reference
	Encode(s string) string
	GetWidth(text string, fontSize float64) float64
	GetAscent(fontSize float64) float64
	GetUnderlineThickness(size float64) float64
	GetUnderlinePosition(size float64) float64
	GetTop(fontSize float64) float64
	GetBottom(fontSize float64) float64
	GetHeight(fontSize float64) float64
	HasGylph(runes []rune) []bool
	FallbackFont() FontHandler
}

// ---------------------------------------------------------------------------------------------------------------------

// StandardFont references a standard font and provides additional function like font metrics
type StandardFont struct {
	reference types.Reference
	encoding  types.Encoding
	metrics   *afm.Font
}

func (q *StandardFont) Encode(text string) string {
	return q.encoding.Encode(text)
}
func (q *StandardFont) Reference() types.Reference {
	return q.reference
}
func (q *StandardFont) GetWidth(text string, fontSize float64) float64 {
	var w int
	for _, c := range text {
		w += q.metrics.GetGlyphAdvance(int(c))
	}
	return float64(w) * fontSize / 1000
}
func (q *StandardFont) GetAscent(fontSize float64) float64 {
	return q.metrics.Ascender.Float64() * fontSize / 1000
}
func (q *StandardFont) GetTop(fontSize float64) float64 {
	return q.metrics.BBox.URY.Float64() * fontSize / 1000
}
func (q *StandardFont) GetBottom(fontSize float64) float64 {
	return q.metrics.BBox.LLY.Float64() * fontSize / 1000
}
func (q *StandardFont) GetHeight(fontSize float64) float64 {
	return (q.metrics.BBox.URY - q.metrics.BBox.LLY).Float64() * fontSize / 1000
}
func (q *StandardFont) GetUnderlineThickness(fontSize float64) float64 {
	return q.metrics.Direction[0].UnderlineThickness.Float64() * fontSize / 1000
}
func (q *StandardFont) GetUnderlinePosition(fontSize float64) float64 {
	return q.metrics.Direction[0].UnderlinePosition.Float64() * fontSize / 1000
}
func (q *StandardFont) HasGylph(runes []rune) []bool {
	dest := make([]bool, 0, len(runes))
	for range runes {
		dest = append(dest, true)
	}
	return dest
}
func (q *StandardFont) FallbackFont() FontHandler {
	return nil
}
func (q *StandardFont) finish() error {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// TrueTypeFont references a TrueType font and provides additional function like font metrics
type TrueTypeFont struct {
	reference types.Reference
	encoding  types.Encoding
	font      *unitype.Font
	metrics   unitype.Metrics
}

func (q *TrueTypeFont) Encode(text string) string {
	return q.encoding.Encode(text)
}
func (q *TrueTypeFont) Reference() types.Reference {
	return q.reference
}
func (q *TrueTypeFont) GetWidth(text string, fontSize float64) float64 {
	var w int
	for _, ind := range q.font.LookupRunes([]rune(text)) {
		w += q.font.GetGlyphAdvance(ind)
	}
	return float64(w) * fontSize / 1000
}
func (q *TrueTypeFont) GetAscent(fontSize float64) float64 {
	return float64(q.metrics.Ascent) * fontSize / 1000
}
func (q *TrueTypeFont) GetTop(fontSize float64) float64 {
	return float64(q.metrics.YMax) * fontSize / 1000
}
func (q *TrueTypeFont) GetBottom(fontSize float64) float64 {
	return float64(q.metrics.YMin) * fontSize / 1000
}
func (q *TrueTypeFont) GetHeight(fontSize float64) float64 {
	return float64(q.metrics.YMax-q.metrics.YMin) * fontSize / 1000
}
func (q *TrueTypeFont) GetUnderlineThickness(size float64) float64 {
	return float64(q.metrics.UnderlineThickness) * size / 1000
}
func (q *TrueTypeFont) GetUnderlinePosition(size float64) float64 {
	return float64(q.metrics.UnderlinePosition) * size / 1000
}
func (q *TrueTypeFont) HasGylph(runes []rune) []bool {
	dest := make([]bool, 0, len(runes))
	for _, ind := range q.font.LookupRunes(runes) {
		dest = append(dest, ind > 0)
	}
	return dest
}
func (q *TrueTypeFont) FallbackFont() FontHandler {
	return nil
}
func (q *TrueTypeFont) finish() error {
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// CompositeFont references a composite font and provides additional function like font metrics
type CompositeFont struct {
	reference    types.Reference
	usedRunes    map[rune]struct{}
	usedRunesMux sync.Mutex
	onFinish     func() error
	font         *unitype.Font
	metrics      unitype.Metrics
	fallbackFont FontHandler
}

func (q *CompositeFont) Reference() types.Reference {
	return q.reference
}
func (q *CompositeFont) Encode(text string) string {
	var repl bool
	q.usedRunesMux.Lock()
	runes := []rune(text)
	for i, r := range runes {
		if r > 0xFFFF { // temporary fix for Case429060
			if q.font.LookupRunes([]rune{r})[0] == 0 {
				r = '?'
				runes[i] = '?'
				repl = true
			}
		}

		q.usedRunes[r] = struct{}{}
	}
	q.usedRunesMux.Unlock()

	if repl {
		text = string(runes)
	}

	sn, _ := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder().String(text)
	return sn
}
func (q *CompositeFont) GetWidth(text string, fontSize float64) float64 {
	var w int
	for _, ind := range q.font.LookupRunes([]rune(text)) {
		w += q.font.GetGlyphAdvance(ind)
	}
	return float64(w) * fontSize / 1000
}
func (q *CompositeFont) GetAscent(fontSize float64) float64 {
	return float64(q.metrics.Ascent) * fontSize / 1000
}
func (q *CompositeFont) GetTop(fontSize float64) float64 {
	return float64(q.metrics.YMax) * fontSize / 1000
}
func (q *CompositeFont) GetBottom(fontSize float64) float64 {
	return float64(q.metrics.YMin) * fontSize / 1000
}
func (q *CompositeFont) GetHeight(fontSize float64) float64 {
	return float64(q.metrics.TextHeight) * fontSize / 1000
}
func (q *CompositeFont) GetUnderlineThickness(size float64) float64 {
	return float64(q.metrics.UnderlineThickness) * size / 1000
}
func (q *CompositeFont) GetUnderlinePosition(size float64) float64 {
	return float64(q.metrics.UnderlinePosition) * size / 1000
}
func (q *CompositeFont) HasGylph(runes []rune) []bool {
	dest := make([]bool, 0, len(runes))
	for _, ind := range q.font.LookupRunes(runes) {
		dest = append(dest, ind > 0)
	}
	return dest
}
func (q *CompositeFont) FallbackFont() FontHandler {
	return q.fallbackFont
}
func (q *CompositeFont) finish() error {
	if q.onFinish == nil {
		return nil
	}
	return q.onFinish()
}

// -------------------------------------------------------------------------------------------------------------------

// CompositeFontOTF references a composite font and provides additional function like font metrics
type CompositeFontOTF struct {
	reference    types.Reference
	usedRunes    map[rune]struct{}
	usedRunesMux sync.Mutex
	onFinish     func() error
	font         *sfnt.Font
	metrics      font.Metrics
	bounds       fixed.Rectangle26_6
}

func (q *CompositeFontOTF) Reference() types.Reference {
	return q.reference
}
func (q *CompositeFontOTF) Encode(text string) string {
	q.usedRunesMux.Lock()
	for _, r := range text {
		q.usedRunes[r] = struct{}{}
	}
	q.usedRunesMux.Unlock()

	sn, _ := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder().String(text)
	return sn
}
func (q *CompositeFontOTF) GetWidth(text string, fontSize float64) float64 {
	var w int
	for _, r := range text {
		ind, err := q.font.GlyphIndex(nil, r)
		if err != nil {
			continue
		}
		gw, err := q.font.GlyphAdvance(nil, ind, fixed.I(1000), font.HintingNone)
		if err != nil {
			continue
		}
		w += gw.Round()
	}
	return float64(w) * fontSize / 1000
}
func (q *CompositeFontOTF) GetAscent(fontSize float64) float64 {
	return float64(q.metrics.Ascent.Round()) * fontSize / 1000
}
func (q *CompositeFontOTF) GetTop(fontSize float64) float64 {
	return float64(-q.bounds.Min.Y.Round()) * fontSize / 1000
}
func (q *CompositeFontOTF) GetBottom(fontSize float64) float64 {
	return float64(-q.bounds.Max.Y.Round()) * fontSize / 1000
}
func (q *CompositeFontOTF) GetHeight(fontSize float64) float64 {
	return float64((q.bounds.Max.Y.Round() - q.bounds.Min.Y.Round())) * fontSize / 1000
}
func (q *CompositeFontOTF) GetUnderlineThickness(size float64) float64 {
	underlineThickness := 100
	return float64(underlineThickness) * size / 1000
}
func (q *CompositeFontOTF) GetUnderlinePosition(size float64) float64 {
	underlinePosition := -100
	return float64(underlinePosition) * size / 1000
}
func (q *CompositeFontOTF) HasGylph(runes []rune) []bool {
	dest := make([]bool, 0, len(runes))
	for _, r := range runes {
		ind, _ := q.font.GlyphIndex(nil, r)
		dest = append(dest, ind > 0)
	}
	return dest
}
func (q *CompositeFontOTF) FallbackFont() FontHandler {
	return nil
}
func (q *CompositeFontOTF) finish() error {
	if q.onFinish == nil {
		return nil
	}
	return q.onFinish()
}
