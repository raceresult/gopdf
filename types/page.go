package types

import (
	"errors"
)

// PDF Reference 1.4, Table 3.18 Entries in a page object

type Page struct {
	// (Required) The type of PDF object that this dictionary describes; must be
	// Page for a page object.
	// Type

	// (Required; must be an indirect reference) The page tree node that is the im-
	// mediate parent of this page object.
	Parent Reference

	// Required if PieceInfo is present; optional otherwise; PDF 1.3) The date and
	// time (see Section 3.8.2, “Dates”) when the page’s contents were most re-
	// cently modified. If a page-piece dictionary (PieceInfo) is present, the
	// modification date is used to ascertain which of the application data dic-
	// tionaries that it contains correspond to the current content of the page
	// (see Section 9.4, “Page-Piece Dictionaries”).
	LastModified Date

	// (Required; inheritable) A dictionary containing any resources required by
	// the page (see Section 3.7.2, “Resource Dictionaries”). If the page requires
	// no resources, the value of this entry should be an empty dictionary; omit-
	// ting the entry entirely indicates that the resources are to be inherited from
	// an ancestor node in the page tree.
	Resources Object

	// (Required; inheritable) A rectangle (see Section 3.8.3, “Rectangles”), ex-
	// pressed in default user space units, defining the boundaries of the physical
	// medium on which the page is intended to be displayed or printed (see
	// Section 9.10.1, “Page Boundaries”).
	MediaBox Rectangle

	// (Optional; inheritable) A rectangle, expressed in default user space units,
	// defining the visible region of default user space. When the page is dis-
	// played or printed, its contents are to be clipped (cropped) to this rectangle
	// and then imposed on the output medium in some implementation-
	// defined manner (see Section 9.10.1, “Page Boundaries”). Default value:
	// the value of MediaBox.
	CropBox *Rectangle

	// (Optional; PDF 1.3) A rectangle, expressed in default user space units, de-
	// fining the region to which the contents of the page should be clipped
	// when output in a production environment (see Section 9.10.1, “Page
	// Boundaries”). Default value: the value of CropBox.
	BleedBox *Rectangle

	// (Optional; PDF 1.3) A rectangle, expressed in default user space units, de-
	// fining the intended dimensions of the finished page after trimming (see
	// Section 9.10.1, “Page Boundaries”). Default value: the value of CropBox.
	TrimBox *Rectangle

	// (Optional; PDF 1.3) A rectangle, expressed in default user space units, de-
	// fining the extent of the page’s meaningful content (including potential
	// white space) as intended by the page’s creator (see Section 9.10.1, “Page
	// Boundaries”). Default value: the value of CropBox.
	ArtBox *Rectangle

	// (Optional) A box color information dictionary specifying the colors and
	// other visual characteristics to be used in displaying guidelines on the
	// screen for the various page boundaries (see “Display of Page Boundaries”
	// on page 679). If this entry is absent, the viewer application should use its
	// own current default settings.
	BoxColorInfo Object

	// (Optional) A content stream (see Section 3.7.1, “Content Streams”) de-
	// scribing the contents of this page. If this entry is absent, the page is empty.
	// The value may be either a single stream or an array of streams. If it is an
	// array, the effect is as if all of the streams in the array were concatenated, in
	// order, to form a single stream. This allows a program generating a PDF
	// file to create image objects and other resources as they occur, even though
	// they interrupt the content stream. The division between streams may
	// occur only at the boundaries between lexical tokens (see Section 3.1, “Lex-
	// ical Conventions”), but is unrelated to the page’s logical content or orga-
	// nization. Applications that consume or produce PDF files are not required
	// to preserve the existing structure of the Contents array. (See implementa-
	// tion note 22 in Appendix H.)
	Contents Object

	// (Optional; inheritable) The number of degrees by which the page should
	// be rotated clockwise when displayed or printed. The value must be a mul-
	// tiple of 90. Default value: 0.
	Rotate Int

	// (Optional; PDF 1.4) A group attributes dictionary specifying the attributes
	//of the page’s page group for use in the transparent imaging model (see
	//Sections 7.3.6, “Page Group,” and 7.5.5, “Transparency Group XObjects”).
	Group Object

	// (Optional) A stream object defining the page’s thumbnail image (see Sec-
	// tion 8.2.3, “Thumbnail Images”).
	Thumb Object

	// (Optional; PDF 1.1; recommended if the page contains article beads) An ar-
	// ray of indirect references to article beads appearing on the page (see Sec-
	// tion 8.3.2, “Articles”; see also implementation note 23 in Appendix H).
	// The beads are listed in the array in natural reading order.
	B Array

	// (Optional; PDF 1.1) The page’s display duration (also called its advance
	// timing): the maximum length of time, in seconds, that the page will be
	// displayed during presentations before the viewer application automati-
	// cally advances to the next page (see Section 8.3.3, “Presentations”). By
	// default, the viewer does not advance automatically.
	Dur Number

	// (Optional; PDF 1.1) A transition dictionary describing the transition effect
	// to be used when displaying the page during presentations (see Section
	// 8.3.3, “Presentations”).
	Trans Object

	// (Optional) An array of annotation dictionaries representing annotations
	// associated with the page (see Section 8.4, “Annotations”).
	Annots Array

	// (Optional; PDF 1.2) An additional-actions dictionary defining actions to
	// be performed when the page is opened or closed (see Section 8.5.2, “Trig-
	// ger Events”; see also implementation note 24 in Appendix H).
	AA Object

	// (Optional; PDF 1.4) A metadata stream containing metadata for the page
	// (see Section 9.2.2, “Metadata Streams”).
	Metadata Object

	// (Optional; PDF 1.3) A page-piece dictionary associated with the page (see
	// Section 9.4, “Page-Piece Dictionaries”).
	PieceInfo Object

	// (Required if the page contains structural content items; PDF 1.3) The inte-
	// ger key of the page’s entry in the structural parent tree (see “Finding Struc-
	// ture Elements from Content Items” on page 600).
	StructParents Int

	// (Optional; PDF 1.3; indirect reference preferred) The digital identifier of the
	// page’s parent Web Capture content set (see Section 9.9.5, “Object At-
	// tributes Related to Web Capture”).
	ID String

	// (Optional; PDF 1.3) The page’s preferred zoom (magnification) factor: the
	// factor by which it should be scaled to achieve the “natural” display magni-
	// fication (see Section 9.9.5, “Object Attributes Related to Web Capture”).
	PZ Number

	// (Optional; PDF 1.3) A separation dictionary containing information need-
	// ed to generate color separations for the page (see Section 9.10.3, “Separa-
	// tion Dictionaries”).
	SeparationInfo Object
}

