package gopdf

import (
	"math"
	"strings"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// TextElement draws a text, may have line breaks
type TextElement struct {
	TextChunk
	Left, Top    Length
	TextAlign    HorizontalAlign
	LineHeight   float64
	Rotate       float64
	Transparency float64
}

// Build adds the element to the content stream
func (q *TextElement) Build(page *pdf.Page) (string, error) {
	// if no text given, ignore element
	if q.Text == "" {
		return "", nil
	}

	// set format of first font before saving graphics state
	warning, err := q.setFontAndColor(page)
	if err != nil {
		return warning, err
	}

	// graphics state
	page.GraphicsState_q()
	defer page.GraphicsState_Q()

	// set coordinate system
	y := float64(page.Data.MediaBox.URY) - q.Top.Pt()
	if q.Rotate == 0 {
		page.GraphicsState_cm(1, 0, 0, 1, q.Left.Pt(), y)
	} else {
		r := q.Rotate * math.Pi / 180
		page.GraphicsState_cm(math.Cos(r), math.Sin(r), -math.Sin(r), math.Cos(r), q.Left.Pt(), y)
	}

	// transparency
	if q.Transparency > 0 && q.Transparency <= 1 {
		n := page.AddExtGState(types.Dictionary{
			"ca": types.Number(1 - q.Transparency),
			"CA": types.Number(1 - q.Transparency),
		})
		page.GraphicsState_gs(n)
	}

	// calculate some values needed below
	lineHeight := q.lineHeight()
	top := 0.0

	// iterate over lines
	for _, line := range strings.Split(q.Text, "\n") {
		left := 0.0
		switch q.TextAlign {
		case HorizontalAlignCenter:
			left -= q.getLineWidth(line) / 2
		case HorizontalAlignRight:
			left -= q.getLineWidth(line)
		}

		warning2, err := q.draw(page, left, top)
		if err != nil {
			return "", err
		}
		if warning2 != "" {
			warning += warning2
		}

		top -= lineHeight
	}

	return warning, nil
}

// TextHeight returns the height of the text, accounting for line breaks
func (q *TextElement) TextHeight() Length {
	lines := strings.Count(q.Text, "\n") + 1
	return Pt(float64(lines) * q.lineHeight())
}

// FontHeight returns the height of the font (bounding box y min to max)
func (q *TextElement) FontHeight() Length {
	return Pt(q.Font.GetHeight(q.FontSize))
}

// getLineWidth returns the width of the given text line consider font, fontsize, text-scaling, char spacing
func (q *TextElement) getLineWidth(line string) float64 {
	if line == "" {
		return 0
	}
	v := q.Font.GetWidth(line, q.FontSize)
	if q.CharSpacing.Value != 0 {
		v += float64(len([]rune(line))-1) * q.CharSpacing.Pt()
	}
	if q.TextScaling != 0 {
		v *= q.TextScaling / 100
	}
	return v
}

func (q *TextElement) lineHeight() float64 {
	if q.LineHeight != 0 {
		return q.LineHeight
	}
	return q.FontHeight().Pt()
}
