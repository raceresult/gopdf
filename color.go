package gopdf

import "github.com/raceresult/gopdf/pdf"

// Color is the interface any type of color needs to fulfill
type Color interface {
	Build(page *pdf.Page, stroke bool)
}

// RGB Color
// --------------------------------------------------------------------------------

// ColorRGB represents a RGB color value
type ColorRGB struct {
	R int
	G int
	B int
}

// NewColorRGB creates a new ColorRGB object
func NewColorRGB(r, g, b int) ColorRGB {
	return ColorRGB{
		R: r,
		G: g,
		B: b,
	}
}

// Build sets the color in the graphics state of the given page
func (q ColorRGB) Build(page *pdf.Page, stroke bool) {
	if stroke {
		page.Color_RG(float64(q.R)/256, float64(q.G)/256, float64(q.B)/256)
	} else {
		page.Color_rg(float64(q.R)/256, float64(q.G)/256, float64(q.B)/256)
	}
}

var ColorRGBBlack = ColorRGB{R: 0, G: 0, B: 0}

// CMYK Color
// --------------------------------------------------------------------------------

// ColorCMYK represents a CMYK color value
type ColorCMYK struct {
	C int
	M int
	Y int
	K int
}

// NewColorCMYK creates a new ColorCMYK object
func NewColorCMYK(c, m, y, k int) ColorCMYK {
	return ColorCMYK{
		C: c,
		M: m,
		Y: y,
		K: k,
	}
}

// Build sets the color in the graphics state of the given page
func (q ColorCMYK) Build(page *pdf.Page, stroke bool) {
	if stroke {
		page.Color_K(float64(q.C)/100, float64(q.M)/100, float64(q.Y)/100, float64(q.K)/100)
	} else {
		page.Color_k(float64(q.C)/100, float64(q.M)/100, float64(q.Y)/100, float64(q.K)/100)
	}
}

// Gray Color
// --------------------------------------------------------------------------------

// ColorGray represents a gray color value
type ColorGray struct {
	Gray int
}

// NewColorGray creates a new ColorGray object
func NewColorGray(gray int) ColorGray {
	return ColorGray{
		Gray: gray,
	}
}

// Build sets the color in the graphics state of the given page
func (q ColorGray) Build(page *pdf.Page, stroke bool) {
	if stroke {
		page.Color_G(float64(q.Gray) / 255)
	} else {
		page.Color_g(float64(q.Gray) / 255)
	}
}
