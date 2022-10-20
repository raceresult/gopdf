package extractor

import (
	"errors"

	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
)

// collectPages is a helper function to GetAllPAges
func collectPages(f *pdffile.File, node types.PageTreeNode) ([]types.Page, error) {
	var res []types.Page
	for _, kidRef := range node.Kids {
		kidObj, err := f.GetObject(kidRef)
		if err != nil {
			return nil, err
		}
		kidN, ok := kidObj.(types.PageTreeNode)
		if ok {
			pp, err := collectPages(f, kidN)
			if err != nil {
				return nil, err
			}
			res = append(res, pp...)
			continue
		}

		kidP, ok := kidObj.(types.Page)
		if ok {
			res = append(res, kidP)
			continue
		}

		return nil, errors.New("page tree kid is neither Page nor PageTreeNode")
	}
	return res, nil
}

// GetAllPages returns all pages from the PageTree
func (q *Extractor) GetAllPages() ([]types.Page, error) {
	catalogObj, err := q.file.GetObject(q.file.Root)
	if err != nil {
		return nil, err
	}
	catalog, ok := catalogObj.(types.DocumentCatalog)
	if !ok {
		return nil, errors.New("catalog invalid")
	}
	pagesObj, err := q.file.GetObject(catalog.Pages)
	if err != nil {
		return nil, err
	}
	pages, ok := pagesObj.(types.PageTreeNode)
	if !ok {
		return nil, errors.New("pages invalid")
	}

	return collectPages(q.file, pages)
}
