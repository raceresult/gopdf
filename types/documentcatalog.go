package types

import "errors"

// PDF Reference 1.4, Table 3.16 Entries in the catalog dictionary

type DocumentCatalog struct {
	// (Required) The type of PDF object that this dictionary describes; must
	// be Catalog for the catalog dictionary.
	// Type

	// (Optional; PDF 1.4) The version of the PDF specification to which the
	// document conforms (for example, 1.4), if later than the version specified
	// in the file’s header (see Section 3.4.1, “File Header”). If the header speci-
	// fies a later version, or if this entry is absent, the document conforms to
	// the version specified in the header. This entry enables a PDF producer
	// application to update the version using an incremental update; see Sec-
	// tion 3.4.5, “Incremental Updates.” (See implementation note 18 in Ap-
	// pendix H.)
	// Note: The value of this entry is a name object, not a number, and so must
	// be preceded by a slash character (/) when written in the PDF file (for ex-
	// ample, /1.4).
	Version Name

	// (Required; must be an indirect reference) The page tree node that is the
	// root of the document’s page tree (see Section 3.6.2, “Page Tree”).
	Pages Reference

	// (Optional; PDF 1.3) A number tree (see Section 3.8.5, “Number Trees”)
	// defining the page labeling for the document. The keys in this tree are
	// page indices; the corresponding values are page label dictionaries (see
	// Section 8.3.1, “Page Labels”). Each page index denotes the first page in a
	// labeling range to which the specified page label dictionary applies. The
	// tree must include a value for page index 0.
	PageLabels Object

	// (Optional; PDF 1.2) The document’s name dictionary (see Section 3.6.3,
	// “Name Dictionary”).
	Names Object

	// (Optional; PDF 1.1; must be an indirect reference) A dictionary of names
	// and corresponding destinations (see “Named Destinations” on page
	// 476).
	Dests Reference

	// (Optional; PDF 1.2) A viewer preferences dictionary (see Section 8.1,
	// “Viewer Preferences”) specifying the way the document is to be dis-
	// played on the screen. If this entry is absent, viewer applications should
	// use their own current user preference settings.
	ViewerPreferences Object

	// (Optional) A name object specifying the page layout to be used when the
	// document is opened:
	// SinglePage Display one page at a time.
	// OneColumn Display the pages in one column.
	// TwoColumnLeft Display the pages in two columns, with odd-
	// numbered pages on the left.
	// TwoColumnRight Display the pages in two columns, with odd-
	// numbered pages on the right.
	// (See implementation note 19 in Appendix H.) Default value: SinglePage.
	PageLayout Name

	// (Optional) A name object specifying how the document should be dis-
	// played when opened:
	// UseNone Neither document outline nor thumbnail im-
	// ages visible
	// UseOutlines Document outline visible
	// UseThumbs Thumbnail images visible
	// FullScreen Full-screen mode, with no menu bar, window
	// controls, or any other window visible
	// Default value: UseNone.
	PageMode Name

	// (Optional; must be an indirect reference) The outline dictionary that is the
	// root of the document’s outline hierarchy (see Section 8.2.2, “Document
	// Outline”).
	Outlines Reference

	// (Optional; PDF 1.1; must be an indirect reference) An array of thread
	// dictionaries representing the document’s article threads (see Section
	// 8.3.2, “Articles”).
	Threads Array

	// (Optional; PDF 1.1) A value specifying a destination to be displayed or
	// an action to be performed when the document is opened. The value is
	// either an array defining a destination (see Section 8.2.1, “Destinations”)
	// or an action dictionary representing an action (Section 8.5, “Actions”). If
	// this entry is absent, the document should be opened to the top of the
	// first page at the default magnification factor.
	OpenAction Object // can be reference as well

	// (Optional; PDF 1.1) A value specifying a destination to be displayed or
	// dictionary an action to be performed when the document is opened. The value is
	// either an array defining a destination (see Section 8.2.1, “Destinations”)
	// or an action dictionary representing an action (Section 8.5, “Actions”). If
	// this entry is absent, the document should be opened to the top of the
	// first page at the default magnification factor.
	AdditionalActions Object

	// (Optional) A URI dictionary containing document-level information for
	// URI (uniform resource identifier) actions (see “URI Actions” on page
	// 523).
	URI Object

	// (Optional; PDF 1.2) The document’s interactive form (AcroForm) dic-
	// tionary (see Section 8.6.1, “Interactive Form Dictionary”).
	AcroForm Object

	// (Optional; PDF 1.4; must be an indirect reference) A metadata stream
	// containing metadata for the document (see Section 9.2.2, “Metadata
	// Streams”).
	Metadata Reference

	// (Optional; PDF 1.3) The document’s structure tree root dictionary (see
	// Section 9.6.1, “Structure Hierarchy”).
	StructTreeRoot Object

	// (Optional; PDF 1.4) A mark information dictionary containing informa-
	// tion about the document’s usage of Tagged PDF conventions (see Sec-
	// tion 9.7.1, “Mark Information Dictionary”).
	MarkInfo Object

	// (Optional; PDF 1.4) A language identifier specifying the natural language
	// for all text in the document except where overridden by language speci-
	// fications for structure elements or marked content (see Section 9.8.1,
	// “Natural Language Specification”). If this entry is absent, the language is
	// considered unknown.
	Lang String

	// (Optional; PDF 1.3) A Web Capture information dictionary containing
	// state information used by the Acrobat Web Capture (AcroSpider) plug-
	// in extension (see Section 9.9.1, “Web Capture Information Dictionary”).
	SpiderInfo Object

	// (Optional; PDF 1.4) An array of output intent dictionaries describing the
	// color characteristics of output devices on which the document might be
	// rendered (see “Output Intents” on page 684).
	OutputIntents Object // seems that it can also be a refernece
}

