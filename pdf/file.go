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
	ID              [2]types.String
	Info            types.InformationDictionary
	Pages           []*Page
	CompressStreams bool
	Version         float64

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
		creator:         pdffile.NewFile(),
		Version:         pdffile.DefaultVersion,
		CompressStreams: true,
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
		pageRef, err := page.create(q.creator, q.CompressStreams)
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
