// Copyright ©2021 The star-tex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package afm

import (
	"fmt"
	"io"

	"github.com/raceresult/gopdf/types/standardfont/fixed"
)

type Direction struct {
	// UnderlinePosition is the distance from the baseline for centering
	// underlining strokes.
	UnderlinePosition fixed.Int16_16

	// UnderlineThickness is the stroke width for underlining.
	UnderlineThickness fixed.Int16_16

	// ItalicAngle is the angle (in degrees counter-clockwise from the vertical)
	// of the dominant vertical stroke of the font.
	ItalicAngle fixed.Int16_16

	// CharWidth is the width vector of this font's program characters.
	CharWidth CharWidth

	// IsFixedPitch indicates whether the program is a fixed pitch (monospace) font.
	IsFixedPitch bool
}

type CharWidth struct {
	x fixed.Int16_16 // x component of the width vector of a font's program characters.
	y fixed.Int16_16 // y component of the width vector of a font's program characters.
}

type charMetric struct {
	// c is the decimal value of default character code.
	// c is -1 if the character is not encoded.
	c int

	// name is the PostScript name of this character.
	name string

	// w0 is the character width vector for writing Direction 0.
	w0 CharWidth

	// w1 is the character width vector for writing Direction 1.
	w1 CharWidth

	// vvector holds the components of a vector from origin 0 to origin 1.
	// origin 0 is the origin for writing Direction 0.
	// origin 1 is the origin for writing Direction 1.
	vv [2]fixed.Int16_16

	// bbox is the character bounding box.
	bbox BBox

	// ligs is a ligature sequence.
	ligs []lig
}

type BBox struct {
	LLX, LLY fixed.Int16_16
	URX, URY fixed.Int16_16
}

// lig is a ligature.
type lig struct {
	// succ is the name of the successor
	succ string
	// name is the name of the composite ligature, consisting
	// of the current character and the successor.
	name string
}

// Font is an Adobe Font metrics.
type Font struct {
	// MetricsSets defines the writing Direction.
	// 0: Direction 0 only.
	// 1: Direction 1 only.
	// 2: both directions.
	MetricsSets int

	FontName   string // FontName is the name of the font program as presented to the PostScript language 'findfont' operator.
	FullName   string // FullName is the full text name of the font.
	FamilyName string // FamilyName is the name of the typeface family to which the font belongs.
	Weight     string // Weight is the Weight of the font (ex: Regular, Bold, Light).
	BBox       BBox   // BBox is the font bounding box.
	Version    string // Version is the font program Version identifier.
	Notice     string // Notice contains the font name trademark or copyright Notice.

	// EncodingScheme specifies the default encoding vector used for this font
	// program (ex: AdobeStandardEncoding, JIS12-88-CFEncoding, ...)
	// Special font program might state FontSpecific.
	EncodingScheme string
	MappingScheme  int
	EscChar        int
	CharacterSet   string // CharacterSet describes the character set (glyph complement) of this font program.
	Characters     int    // Characters describes the number of Characters defined in this font program.
	IsBaseFont     bool   // IsBaseFont indicates whether this font is a base font program.

	// vvector holds the components of a vector from origin 0 to origin 1.
	// origin 0 is the origin for writing Direction 0.
	// origin 1 is the origin for writing Direction 1.
	// vvector is required when metricsSet is 2.
	vvector [2]fixed.Int16_16

	IsFixedV  bool // IsFixedV indicates whether vvector is the same for every character in this font.
	IsCIDFont bool // IsCIDFont indicates whether the font is a CID-keyed font.

	CapHeight fixed.Int16_16 // CapHeight is usually the y-value of the top of the capital 'H'.
	XHeight   fixed.Int16_16 // XHeight is typically the y-value of the top of the lowercase 'x'.
	Ascender  fixed.Int16_16 // Ascender is usually the y-value of the top of the lowercase 'd'.
	Descender fixed.Int16_16 // Descender is typically the y-value of the bottom of the lowercase 'p'.
	StdHW     fixed.Int16_16 // StdHW specifies the dominant width of horizontal stems.
	StdVW     fixed.Int16_16 // StdVW specifies the dominant width of vertical stems.

	blendAxisTypes       []string
	blendDesignPositions [][]fixed.Int16_16
	blendDesignMap       [][][]fixed.Int16_16
	weightVector         []fixed.Int16_16

	Direction   [3]Direction
	charMetrics []charMetric
	composites  []composite

	tkerns []trackKern
	pkerns []kernPair
}

func newFont() Font {
	return Font{
		IsBaseFont: true,
	}
}

// Parse parses an AFM file.
func Parse(r io.Reader) (Font, error) {
	var (
		fnt = newFont()
		p   = newParser(r)
	)
	err := p.parse(&fnt)
	if err != nil {
		return fnt, fmt.Errorf("could not parse AFM file: %w", err)
	}
	return fnt, nil
}

func (q Font) GetGlyphAdvance(charcode int) int {
	for _, c := range q.charMetrics {
		if c.c == charcode {
			return int(c.w0.x.Float64())
		}
	}
	return 0
}
