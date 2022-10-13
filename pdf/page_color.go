package pdf

import "github.com/raceresult/gopdf/types"

// PDF Reference 1.4, Table 4.21 Color operators

// Color_CS sets the current color space to use for stroking operations. The oper-
// and name must be a name object. If the color space is one that can be specified
// by a name and no additional parameters (DeviceGray, DeviceRGB, DeviceCMYK,
// and certain cases of Pattern), the name may be specified directly. Otherwise, it
// must be a name defined in the ColorSpace subdictionary of the current resource
// dictionary (see Section 3.7.2, “Resource Dictionaries”); the associated value is an
// array describing the color space (see Section 4.5.2, “Color Space Families”).
// No-Op if identical to value in current graphics state.
func (q *Page) Color_CS(name types.ColorSpaceFamily) {
	if q.graphicsState.StrokingColor.Name == name {
		return
	}
	q.graphicsState.StrokingColor.Name = name

	q.AddCommand("CS", name)
}

// Color_cs Same as CS, but for nonstroking operations.
// No-Op if identical to value in current graphics state.
func (q *Page) Color_cs(name types.ColorSpaceFamily) {
	if q.graphicsState.NonStrokingColor.Name == name {
		return
	}
	q.graphicsState.NonStrokingColor.Name = name

	q.AddCommand("cs", name)
}

// Color_SC sets the color to use for stroking operations in a device, CIE-based
// (other than ICCBased), or Indexed color space. The number of operands re-
// quired and their interpretation depends on the current stroking color space:
// • For DeviceGray, CalGray, and Indexed color spaces, one operand is required (n = 1).
// • For DeviceRGB, CalRGB, and Lab color spaces, three operands are required (n = 3).
// • For DeviceCMYK, four operands are required (n = 4).
// No-Op if identical to value in current graphics state.
func (q *Page) Color_SC(c ...float64) {
	if !q.graphicsState.StrokingColor.SetIfNotEqual(q.graphicsState.StrokingColor.Name, "", c...) {
		return
	}

	arr := make([]types.Object, 0, len(c))
	for _, v := range c {
		arr = append(arr, types.Number(v))
	}
	q.AddCommand("SC", arr...)
}

// Color_SCN Same as SC, but also supports Pattern, Separation, DeviceN, and ICCBased color spaces.
// If the current stroking color space is a Separation, DeviceN, or ICCBased color
// space, the operands c1 ...cn are numbers. The number of operands and their interpretation depends on the color space.
// If the current stroking color space is a Pattern color space, name is the name of
// an entry in the Pattern subdictionary of the current resource dictionary (see
// Section 3.7.2, “Resource Dictionaries”). For an uncolored tiling pattern
// (PatternType = 1 and PaintType = 2), c1 ...cn are component values specifying a
// color in the pattern’s underlying color space. For other types of pattern, these
// operands must not be specified.
// No-Op if identical to value in current graphics state.
func (q *Page) Color_SCN(name types.Name, c ...float64) {
	if !q.graphicsState.StrokingColor.SetIfNotEqual(q.graphicsState.StrokingColor.Name, name, c...) {
		return
	}

	arr := make([]types.Object, 0, len(c)+1)
	for _, v := range c {
		arr = append(arr, types.Number(v))
	}
	arr = append(arr, name)
	q.AddCommand("SCN", arr...)
}

// Color_sc : Same as SC, but for nonstroking operations.
// No-Op if identical to value in current graphics state.
func (q *Page) Color_sc(c ...float64) {
	if !q.graphicsState.NonStrokingColor.SetIfNotEqual(q.graphicsState.NonStrokingColor.Name, "", c...) {
		return
	}

	arr := make([]types.Object, 0, len(c))
	for _, v := range c {
		arr = append(arr, types.Number(v))
	}
	q.AddCommand("sc", arr...)
}

// Color_scn : Same as SCN, but for nonstroking operations.
// No-Op if identical to value in current graphics state.
func (q *Page) Color_scn(name types.Name, c ...float64) {
	if !q.graphicsState.NonStrokingColor.SetIfNotEqual(q.graphicsState.NonStrokingColor.Name, name, c...) {
		return
	}

	arr := make([]types.Object, 0, len(c)+1)
	for _, v := range c {
		arr = append(arr, types.Number(v))
	}
	arr = append(arr, name)
	q.AddCommand("scn", arr...)
}

// Color_G sets the stroking color space to DeviceGray (or the DefaultGray color space; see
// “Default Color Spaces” on page 194) and set the gray level to use for stroking
// operations. gray is a number between 0.0 (black) and 1.0 (white).
// No-Op if identical to value in current graphics state.
func (q *Page) Color_G(gray float64) {
	if !q.graphicsState.StrokingColor.SetIfNotEqual(types.ColorSpace_DeviceGray, "", gray) {
		return
	}

	q.AddCommand("G", types.Number(gray))
}

// Color_g : Same as G, but for nonstroking operations.
// No-Op if identical to value in current graphics state.
func (q *Page) Color_g(gray float64) {
	if !q.graphicsState.NonStrokingColor.SetIfNotEqual(types.ColorSpace_DeviceGray, "", gray) {
		return
	}

	q.AddCommand("g", types.Number(gray))
}

// Color_RG sets the stroking color space to DeviceRGB (or the DefaultRGB color space; see
// “Default Color Spaces” on page 194) and set the color to use for stroking opera-
// tions. Each operand must be a number between 0.0 (minimum intensity) and
// 1.0 (maximum intensity).
// No-Op if identical to value in current graphics state.
func (q *Page) Color_RG(r, g, b float64) {
	if !q.graphicsState.StrokingColor.SetIfNotEqual(types.ColorSpace_CalRGB, "", r, g, b) {
		return
	}

	q.AddCommand("RG", types.Number(r), types.Number(g), types.Number(b))
}

// Color_rg : Same as RG, but for nonstroking operations.
// No-Op if identical to value in current graphics state.
func (q *Page) Color_rg(r, g, b float64) {
	if !q.graphicsState.NonStrokingColor.SetIfNotEqual(types.ColorSpace_CalRGB, "", r, g, b) {
		return
	}

	q.AddCommand("rg", types.Number(r), types.Number(g), types.Number(b))
}

// Color_K sets the stroking color space to DeviceCMYK (or the DefaultCMYK color space; see
// “Default Color Spaces” on page 194) and set the color to use for stroking opera-
// tions. Each operand must be a number between 0.0 (zero concentration) and 1.0
// (maximum concentration). The behavior of this operator is affected by the over-
// print mode (see Section 4.5.6, “Overprint Control”).
// No-Op if identical to value in current graphics state.
func (q *Page) Color_K(c, m, y, k float64) {
	if !q.graphicsState.StrokingColor.SetIfNotEqual(types.ColorSpace_DeviceCMYK, "", c, m, y, k) {
		return
	}

	q.AddCommand("K", types.Number(c), types.Number(m), types.Number(y), types.Number(k))
}

// Color_k : Same as K, but for nonstroking operations.
// No-Op if identical to value in current graphics state.
func (q *Page) Color_k(c, m, y, k float64) {
	if !q.graphicsState.NonStrokingColor.SetIfNotEqual(types.ColorSpace_DeviceCMYK, "", c, m, y, k) {
		return
	}

	q.AddCommand("k", types.Number(c), types.Number(m), types.Number(y), types.Number(k))
}
