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
	ID              [2]string
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
}

// NewFile creates a new File object
func NewFile() *File {
	q := &File{
		creator: pdffile.NewFile(),
		Info: types.InformationDictionary{
			Creator:      "race result gopdf",
			CreationDate: types.Date(time.Now()),
		},
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
	p := newPage(width, height, q.catalog.Pages)
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
	q.creator.Version = q.Version
	q.creator.ID = q.ID
	q.creator.Info = q.creator.AddObject(q.Info)

	// pages
	for _, page := range q.Pages {
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
