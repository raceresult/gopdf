package types

type StandardFontName Name

const (
	StandardFont_Helvetica            StandardFontName = "Helvetica"
	StandardFont_HelveticaBold        StandardFontName = "Helvetica-Bold"
	StandardFont_HelveticaOblique     StandardFontName = "Helvetica-Oblique"
	StandardFont_HelveticaBoldOblique StandardFontName = "Helvetica-BoldOblique"
	StandardFont_TimesRoman           StandardFontName = "Times-Roman"
	StandardFont_TimesBold            StandardFontName = "Times-Bold"
	StandardFont_TimesItalic          StandardFontName = "Times-Italic"
	StandardFont_TimesBoldItalic      StandardFontName = "Times-BoldItalic"
	StandardFont_Courier              StandardFontName = "Courier"
	StandardFont_CourierBold          StandardFontName = "Courier-Bold"
	StandardFont_CourierOblique       StandardFontName = "Courier-Oblique"
	StandardFont_CourierBoldOblique   StandardFontName = "Courier-BoldOblique"
	StandardFont_Symbol               StandardFontName = "Symbol"
	StandardFont_ZapfDingbats         StandardFontName = "ZapfDingbats"
)

func (q StandardFontName) ToRawBytes() []byte {
	return Name(q).ToRawBytes()
}

func (q StandardFontName) Copy(_ func(reference Reference) Reference) Object {
	return q
}
