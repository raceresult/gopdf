package types

import (
	"bytes"
	"errors"
)

// PDF Reference 1.4, Table 4.35 Additional entries specific to an image dictionary

type Image struct {
	Stream StreamObject

	// (Optional) The type of PDF object that this dictionary describes; if
	// present, must be XObject for an image XObject.
	// Type

	// (Required) The type of XObject that this dictionary describes; must be
	// Image for an image XObject.
	// Subtype  types.Name

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
	sb := bytes.Buffer{}
	d := q.Stream.Dictionary.createDict()
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
	sb.Write(q.Stream.Stream)
	sb.WriteString("\n")
	sb.WriteString("endstream\n")

	return sb.Bytes()
}

func (q *Image) Read(dict Dictionary) error {
	// Type
	v, ok := dict["Type"]
	if !ok {
		return errors.New("XObject missing Type")
	}
	dtype, ok := v.(Name)
	if !ok {
		return errors.New("XObject field Type invalid")
	}
	if dtype != "XObject" {
		return errors.New("unexpected value in XObject field Type")
	}

	// Subtype
	v, ok = dict["Subtype"]
	if !ok {
		return errors.New("XObject field Subtype missing")
	}
	vt, ok := v.(Name)
	if !ok {
		return errors.New("XObject field SubType invalid")
	}
	if vt != "Image" {
		return errors.New("XObject field SubType invalid")
	}

	// general stream dictionary
	if err := q.Stream.Dictionary.Read(dict); err != nil {
		return err
	}

	// Width
	v, ok = dict["Width"]
	if !ok {
		return errors.New("XObject field Width missing")
	}
	q.Width, ok = v.(Int)
	if !ok {
		return errors.New("XObject field Width invalid")
	}

	// Height
	v, ok = dict["Height"]
	if !ok {
		return errors.New("XObject field Height missing")
	}
	q.Height, ok = v.(Int)
	if !ok {
		return errors.New("XObject field Height invalid")
	}

	// ColorSpace
	v, ok = dict["ColorSpace"]
	if ok {
		q.ColorSpace = v
	}

	// BitsPerComponent
	v, ok = dict["BitsPerComponent"]
	if ok {
		q.BitsPerComponent, ok = v.(Int)
		if !ok {
			return errors.New("XObject field BitsPerComponent invalid")
		}
	}

	// Intent
	v, ok = dict["Intent"]
	if ok {
		q.Intent, ok = v.(String)
		if !ok {
			return errors.New("XObject field Intent invalid")
		}
	}

	// ImageMask
	v, ok = dict["ImageMask"]
	if ok {
		q.ImageMask, ok = v.(Boolean)
		if !ok {
			return errors.New("XObject field ImageMask invalid")
		}
	}

	// Mask
	v, ok = dict["Mask"]
	if ok {
		q.Mask = v
	}

	// SMask
	v, ok = dict["SMask"]
	if ok {
		q.SMask = v
	}

	// Decode
	v, ok = dict["Decode"]
	if ok {
		q.Decode, ok = v.(Array)
		if !ok {
			return errors.New("XObject field Decode invalid")
		}
	}

	// Interpolate
	v, ok = dict["Interpolate"]
	if ok {
		q.Interpolate, ok = v.(Boolean)
		if !ok {
			return errors.New("XObject field Interpolate invalid")
		}
	}

	// Alternates
	v, ok = dict["Alternates"]
	if ok {
		q.Alternates, ok = v.(Array)
		if !ok {
			return errors.New("XObject field Alternates invalid")
		}
	}

	// Name
	v, ok = dict["Name"]
	if ok {
		q.Name, ok = v.(Name)
		if !ok {
			return errors.New("XObject field Name invalid")
		}
	}

	// StructParent
	v, ok = dict["StructParent"]
	if ok {
		q.StructParent, ok = v.(Int)
		if !ok {
			return errors.New("XObject field StructParent invalid")
		}
	}

	// ID
	v, ok = dict["ID"]
	if ok {
		q.ID, ok = v.(String)
		if !ok {
			return errors.New("XObject field ID invalid")
		}
	}

	// OPI
	v, ok = dict["OPI"]
	if ok {
		q.OPI = v
	}

	// Metadata
	v, ok = dict["Metadata"]
	if ok {
		q.Metadata = v
	}

	// return without error
	return nil
}

func (q Image) Copy(copyRef func(reference Reference) Reference) Object {
	return Image{
		Stream:           q.Stream.Copy(copyRef).(StreamObject),
		Width:            q.Width.Copy(copyRef).(Int),
		Height:           q.Height.Copy(copyRef).(Int),
		ColorSpace:       Copy(q.ColorSpace, copyRef),
		BitsPerComponent: q.BitsPerComponent.Copy(copyRef).(Int),
		Intent:           q.Intent.Copy(copyRef).(String),
		ImageMask:        q.ImageMask.Copy(copyRef).(Boolean),
		Mask:             Copy(q.Mask, copyRef),
		SMask:            Copy(q.SMask, copyRef),
		Decode:           q.Decode.Copy(copyRef).(Array),
		Interpolate:      q.Interpolate.Copy(copyRef).(Boolean),
		Alternates:       q.Alternates.Copy(copyRef).(Array),
		Name:             q.Name.Copy(copyRef).(Name),
		StructParent:     q.StructParent.Copy(copyRef).(Int),
		ID:               q.ID.Copy(copyRef).(String),
		OPI:              Copy(q.OPI, copyRef),
		Metadata:         Copy(q.Metadata, copyRef),
	}
}
