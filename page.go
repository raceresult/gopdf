package gopdf

import (
	"github.com/raceresult/gopdf/pdf"
	"github.com/raceresult/gopdf/types"
)

// Page represents on page in the PDF document to which elements can be added arbitrarily
type Page struct {
	Width  Length
	Height Length
	Rotate int

	elements []Element
}

// NewPage creates a new Page object with the given page size
func NewPage(size PageSize) *Page {
	return &Page{
		Width:  size[0],
		Height: size[1],
	}
}

// AddElement adds one or more elements to the page
func (q *Page) AddElement(item ...Element) {
	q.elements = append(q.elements, item...)
}

// build is called when the PDF file is created and calls the Build function on all elements
func (q *Page) build(page *pdf.Page) ([]string, error) {
	page.Data.Rotate = types.Int(q.Rotate)

	var warnings []string
	for _, item := range q.elements {
		w, err := item.Build(page)
		if err != nil {
			return nil, err
		}
		if w != "" {
			warnings = append(warnings, w)
		}
	}
	return warnings, nil
}
