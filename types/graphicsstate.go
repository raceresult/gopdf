package types

// PDF Reference 1.4, Table 4.8 Entries in a graphics state parameter dictionary

type GraphicsState struct {
	// Type // Required

	// (Optional; PDF 1.3) The line width (see “Line Width” on page 152).
	W Number

	// (Optional; PDF 1.3) The line cap style (see “Line Cap Style” on page 153).
	LC Int

	// (Optional; PDF 1.3) The line join style (see “Line Join Style” on page 153).
	LJ Int

	// (Optional; PDF 1.3) The miter limit (see “Miter Limit” on page 153).
	ML Number

	// (Optional; PDF 1.3) The line dash pattern, expressed as an array of the form
	// [dashArray dashPhase], where dashArray is itself an array and dashPhase is an
	// 	integer (see “Line Dash Pattern” on page 155).
	D Array

	// (Optional; PDF 1.3) The name of the rendering intent (see “Rendering
	// Intents” on page 197).
	RI Number

	// (Optional) A flag specifying whether to apply overprint (see Section 4.5.6,
	// “Overprint Control”). In PDF 1.2 and earlier, there is a single overprint
	// parameter that applies to all painting operations. Beginning with PDF 1.3,
	// there are two separate overprint parameters: one for stroking and one for all
	// other painting operations. Specifying an OP entry sets both parameters un-
	// less there is also an op entry in the same graphics state parameter dictionary,
	// in which case the OP entry sets only the overprint parameter for stroking.
	OP Boolean

	// (Optional; PDF 1.3) A flag specifying whether to apply overprint (see Section
	// 4.5.6, “Overprint Control”) for painting operations other than stroking. If
	// this entry is absent, the OP entry, if any, sets this parameter.
	Op Boolean

	// (Optional; PDF 1.3) The overprint mode (see Section 4.5.6, “Overprint Control”).
	OPM Int

	// (Optional; PDF 1.3) An array of the form [font size], where font is an indirect
	// reference to a font dictionary and size is a number expressed in text space
	// units. These two objects correspond to the operands of the Tf operator (see
	// Section 5.2, “Text State Parameters and Operators”); however, the first oper-
	// and is an indirect object reference instead of a resource name.
	// BG function (Optional) The black-generation function, which maps the interval [0.0 1.0]
	// to the interval [0.0 1.0] (see Section 6.2.3, “Conversion from DeviceRGB to
	// DeviceCMYK”).
	Font Array

	// (Optional; PDF 1.3) Same as BG except that the value may also be the name
	// Default, denoting the black-generation function that was in effect at the start
	// of the page. If both BG and BG2 are present in the same graphics state param-
	// eter dictionary, BG2 takes precedence.
	BG2 Object

	// (Optional) The undercolor-removal function, which maps the interval
	// [0.0 1.0] to the interval [−1.0 1.0] (see Section 6.2.3, “Conversion from
	// DeviceRGB to DeviceCMYK”).
	UCR Function

	// (Optional; PDF 1.3) Same as UCR except that the value may also be the name
	// Default, denoting the undercolor-removal function that was in effect at the
	// start of the page. If both UCR and UCR2 are present in the same graphics state
	// parameter dictionary, UCR2 takes precedence.
	UCR2 Object

	// (Optional) The transfer function, which maps the interval [0.0 1.0] to the
	// or name interval [0.0 1.0] (see Section 6.3, “Transfer Functions”). The value is either
	// a single function (which applies to all process colorants) or an array of four
	// functions (which apply to the process colorants individually). The name
	// Identity may be used to represent the identity function.
	// Graphics StateS E C T I O N 4 . 3 159
	TR Object

	// (Optional; PDF 1.3) Same as TR except that the value may also be the name
	// or name Default, denoting the transfer function that was in effect at the start of the
	// page. If both TR and TR2 are present in the same graphics state parameter dic-
	// tionary, TR2 takes precedence.
	TR2 Object

	// (Optional) The halftone dictionary or stream (see Section 6.4, “Halftones”)
	// 	stream, or name or the name Default, denoting the halftone that was in effect at the start of the
	// 	page.
	HT Dictionary

	// (Optional; PDF 1.3) The flatness tolerance (see Section 6.5.1, “Flatness Tolerance”).
	FL Number

	// (Optional; PDF 1.3) The smoothness tolerance (see Section 6.5.2, “Smooth-
	// ness Tolerance”).
	SM Number

	// (Optional) A flag specifying whether to apply automatic stroke adjustment
	// (see Section 6.5.4, “Automatic Stroke Adjustment”).
	SA Boolean

	// (Optional; PDF 1.4) The current blend mode to be used in the transparent
	// imaging model (see Sections 7.2.4, “Blend Mode,” and 7.5.2, “Specifying
	// Blending Color Space and Blend Mode”).
	BM Object

	// (Optional; PDF 1.4) The current soft mask, specifying the mask shape or
	// mask opacity values to be used in the transparent imaging model (see
	// “Source Shape and Opacity” on page 421 and “Mask Shape and Opacity” on
	// page 443).
	// Note: Although the current soft mask is sometimes referred to as a “soft clip,”
	// altering it with the gs operator completely replaces the old value with the new
	// one, rather than intersecting the two as is done with the current clipping path
	// parameter (see Section 4.4.3, “Clipping Path Operators”).
	SMask Object

	// (Optional; PDF 1.4) The current stroking alpha constant, specifying the con-
	// stant shape or constant opacity value to be used for stroking operations in the
	// transparent imaging model (see “Source Shape and Opacity” on page 421
	// and “Constant Shape and Opacity” on page 444).
	CA Number

	// (Optional; PDF 1.4) Same as CA, but for nonstroking operations.
	Ca Number

	// (Optional; PDF 1.4) The alpha source flag (“alpha is shape”), specifying
	// whether the current soft mask and alpha constant are to be interpreted as
	// shape values (true) or opacity values (false).
	AIS Boolean

	// (Optional; PDF 1.4) The text knockout flag, which determines the behavior
	// of overlapping glyphs within a text object in the transparent imaging model
	// (see Section 5.2.7, “Text Knockout”
	TK Boolean
}