func (q Page) ToRawBytes() []byte {
	if q.Resources == nil {
		q.Resources = Dictionary{}
	}
	d := Dictionary{
		"Type":      Name("Page"),
		"Parent":    q.Parent,
		"MediaBox":  q.MediaBox,
		"Resources": q.Resources,
	}
	if !q.LastModified.IsZero() {
		d["LastModified"] = q.LastModified
	}
	if q.CropBox != nil {
		d["CropBox"] = q.CropBox
	}
	if q.BleedBox != nil {
		d["BleedBox"] = q.BleedBox
	}
	if q.TrimBox != nil {
		d["TrimBox"] = q.TrimBox
	}
	if q.ArtBox != nil {
		d["ArtBox"] = q.ArtBox
	}
	if q.BoxColorInfo != nil {
		d["BoxColorInfo"] = q.BoxColorInfo
	}
	if q.Contents != nil {
		d["Contents"] = q.Contents
	}
	if q.Rotate != 0 {
		d["Rotate"] = q.Rotate
	}
	if q.Group != nil {
		d["Group"] = q.Group
	}
	if q.Thumb != nil {
		d["Thumb"] = q.Thumb
	}
	if q.B != nil {
		d["B"] = q.B
	}
	if q.Dur != 0 {
		d["Dur"] = q.Dur
	}
	if q.Trans != nil {
		d["Trans"] = q.Trans
	}

	if len(q.Annots) != 0 {
		d["Annots"] = q.Annots
	}
	if q.AA != nil {
		d["AA"] = q.AA
	}
	if q.Metadata != nil {
		d["Metadata"] = q.Metadata
	}
	if q.PieceInfo != nil {
		d["PieceInfo"] = q.PieceInfo
	}
	if q.StructParents != 0 {
		d["StructParents"] = q.StructParents
	}
	if q.ID != "" {
		d["ID"] = q.ID
	}
	if q.PZ != 0 {
		d["PZ"] = q.PZ
	}
	if q.SeparationInfo != nil {
		d["SeparationInfo"] = q.SeparationInfo
	}

	return d.ToRawBytes()
}

