package pdf

import (
	"bytes"

	"github.com/raceresult/gopdf/pdf/unitype"
	"github.com/raceresult/gopdf/types"
)

// NewStandardFont adds and returns a new standard font (expected to be available on all pdf consuming systems)
func (q *File) NewStandardFont(name types.StandardFontName, encoding types.Encoding) (*StandardFont, error) {
	f := types.StandardFont{
		BaseFont: name,
		Encoding: encoding,
	}
	metrics, err := f.Metrics()
	if err != nil {
		return nil, err
	}
	fh := &StandardFont{
		reference: q.creator.AddObject(f),
		encoding:  encoding,
		metrics:   metrics,
	}
	q.fonts = append(q.fonts, fh)
	return fh, nil
}

// NewTrueTypeFont adds and returns a new true type font
func (q *File) NewTrueTypeFont(ttf []byte, encoding types.Encoding, embed bool) (*TrueTypeFont, error) {
	// parse font by unitype
	font, err := unitype.Parse(bytes.NewReader(ttf))
	if err != nil {
		return nil, err
	}
	metrics := font.GetMetrics()

	// embed font
	var fontFileRef types.Reference
	if embed {
		fontFileRef = q.creator.AddObject(types.StreamFont{
			Stream: types.Stream{
				Data: ttf,
			},
			Length1: types.Int(len(ttf)),
		})
	}

	// create font descriptor
	var flags int
	if metrics.IsFixedPitch {
		flags += 1 << 0
	}
	if metrics.ItalicAngle != 0 {
		flags += 1 << 6
	}
	flags += 1 << 5
	fd := types.FontDescriptor{
		FontName: types.Name(font.GetNameByID(1)),
		Flags:    types.Int(flags),
		FontBBox: types.Rectangle{
			LLX: types.Number(metrics.XMin),
			LLY: types.Number(metrics.YMin),
			URX: types.Number(metrics.XMax),
			URY: types.Number(metrics.YMax),
		},
		ItalicAngle:  types.Number(metrics.ItalicAngle),
		Ascent:       types.Number(metrics.Ascent),
		Descent:      types.Number(metrics.Descent),
		CapHeight:    types.Number(metrics.CapHeight),
		XHeight:      types.Number(metrics.XHeight),
		StemV:        types.Number(calcStemV(metrics.Weight)),
		MissingWidth: types.Number(font.GetGlyphAdvance(0)),
		FontFile2:    fontFileRef,
	}
	fdRef := q.creator.AddObject(fd)

	// create font
	f := types.Font{
		Subtype:        types.FontSub_TrueType,
		BaseFont:       fd.FontName,
		Encoding:       encoding,
		FirstChar:      32,
		LastChar:       255,
		FontDescriptor: fdRef,
	}
	var widths types.Array
	for i := f.FirstChar; i <= f.LastChar; i++ {
		index := font.LookupRunes([]rune{rune(i)})
		w := font.GetGlyphAdvance(index[0])
		widths = append(widths, types.Int(w))
	}
	f.Widths = widths

	// create TrueTypeFont object
	fh := &TrueTypeFont{
		reference: q.creator.AddObject(f),
		encoding:  encoding,
		font:      font,
		metrics:   metrics,
	}
	q.fonts = append(q.fonts, fh)
	return fh, nil
}

