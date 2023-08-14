package gopdf

import (
	"github.com/raceresult/gopdf/pdf"
)

// TextBoxElement is similar to TextElement, but can have a maximum width and height
type TextBoxElement struct {
	TextElement
	Width, Height   Length
	VerticalAlign   VerticalAlign
	HeightBufferRel float64
}

// Build adds the element to the content stream
func (q *TextBoxElement) Build(page *pdf.Page) (string, error) {
	return q.toChunkBox().Build(page)
}

// toChunkBox converts the element to a TextChunkBoxElement
func (q *TextBoxElement) toChunkBox() *TextChunkBoxElement {
	return &TextChunkBoxElement{
		Chunks: []TextChunk{{
			Text:          q.Text,
			Font:          q.Font,
			FontSize:      q.FontSize,
			Color:         q.Color,
			OutlineColor:  q.OutlineColor,
			OutlineWidth:  q.OutlineWidth,
			DashPattern:   q.DashPattern,
			Bold:          q.Bold,
			Italic:        q.Italic,
			Underline:     q.Underline,
			StrikeThrough: q.StrikeThrough,
			CharSpacing:   q.CharSpacing,
			TextScaling:   q.TextScaling,
		}},
		Transparency:    q.Transparency,
		LineHeight:      q.LineHeight,
		Left:            q.Left,
		Top:             q.Top,
		Width:           q.Width,
		Height:          q.Height,
		Rotate:          q.Rotate,
		TextAlign:       q.TextAlign,
		VerticalAlign:   q.VerticalAlign,
		HeightBufferRel: q.HeightBufferRel,
	}
}

// TextHeight returns the height of the text, accounting for line breaks and max width
func (q *TextBoxElement) TextHeight() Length {
	return q.toChunkBox().TextHeight()
}

// ShrinkToFit reduces the font size so that the entire text fits into the box
func (q *TextBoxElement) ShrinkToFit() {
	elem := q.toChunkBox()
	elem.ShrinkToFit()
	q.FontSize = elem.Chunks[0].FontSize
}