func (q *Page) Read(dict Dictionary) error {
	// Type
	v, ok := dict["Type"]
	if !ok {
		return errors.New("page missing Type")
	}
	dtype, ok := v.(Name)
	if !ok {
		return errors.New("page field Type invalid")
	}
	if dtype != "Page" {
		return errors.New("unexpected value in page field Type")
	}

	// Parent
	v, ok = dict["Parent"]
	if !ok {
		return errors.New("page field Parent missing")
	}
	q.Parent, ok = v.(Reference)
	if !ok {
		return errors.New("page field Parent invalid")
	}

	// LastModified
	v, ok = dict["LastModified"]
	if ok {
		q.LastModified, ok = v.(Date)
		if !ok {
			return errors.New("page field LastModified invalid")
		}
	}

	// Resources
	v, ok = dict["Resources"]
	if !ok {
		return errors.New("page field Resource missing")
	}
	q.Resources = v

	// MediaBox
	v, ok = dict["MediaBox"]
	if !ok {
		return errors.New("page field MediaBox missing")
	}
	q.MediaBox, ok = v.(Rectangle)
	if !ok {
		va, ok := v.(Array)
		if !ok {
			return errors.New("page field MediaBox invalid")
		}
		if err := q.MediaBox.Read(va); err != nil {
			return err
		}
	}

	// CropBox
	v, ok = dict["CropBox"]
	if ok {
		r, ok := v.(Rectangle)
		if ok {
			q.CropBox = &r
		} else {
			va, ok := v.(Array)
			if !ok {
				return errors.New("page field CropBox invalid")
			}
			q.CropBox = &Rectangle{}
			if err := q.CropBox.Read(va); err != nil {
				return err
			}
		}
	}

	// BleedBox
	v, ok = dict["BleedBox"]
	if ok {
		r, ok := v.(Rectangle)
		if ok {
			q.BleedBox = &r
		} else {
			va, ok := v.(Array)
			if !ok {
				return errors.New("page field BleedBox invalid")
			}
			q.BleedBox = &Rectangle{}
			if err := q.BleedBox.Read(va); err != nil {
				return err
			}
		}
	}

	// TrimBox
	v, ok = dict["TrimBox"]
	if ok {
		r, ok := v.(Rectangle)
		if ok {
			q.TrimBox = &r
		} else {
			va, ok := v.(Array)
			if !ok {
				return errors.New("page field TrimBox invalid")
			}
			q.TrimBox = &Rectangle{}
			if err := q.TrimBox.Read(va); err != nil {
				return err
			}
		}
	}

	// ArtBox
	v, ok = dict["ArtBox"]
	if ok {
		r, ok := v.(Rectangle)
		if ok {
			q.ArtBox = &r
		} else {
			va, ok := v.(Array)
			if !ok {
				return errors.New("page field ArtBox invalid")
			}
			q.ArtBox = &Rectangle{}
			if err := q.ArtBox.Read(va); err != nil {
				return err
			}
		}
	}

	// BoxColorInfo
	v, ok = dict["BoxColorInfo"]
	if ok {
		q.BoxColorInfo = v
	}

	// Contents
	v, ok = dict["Contents"]
	if ok {
		q.Contents = v
	}

	// Rotate
	v, ok = dict["Rotate"]
	if ok {
		q.Rotate, ok = v.(Int)
		if !ok {
			return errors.New("page field Rotate invalid")
		}
	}

	// Group
	v, ok = dict["Group"]
	if ok {
		q.Group = v
	}

	// Thumb
	v, ok = dict["Thumb"]
	if ok {
		q.Thumb = v
	}

	// B
	v, ok = dict["B"]
	if ok {
		q.B, ok = v.(Array)
		if !ok {
			return errors.New("page field B invalid")
		}
	}

	// B
	v, ok = dict["B"]
	if ok {
		q.B, ok = v.(Array)
		if !ok {
			return errors.New("page field B invalid")
		}
	}

	// Dur
	v, ok = dict["Dur"]
	if ok {
		q.Dur, ok = v.(Number)
		if !ok {
			dur, ok := v.(Int)
			if !ok {
				return errors.New("page field Dur invalid")
			}
			q.Dur = Number(dur)
		}
	}

	// Trans
	v, ok = dict["Trans"]
	if ok {
		q.Trans = v
	}

	// Annots
	v, ok = dict["Annots"]
	if ok {
		q.Annots, ok = v.(Array)
		if !ok {
			return errors.New("page field Annots invalid")
		}
	}

	// AA
	v, ok = dict["AA"]
	if ok {
		q.AA = v
	}

	// Metadata
	v, ok = dict["Metadata"]
	if ok {
		q.Metadata = v
	}

	// PieceInfo
	v, ok = dict["PieceInfo"]
	if ok {
		q.PieceInfo = v
	}

	// StructParents
	v, ok = dict["StructParents"]
	if ok {
		q.StructParents, ok = v.(Int)
		if !ok {
			return errors.New("page field StructParents invalid")
		}
	}

	// ID
	v, ok = dict["ID"]
	if ok {
		q.ID, ok = v.(String)
		if !ok {
			return errors.New("page field ID invalid")
		}
	}

	// PZ
	v, ok = dict["PZ"]
	if ok {
		q.PZ, ok = v.(Number)
		if !ok {
			pz, ok := v.(Int)
			if !ok {
				return errors.New("page field PZ invalid")
			}
			q.PZ = Number(pz)
		}
	}

	// SeparationInfo
	v, ok = dict["SeparationInfo"]
	if ok {
		q.SeparationInfo = v
	}

	// return without error
	return nil
}