// NewCompositeFont creates a new composite front from the given true type font
func (q *File) NewCompositeFont(ttf []byte) (*CompositeFont, error) {
	// parse font by unitype
	font, err := unitype.Parse(bytes.NewReader(ttf))
	if err != nil {
		return nil, err
	}
	metrics := font.GetMetrics()

	// create font descriptor
	var flags int
	if metrics.IsFixedPitch {
		flags += 1 << 0
	}
	if metrics.ItalicAngle != 0 {
		flags += 1 << 6
	}
	flags += 1 << 5
	fd := types.FontDescriptor{
		FontName: types.Name(font.GetNameByID(1)),
		Flags:    types.Int(flags),
		FontBBox: types.Rectangle{
			LLX: types.Number(metrics.XMin),
			LLY: types.Number(metrics.YMin),
			URX: types.Number(metrics.XMax),
			URY: types.Number(metrics.YMax),
		},
		ItalicAngle:  types.Number(metrics.ItalicAngle),
		Ascent:       types.Number(metrics.Ascent),
		Descent:      types.Number(metrics.Descent),
		CapHeight:    types.Number(metrics.CapHeight),
		XHeight:      types.Number(metrics.XHeight),
		StemV:        types.Number(calcStemV(metrics.Weight)),
		MissingWidth: types.Number(font.GetGlyphAdvance(0)),
	}
	fdRef := q.creator.AddObject(&fd)

	// create CID font
	cid := &types.CIDFont{
		Subtype:        types.FontSub_CIDFontType2,
		BaseFont:       fd.FontName,
		CIDSystemInfo:  q.getCIDSystemInfo(),
		FontDescriptor: fdRef,
		DW:             types.Int(750),
	}
	cidRef := q.creator.AddObject(cid)
	f := types.Type0Font{
		BaseFont:        cid.BaseFont,
		Encoding:        types.Name("Identity-H"),
		DescendantFonts: types.Array{cidRef},
		ToUnicode:       q.getToUnicode(),
	}

	// create CompositeFont object
	fh := &CompositeFont{
		reference: q.creator.AddObject(f),
		usedRunes: make(map[rune]struct{}),
		font:      font,
		metrics:   font.GetMetrics(),
	}
	fh.onFinish = func() error {
		// determine highest rune number
		var maxRune rune
		runes := make([]rune, 0, len(fh.usedRunes))
		for r := range fh.usedRunes {
			runes = append(runes, r)
			if maxRune < r {
				maxRune = r
			}
		}

		// add subsetted font to document
		indices := font.LookupRunes(runes)
		newFont, err := font.SubsetKeepIndices(indices)
		if err != nil {
			return err
		}
		var bts bytes.Buffer
		if err := newFont.Write(&bts); err != nil {
			return err
		}
		fd.FontFile2 = q.creator.AddObject(types.StreamFont{
			Stream: types.Stream{
				Data: bts.Bytes(),
			},
			Length1: types.Int(bts.Len()),
		})

		// build widths and cidToGID arrays
		cidToGID := make([]byte, maxRune*2+5)
		widths := types.Array{}
		for i, r := range runes {
			pos := indices[i]
			cidToGID[r*2] = byte(pos / 256)
			cidToGID[r*2+1] = byte(pos % 256)

			w := font.GetGlyphAdvance(pos)
			widths = append(widths, types.Int(r), types.Int(r), types.Int(w))
		}
		cid.CIDToGIDMap = q.creator.AddObject(types.Stream{Data: cidToGID})
		cid.W = widths
		return nil
	}
	q.fonts = append(q.fonts, fh)
	return fh, nil
}

// copied from fpdf.php
func (q *File) getCIDSystemInfo() types.Reference {
	if q.cidSystemInfo.Number == 0 {
		q.cidSystemInfo = q.creator.AddObject(types.CIDSystemInfoDictionary{
			Registry:   "Adobe",
			Ordering:   "UCS",
			Supplement: 0,
		})
	}
	return q.cidSystemInfo
}

// copied from fpdf.php
func (q *File) getToUnicode() types.Reference {
	if q.toUnicode.Number == 0 {
		q.toUnicode = q.creator.AddObject(types.Stream{Data: []byte("/CIDInit /ProcSet findresource begin\n12 dict begin\nbegincmap\n/CIDSystemInfo\n<</Registry (Adobe)\n/Ordering (UCS)\n/Supplement 0\n>> def\n/CMapName /Adobe-Identity-UCS def\n/CMapType 2 def\n1 begincodespacerange\n<0000> <FFFF>\nendcodespacerange\n1 beginbfrange\n<0000> <FFFF> <0000>\nendbfrange\nendcmap\nCMapName currentdict /CMap defineresource pop\nend\nend")})
	}
	return q.toUnicode
}

func calcStemV(weight int) int {
	f := float64(weight) / 65
	return int(50 + (f * f) + 0.5)
}
