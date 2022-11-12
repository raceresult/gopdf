package gopdf

import (
	"errors"
	"strconv"
	"strings"

	"github.com/raceresult/gopdf/pdf"
)

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

// Gray Color
// --------------------------------------------------------------------------------

// ParseColor parses a string to a color. Can be r,g,b or c,m,y,k or #RRGGBB
func ParseColor(s string) (Color, error) {
	if len(s) >= 7 && s[0] == '#' {
		r, err1 := strconv.ParseInt(s[1:3], 16, 64)
		g, err2 := strconv.ParseInt(s[3:5], 16, 64)
		b, err3 := strconv.ParseInt(s[5:7], 16, 64)
		if err1 == nil && err2 == nil && err3 == nil {
			return NewColorRGB(int(r), int(g), int(b)), nil
		}
	}

	arr := strings.Split(s, ",")
	switch len(arr) {
	case 1:
		n, err := strconv.Atoi(arr[0])
		if err == nil {
			r := n % 256
			g := (n / 256) % 256
			b := n / (256 * 256)
			return NewColorRGB(r, g, b), nil
		}

	case 3: // r,g,b
		r, err1 := strconv.Atoi(arr[0])
		g, err2 := strconv.Atoi(arr[1])
		b, err3 := strconv.Atoi(arr[2])
		if err1 == nil && err2 == nil && err3 == nil {
			if r >= 0 && r <= 255 && g >= 0 && g <= 255 && b >= 0 && b <= 255 {
				return NewColorRGB(r, g, b), nil
			}
		}

	case 4: // c,m,y,k
		c, err1 := strconv.Atoi(arr[0])
		m, err2 := strconv.Atoi(arr[1])
		y, err3 := strconv.Atoi(arr[2])
		k, err4 := strconv.Atoi(arr[3])
		if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
			if c >= 0 && c <= 100 && m >= 0 && m <= 100 && y >= 0 && y <= 100 && k >= 0 && k <= 100 {
				return NewColorCMYK(c, m, y, k), nil
			}
		}
	}
	return nil, errors.New("unknown color format")
}
