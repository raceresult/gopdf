package types

import (
	"bytes"
)

// PDF Reference 1.4, Table 4.35 Additional entries specific to an image dictionary

type Image struct {
	Stream

	// (Optional) The type of PDF object that this dictionary describes; if
	// present, must be XObject for an image XObject.
	// Type

	// (Required) The type of XObject that this dictionary describes; must be
	// Image for an image XObject.
	//Subtype  types.Name

	// (Required) The width of the image, in samples.
	Width Int

	// (Required) The height of the image, in samples.
	Height Int

	// (Required except for image masks; not allowed for image masks) The color
	// space in which image samples are specified. This may be any type of color
	// space except Pattern.
	ColorSpace Object

	// (Required except for image masks; optional for image masks) The number of
	// bits used to represent each color component. Only a single value may be
	// specified; the number of bits is the same for all color components. Valid
	// values are 1, 2, 4, and 8. If ImageMask is true, this entry is optional, and if
	// specified, its value must be 1.
	// If the image stream uses a filter, the value of BitsPerComponent must be
	// consistent with the size of the data samples that the filter delivers. In par-
	// ticular, a CCITTFaxDecode or JBIG2Decode filter always delivers 1-bit sam-
	// ples, a RunLengthDecode or DCTDecode filter delivers 8-bit samples, and
	// an LZWDecode or FlateDecode filter delivers samples of a specified size if
	// a predictor function is used.
	BitsPerComponent Int

	// (Optional; PDF 1.1) The name of a color rendering intent to be used in
	//rendering the image (see “Rendering Intents” on page 197). Default value:
	//the current rendering intent in the graphics state.
	Intent String

	// (Optional) A flag indicating whether the image is to be treated as an image
	//mask (see Section 4.8.5, “Masked Images”). If this flag is true, the value of
	//BitsPerComponent must be 1 and Mask and ColorSpace should not be
	//specified; unmasked areas will be painted using the current nonstroking
	//color. Default value: false.
	ImageMask Boolean

	// (Optional except for image masks; not allowed for image masks; PDF 1.3) An
	// image XObject defining an image mask to be applied to this image (see
	// “Explicit Masking” on page 277), or an array specifying a range of colors
	// to be applied to it as a color key mask (see “Color Key Masking” on page
	// 277). If ImageMask is true, this entry must not be present. (See
	// implementation note 35 in Appendix H.)
	Mask Object

	// (Optional; PDF 1.4) A subsidiary image XObject defining a soft-mask
	// image (see “Soft-Mask Images” on page 447) to be used as a source of
	// mask shape or mask opacity values in the transparent imaging model. The
	// alpha source parameter in the graphics state determines whether the mask
	// values are interpreted as shape or opacity.
	// If present, this entry overrides the current soft mask in the graphics state,
	// as well as the image’s Mask entry, if any. (However, the other transparency-
	// related graphics state parameters—blend mode and alpha constant—
	// remain in effect.) If SMask is absent, the image has no associated soft mask
	// (although the current soft mask in the graphics state may still apply).
	SMask Object

	// (Optional) An array of numbers describing how to map image samples
	// into the range of values appropriate for the image’s color space (see
	// “Decode Arrays” on page 271). If ImageMask is true, the array must be
	// either [0 1] or [1 0]; otherwise, its length must be twice the number of
	// color components required by ColorSpace. Default value: see “Decode
	// Arrays” on page 271.
	Decode Array

	// (Optional) A flag indicating whether image interpolation is to be per-
	//formed (see “Image Interpolation” on page 273). Default value: false.
	Interpolate Boolean

	// (Optional; PDF 1.3) An array of alternate image dictionaries for this image
	// (see “Alternate Images” on page 273). The order of elements within the
	// array has no significance. This entry may not be present in an image
	// XObject that is itself an alternate image.
	Alternates Array

	// (Required in PDF 1.0; optional otherwise) The name by which this image
	// XObject is referenced in the XObject subdictionary of the current resource
	// dictionary (see Section 3.7.2, “Resource Dictionaries”).
	// Note: This entry is obsolescent and its use is no longer recommended. (See
	// implementation note 36 in Appendix H.)
	Name Name

	// (Required if the image is a structural content item; PDF 1.3) The integer key
	// of the image’s entry in the structural parent tree (see “Finding Structure
	// Elements from Content Items” on page 600).
	StructParent Int

	// (Optional; PDF 1.3; indirect reference preferred) The digital identifier of the
	//image’s parent Web Capture content set (see Section 9.9.5, “Object At-
	//tributes Related to Web Capture”).
	ID String

	// (Optional; PDF 1.2) An OPI version dictionary for the image (see Section
	// 9.10.6, “Open Prepress Interface (OPI)”). If ImageMask is true, this entry
	// is ignored.
	OPI Object

	// (Optional; PDF 1.4) A metadata stream containing metadata for the image
	// (see Section 9.2.2, “Metadata Streams”).
	Metadata Object
}

func (q Image) ToRawBytes() []byte {
	q.Data = bytes.TrimSpace(q.Data)

	sb := bytes.Buffer{}
	d := q.Stream.createDict()
	d["Type"] = Name("XObject")
	d["Subtype"] = Name("Image")
	d["Width"] = q.Width
	d["Height"] = q.Height
	d["ColorSpace"] = q.ColorSpace
	d["BitsPerComponent"] = q.BitsPerComponent

	if q.Intent != "" {
		d["Intent"] = Name(q.Intent)
	}
	if q.ImageMask {
		d["ImageMask"] = q.ImageMask
	}
	if q.Mask != nil {
		d["Mask"] = q.Mask
	}
	if q.SMask != nil {
		d["SMask"] = q.SMask
	}
	if len(q.Decode) != 0 {
		d["Decode"] = q.Decode
	}
	if q.Interpolate {
		d["Interpolate"] = q.Interpolate
	}
	if len(q.Alternates) != 0 {
		d["Alternates"] = q.Alternates
	}
	if q.Name != "" {
		d["Name"] = q.Name
	}
	if q.StructParent != 0 {
		d["StructParent"] = q.StructParent
	}
	if q.ID != "" {
		d["ID"] = q.ID
	}
	if q.OPI != nil {
		d["OPI"] = q.OPI
	}
	if q.Metadata != nil {
		d["Metadata"] = q.Metadata
	}

	sb.Write(d.ToRawBytes())

	sb.WriteString("stream\n")
	sb.Write(q.Data)
	sb.WriteString("\n")
	sb.WriteString("endstream\n")

	return sb.Bytes()
}