func (q GraphicsState) ToRawBytes() []byte {
	d := Dictionary{
		"Type": Name("ExtGState"),
		// todo
	}

	return d.ToRawBytes()
}

func (q GraphicsState) Copy(copyRef func(reference Reference) Reference) Object {
	return GraphicsState{
		W:     q.W.Copy(copyRef).(Number),
		LC:    q.LC.Copy(copyRef).(Int),
		LJ:    q.LJ.Copy(copyRef).(Int),
		ML:    q.ML.Copy(copyRef).(Number),
		D:     q.D.Copy(copyRef).(Array),
		RI:    q.RI.Copy(copyRef).(Number),
		OP:    q.OP.Copy(copyRef).(Boolean),
		Op:    q.Op.Copy(copyRef).(Boolean),
		OPM:   q.OPM.Copy(copyRef).(Int),
		Font:  q.Font.Copy(copyRef).(Array),
		BG2:   Copy(q.BG2, copyRef),
		UCR:   q.UCR.Copy(copyRef).(Function),
		UCR2:  Copy(q.UCR2, copyRef),
		TR:    Copy(q.TR, copyRef),
		TR2:   Copy(q.TR2, copyRef),
		HT:    q.HT.Copy(copyRef).(Dictionary),
		FL:    q.FL.Copy(copyRef).(Number),
		SM:    q.SM.Copy(copyRef).(Number),
		SA:    q.SA.Copy(copyRef).(Boolean),
		BM:    Copy(q.BM, copyRef),
		SMask: Copy(q.SMask, copyRef),
		CA:    q.CA.Copy(copyRef).(Number),
		Ca:    q.Ca.Copy(copyRef).(Number),
		AIS:   q.AIS.Copy(copyRef).(Boolean),
		TK:    q.TK.Copy(copyRef).(Boolean),
	}
}

func (q GraphicsState) Equal(obj Object) bool {
	a, ok := obj.(GraphicsState)
	if !ok {
		return false
	}
	if !Equal(q.W, a.W) {
		return false
	}
	if !Equal(q.LC, a.LC) {
		return false
	}
	if !Equal(q.LJ, a.LJ) {
		return false
	}
	if !Equal(q.ML, a.ML) {
		return false
	}
	if !Equal(q.D, a.D) {
		return false
	}
	if !Equal(q.RI, a.RI) {
		return false
	}
	if !Equal(q.OP, a.OP) {
		return false
	}
	if !Equal(q.Op, a.Op) {
		return false
	}
	if !Equal(q.OPM, a.OPM) {
		return false
	}
	if !Equal(q.Font, a.Font) {
		return false
	}
	if !Equal(q.BG2, a.BG2) {
		return false
	}
	if !Equal(q.UCR, a.UCR) {
		return false
	}
	if !Equal(q.UCR2, a.UCR2) {
		return false
	}
	if !Equal(q.TR, a.TR) {
		return false
	}
	if !Equal(q.TR2, a.TR2) {
		return false
	}
	if !Equal(q.HT, a.HT) {
		return false
	}
	if !Equal(q.FL, a.FL) {
		return false
	}
	if !Equal(q.SM, a.SM) {
		return false
	}
	if !Equal(q.SA, a.SA) {
		return false
	}
	if !Equal(q.BM, a.BM) {
		return false
	}
	if !Equal(q.SMask, a.SMask) {
		return false
	}
	if !Equal(q.CA, a.CA) {
		return false
	}
	if !Equal(q.Ca, a.Ca) {
		return false
	}
	if !Equal(q.AIS, a.AIS) {
		return false
	}
	if !Equal(q.TK, a.TK) {
		return false
	}
	return true
}
