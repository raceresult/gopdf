package types

// PDF Reference 1.4, Table 3.19 Entries in the name dictionary

type NameDictionary struct {
	// (Optional; PDF 1.2) A name tree mapping name strings to destinations (see
	// “Named Destinations” on page 476).
	Dests Object

	// (Optional; PDF 1.3) A name tree mapping name strings to annotation
	// appearance streams (see Section 8.4.4, “Appearance Streams”).
	AP Object

	// (Optional; PDF 1.3) A name tree mapping name strings to document-level
	// JavaScript® actions (see “JavaScript Actions” on page 556).
	JavaScript Object

	// (Optional; PDF 1.3) A name tree mapping name strings to visible pages for
	// use in interactive forms (see Section 8.6.5, “Named Pages”).
	Pages Object

	// (Optional; PDF 1.3) A name tree mapping name strings to invisible (tem-
	// plate) pages for use in interactive forms (see Section 8.6.5, “Named Pages”).
	Templates Object

	// (Optional; PDF 1.3) A name tree mapping digital identifiers to Web Capture
	// content sets (see Section 9.9.3, “Content Sets”).
	IDS Object

	// (Optional; PDF 1.3) A name tree mapping uniform resource locators (URLs)
	// to Web Capture content sets (see Section 9.9.3, “Content Sets”).
	URLS Object

	// (Optional; PDF 1.4) A name tree mapping name strings to embedded file
	// streams (see Section 3.10.3, “Embedded File Streams”).
	EmbeddedFiles Object
}

func (q NameDictionary) ToRawBytes() []byte {
	d := Dictionary{}
	if q.Dests != nil {
		d["Dests"] = q.Dests
	}
	if q.AP != nil {
		d["AP"] = q.AP
	}
	if q.JavaScript != nil {
		d["JavaScript"] = q.JavaScript
	}
	if q.Templates != nil {
		d["Templates"] = q.Templates
	}
	if q.IDS != nil {
		d["IDS"] = q.IDS
	}
	if q.URLS != nil {
		d["URLS"] = q.URLS
	}
	if q.EmbeddedFiles != nil {
		d["EmbeddedFiles"] = q.EmbeddedFiles
	}
	return d.ToRawBytes()
}

func (q NameDictionary) Copy(copyRef func(reference Reference) Reference) Object {
	return NameDictionary{
		Dests:         Copy(q.Dests, copyRef),
		AP:            Copy(q.AP, copyRef),
		JavaScript:    Copy(q.JavaScript, copyRef),
		Pages:         Copy(q.Pages, copyRef),
		Templates:     Copy(q.Templates, copyRef),
		IDS:           Copy(q.IDS, copyRef),
		URLS:          Copy(q.URLS, copyRef),
		EmbeddedFiles: Copy(q.EmbeddedFiles, copyRef),
	}
}
