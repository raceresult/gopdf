package gopdf

import "github.com/raceresult/gopdf/pdf"

// Element is the interface any object needs to fulfill to be added to the content stream of a page
type Element interface {
	Build(page *pdf.Page) (string, error)
}
