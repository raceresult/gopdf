package gopdf

import (
	"errors"
	"strconv"
	"sync"

	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
)

// Builder is the main object to build a PDF file
type Builder struct {
	Info types.InformationDictionary
	ID   [2]string

	// Threshold length for compressing content streams
	CompressStreamsThreshold int

	// PDF Version number
	Version float64

	// number of worker routines used to generate content streams of pages
	WorkerRoutines int

	// internals
	file     *pdf.File
	pages    []*Page
	currPage *Page
}

// New creates a new Builder object
func New() *Builder {
	return &Builder{
		file:                     pdf.NewFile(),
		CompressStreamsThreshold: 500,
		Version:                  pdffile.DefaultVersion,
	}
}

// AddElement adds one or more elements to the current page
func (q *Builder) AddElement(item ...Element) {
	if q.currPage == nil {
		q.NewPage(GetStandardPageSize(PageSizeA4, false))
	}
	q.currPage.AddElement(item...)
}

// Build builds the PDF document and returns the file as byte slice
func (q *Builder) Build() ([]byte, error) {
	// settings
	q.file.Version = q.Version
	q.file.Info = q.Info
	q.file.ID = [2]types.String{types.String(q.ID[0]), types.String(q.ID[1])}
	q.file.CompressStreamsThreshold = q.CompressStreamsThreshold

	// create pages
	pdfPages := make([]*pdf.Page, 0, len(q.pages))
	for _, p := range q.pages {
		pdfPages = append(pdfPages, q.file.NewPage(p.Width.Pt(), p.Height.Pt()))
	}

	// determine number of workers
	workers := q.WorkerRoutines
	if workers < 1 {
		workers = 1
	}
	if workers > len(q.pages) {
		workers = len(q.pages)
	}

	// start workers
	var wg sync.WaitGroup
	wg.Add(workers)
	errs := make([]error, workers)
	for i := 0; i < workers; i++ {
		go func(z int) {
			defer wg.Done()

			for k := z; k < len(q.pages); k += workers {
				if err := q.pages[k].build(pdfPages[k]); err != nil {
					errs[z] = err
					return
				}
			}
		}(i)
	}

	// wait and check for error
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	// create byte stream
	return q.file.Write()
}

// NewPage adds a new page to the pdf
func (q *Builder) NewPage(size PageSize) *Page {
	q.currPage = NewPage(size)
	q.pages = append(q.pages, q.currPage)
	return q.currPage
}

// NewPageBefore inserts a new page before the given pageNo to the pdf
func (q *Builder) NewPageBefore(size PageSize, beforePageNo int) *Page {
	if beforePageNo > len(q.pages) {
		return q.NewPage(size)
	}
	if beforePageNo < 1 {
		beforePageNo = 1
	}

	q.currPage = NewPage(size)
	q.pages = append(q.pages, nil)
	copy(q.pages[beforePageNo:], q.pages[beforePageNo-1:])
	q.pages[beforePageNo-1] = q.currPage

	return q.currPage
}

// NewFormFromPage creates a new form object from the give page
func (q *Builder) NewFormFromPage(page *Page) (*Form, error) {
	p := pdf.NewPage(page.Width.Pt(), page.Height.Pt())
	for _, item := range page.elements {
		if err := item.Build(p); err != nil {
			return nil, err
		}
	}

	ref, err := q.file.NewFormFromPage(p)
	if err != nil {
		return nil, err
	}

	return &Form{
		BBox: types.Rectangle{URX: types.Number(page.Width.Pt()), URY: types.Number(page.Height.Pt())},
		Form: ref,
	}, nil
}

// PageCount returns the number of pages already added
func (q *Builder) PageCount() int {
	return len(q.pages)
}

// SelectPage set the current page to the page with the given number (starting from 1)
func (q *Builder) SelectPage(pageNo int) error {
	if pageNo < len(q.pages) || pageNo > len(q.pages) {
		return errors.New("page " + strconv.Itoa(pageNo) + " not found")
	}
	q.currPage = q.pages[pageNo-1]
	return nil
}

// CurrPage returns the current page
func (q *Builder) CurrPage() *Page {
	return q.currPage
}

// NewImage adds a new image to the PDF file
func (q *Builder) NewImage(bts []byte) (*pdf.Image, error) {
	return q.file.NewImage(bts)
}

// NewCapturedPage adds a new captured page to the PDF file
func (q *Builder) NewCapturedPage(sourcePage types.Page, sourceFile *pdffile.File) (*Form, error) {
	cp, err := q.file.NewCapturedPage(sourcePage, sourceFile)
	if err != nil {
		return nil, err
	}
	return &Form{
		BBox: sourcePage.MediaBox,
		Form: cp,
	}, nil
}

// NewStandardFont adds a new standard font (expected to be available in all PDF consuming systems) to the pdf
func (q *Builder) NewStandardFont(name types.StandardFontName, encoding types.Encoding) (*pdf.StandardFont, error) {
	return q.file.NewStandardFont(name, encoding)
}

// NewTrueTypeFont adds a new TrueType font to the pdf
func (q *Builder) NewTrueTypeFont(ttf []byte, encoding types.Encoding, embed bool) (*pdf.TrueTypeFont, error) {
	return q.file.NewTrueTypeFont(ttf, encoding, embed)
}

// NewCompositeFont adds a font as composite font to the pdf, i.e. with Unicode support
func (q *Builder) NewCompositeFont(ttf []byte) (*pdf.CompositeFont, error) {
	return q.file.NewCompositeFontFromTTF(ttf)
}

// NewCompositeFontFromOTF adds an otf font as composite font to the pdf, i.e. with Unicode support
func (q *Builder) NewCompositeFontFromOTF(otf []byte) (*pdf.CompositeFontOTF, error) {
	return q.file.NewCompositeFontFromOTF(otf)
}