func (q DocumentCatalog) ToRawBytes() []byte {
	d := Dictionary{
		"Type":  Name("Catalog"),
		"Pages": q.Pages,
	}

	if q.Version != "" {
		d["Version"] = q.Version
	}
	if q.PageLabels != nil {
		d["PageLabels"] = q.PageLabels
	}
	if q.Names != nil {
		d["Names"] = q.Names
	}
	if q.Dests.Number > 0 {
		d["Dests"] = q.Dests
	}
	if q.ViewerPreferences != nil {
		d["ViewerPreferences"] = q.ViewerPreferences
	}
	if q.PageLayout != "" {
		d["PageLayout"] = q.PageLayout
	}
	if q.PageMode != "" {
		d["PageMode"] = q.PageMode
	}
	if q.Outlines.Number > 0 {
		d["Outlines"] = q.Outlines
	}
	if q.Threads != nil {
		d["Threads"] = q.Threads
	}
	if q.OpenAction != nil {
		d["OpenAction"] = q.OpenAction
	}
	if q.AdditionalActions != nil {
		d["AA"] = q.AdditionalActions
	}
	if q.URI != nil {
		d["URI"] = q.URI
	}
	if q.AcroForm != nil {
		d["AcroForm"] = q.AcroForm
	}
	if q.Metadata.Number > 0 {
		d["Metadata"] = q.Metadata
	}
	if q.StructTreeRoot != nil {
		d["StructTreeRoot"] = q.StructTreeRoot
	}
	if q.MarkInfo != nil {
		d["MarkInfo"] = q.MarkInfo
	}
	if q.Lang != "" {
		d["Lang"] = q.Lang
	}
	if q.SpiderInfo != nil {
		d["SpiderInfo"] = q.SpiderInfo
	}
	if q.OutputIntents != nil {
		d["OutputIntents"] = q.OutputIntents
	}

	return d.ToRawBytes()
}

func (q *DocumentCatalog) Read(dict Dictionary) error {
	// Type
	v, ok := dict["Type"]
	if !ok {
		return errors.New("catalog missing Type")
	}
	dtype, ok := v.(Name)
	if !ok {
		return errors.New("catalog field Type invalid")
	}
	if dtype != "Catalog" {
		return errors.New("unexpected value in catalog field Type")
	}

	// Pages
	v, ok = dict["Pages"]
	if !ok {
		return errors.New("catalog field Pages missing")
	}
	pages, ok := v.(Reference)
	if !ok {
		return errors.New("catalog field Pages invalid")
	}
	q.Pages = pages

	// Version
	v, ok = dict["Version"]
	if ok {
		vt, ok := v.(Name)
		if !ok {
			return errors.New("unexpected value type in catalog field Value")
		}
		q.Version = vt
	}

	// PageLabels
	v, ok = dict["PageLabels"]
	if ok {
		q.PageLabels = v
	}

	// Names
	v, ok = dict["Names"]
	if ok {
		q.Names = v
	}

	// Dests
	v, ok = dict["Dests"]
	if ok {
		vt, ok := v.(Reference)
		if !ok {
			return errors.New("unexpected value type in catalog field Dests")
		}
		q.Dests = vt
	}

	// ViewerPreferences
	v, ok = dict["ViewerPreferences"]
	if ok {
		q.ViewerPreferences = v
	}

	// PageLayout
	v, ok = dict["PageLayout"]
	if ok {
		vt, ok := v.(Name)
		if !ok {
			return errors.New("unexpected value type in catalog field PageLayout")
		}
		q.PageLayout = vt
	}

	// PageMode
	v, ok = dict["PageMode"]
	if ok {
		vt, ok := v.(Name)
		if !ok {
			return errors.New("unexpected value type in catalog field PageMode")
		}
		q.PageMode = vt
	}

	// Version
	v, ok = dict["Outlines"]
	if ok {
		vt, ok := v.(Reference)
		if !ok {
			return errors.New("unexpected value type in catalog field Outlines")
		}
		q.Outlines = vt
	}

	// Threads
	v, ok = dict["Threads"]
	if ok {
		vt, ok := v.(Array)
		if !ok {
			return errors.New("unexpected value type in catalog field Threads")
		}
		q.Threads = vt
	}

	// OpenAction
	v, ok = dict["OpenAction"]
	if ok {
		q.OpenAction = v
	}

	// AdditionalActions
	v, ok = dict["AdditionalActions"]
	if ok {
		q.AdditionalActions = v
	}

	// URI
	v, ok = dict["URI"]
	if ok {
		q.URI = v
	}

	// AcroForm
	v, ok = dict["AcroForm"]
	if ok {
		q.AcroForm = v
	}

	// Metadata
	v, ok = dict["Metadata"]
	if ok {
		vt, ok := v.(Reference)
		if !ok {
			return errors.New("unexpected value type in catalog field Metadata")
		}
		q.Metadata = vt
	}

	// StructTreeRoot
	v, ok = dict["StructTreeRoot"]
	if ok {
		q.StructTreeRoot = v
	}

	// MarkInfo
	v, ok = dict["MarkInfo"]
	if ok {
		q.MarkInfo = v
	}

	// Lang
	v, ok = dict["Lang"]
	if ok {
		vt, ok := v.(String)
		if !ok {
			return errors.New("unexpected value type in catalog field Lang")
		}
		q.Lang = vt
	}

	// SpiderInfo
	v, ok = dict["SpiderInfo"]
	if ok {
		q.SpiderInfo = v
	}

	// OutputIntents
	v, ok = dict["OutputIntents"]
	if ok {
		q.OutputIntents = v
	}

	// return without error
	return nil
}

