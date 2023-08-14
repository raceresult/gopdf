package pdf

import (
	"bytes"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"

	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
	_ "golang.org/x/image/bmp"
)

// File is the main object to create a pdf file
type File struct {
	ID    [2]types.String
	Info  types.InformationDictionary
	Pages []*Page

	// Threshold length for compressing content streams
	CompressStreamsThreshold int

	// PDF Version number
	Version float64

	// internals
	fonts         []FontHandler
	toUnicode     types.Reference
	cidSystemInfo types.Reference
	catalog       types.DocumentCatalog
	pageTree      types.PageTreeNode
	creator       *pdffile.File
	copiedObjects map[*pdffile.File]map[types.Reference]types.Reference
}

// NewFile creates a new File object
func NewFile() *File {
	q := &File{
		creator:                  pdffile.NewFile(),
		Version:                  pdffile.DefaultVersion,
		CompressStreamsThreshold: 500,
	}

	// catalog and page tree
	q.creator.Root = q.creator.AddObject(&q.catalog)
	q.catalog.Pages = q.creator.AddObject(&q.pageTree)

	return q
}

// NewPage adds and returns a new Page
func (q *File) NewPage(width, height float64) *Page {
	p := NewPage(width, height)
	q.Pages = append(q.Pages, p)
	return p
}

// CopyPage adds and returns a new Page
func (q *File) CopyPage(sourcePage types.Page, sourceFile *pdffile.File) *Page {
	if q.copiedObjects == nil {
		q.copiedObjects = make(map[*pdffile.File]map[types.Reference]types.Reference)
	}
	copiedMap, ok := q.copiedObjects[sourceFile]
	if !ok {
		copiedMap = make(map[types.Reference]types.Reference)
	}
	defer func() { q.copiedObjects[sourceFile] = copiedMap }()
	var copyRef func(ref types.Reference) types.Reference
	copyRef = func(ref types.Reference) types.Reference {
		if newRef, ok := copiedMap[ref]; ok {
			return newRef
		}

		obj, _ := sourceFile.GetObject(ref)
		copiedMap[ref] = types.Reference{} // to avoid endless recursion
		newRef := q.creator.AddObject(types.Copy(obj, copyRef))
		copiedMap[ref] = newRef
		return newRef
	}

	p := &Page{
		Data: sourcePage.Copy(copyRef).(types.Page),
	}
	q.Pages = append(q.Pages, p)
	return p
}

// WriteTo writes the parsed to the given writer
func (q *File) WriteTo(w io.Writer) (int64, error) {
	// finish fonts
	for _, f := range q.fonts {
		if err := f.finish(); err != nil {
			return 0, err
		}
	}

	// info
	if q.Info.Producer == "" {
		q.Info.Producer = "race result gopdf"
	}
	if q.Info.CreationDate.IsZero() {
		q.Info.CreationDate = types.Date(time.Now())
	}
	q.creator.Version = q.Version
	q.creator.ID = q.ID
	q.creator.Info = q.creator.AddObject(q.Info)

	// pages
	for _, page := range q.Pages {
		page.Data.Parent = q.catalog.Pages
		pageRef, err := page.create(q.creator, q.CompressStreamsThreshold)
		if err != nil {
			return 0, err
		}
		q.pageTree.Kids = append(q.pageTree.Kids, pageRef)
	}

	// output
	return q.creator.WriteTo(w)
}

// Write returns the PDF as byte slice
func (q *File) Write() ([]byte, error) {
	var bts bytes.Buffer
	_, err := q.WriteTo(&bts)
	return bts.Bytes(), err
}

// AddMetaData adds meta data to the document catalog
func (q *File) AddMetaData(data []byte) error {
	st, err := types.NewStream(data)
	if err != nil {
		return err
	}
	q.catalog.Metadata = q.creator.AddObject(st)
	return nil
}