func (q DocumentCatalog) Copy(copyRef func(reference Reference) Reference) Object {
	return DocumentCatalog{
		Version:           q.Version.Copy(copyRef).(Name),
		Pages:             q.Pages.Copy(copyRef).(Reference),
		PageLabels:        Copy(q.PageLabels, copyRef),
		Names:             Copy(q.Names, copyRef),
		Dests:             q.Dests.Copy(copyRef).(Reference),
		ViewerPreferences: Copy(q.ViewerPreferences, copyRef),
		PageLayout:        q.PageLayout.Copy(copyRef).(Name),
		PageMode:          q.PageMode.Copy(copyRef).(Name),
		Outlines:          q.Outlines.Copy(copyRef).(Reference),
		Threads:           q.Threads.Copy(copyRef).(Array),
		OpenAction:        Copy(q.OpenAction, copyRef),
		AdditionalActions: Copy(q.AdditionalActions, copyRef),
		URI:               Copy(q.URI, copyRef).(Array),
		AcroForm:          Copy(q.AcroForm, copyRef),
		Metadata:          q.Metadata.Copy(copyRef).(Reference),
		StructTreeRoot:    Copy(q.StructTreeRoot, copyRef),
		MarkInfo:          Copy(q.MarkInfo, copyRef),
		Lang:              q.Lang.Copy(copyRef).(String),
		SpiderInfo:        Copy(q.SpiderInfo, copyRef),
		OutputIntents:     Copy(q.OutputIntents, copyRef),
	}
}

func (q DocumentCatalog) Equal(obj Object) bool {
	a, ok := obj.(DocumentCatalog)
	if !ok {
		return false
	}
	if !Equal(q.Version, a.Version) {
		return false
	}
	if !Equal(q.Pages, a.Pages) {
		return false
	}
	if !Equal(q.PageLabels, a.PageLabels) {
		return false
	}
	if !Equal(q.Names, a.Names) {
		return false
	}
	if !Equal(q.Dests, a.Dests) {
		return false
	}
	if !Equal(q.ViewerPreferences, a.ViewerPreferences) {
		return false
	}
	if !Equal(q.PageLayout, a.PageLayout) {
		return false
	}
	if !Equal(q.PageMode, a.PageMode) {
		return false
	}
	if !Equal(q.Outlines, a.Outlines) {
		return false
	}
	if !Equal(q.Threads, a.Threads) {
		return false
	}
	if !Equal(q.OpenAction, a.OpenAction) {
		return false
	}
	if !Equal(q.AdditionalActions, a.AdditionalActions) {
		return false
	}
	if !Equal(q.URI, a.URI) {
		return false
	}
	if !Equal(q.AcroForm, a.AcroForm) {
		return false
	}
	if !Equal(q.Metadata, a.Metadata) {
		return false
	}
	if !Equal(q.StructTreeRoot, a.StructTreeRoot) {
		return false
	}
	if !Equal(q.MarkInfo, a.MarkInfo) {
		return false
	}
	if !Equal(q.Lang, a.Lang) {
		return false
	}
	if !Equal(q.SpiderInfo, a.SpiderInfo) {
		return false
	}
	if !Equal(q.OutputIntents, a.OutputIntents) {
		return false
	}
	return true
}
