package parser

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
func (q *Parser) GetAllPages() ([]types.Page, error) {
	// get page tree root
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

	// collect all pages recursively
	return collectPages(q.file, pages)
}

// GetPage returns one page from the PageTree (first page = pageNo 1)
func (q *Parser) GetPage(pageNo int) (types.Page, error) {
	// check parameter
	if pageNo < 1 {
		return types.Page{}, errors.New("invalid page no")
	}

	// get page tree root
	catalogObj, err := q.file.GetObject(q.file.Root)
	if err != nil {
		return types.Page{}, err
	}
	catalog, ok := catalogObj.(types.DocumentCatalog)
	if !ok {
		return types.Page{}, errors.New("catalog invalid")
	}
	pagesObj, err := q.file.GetObject(catalog.Pages)
	if err != nil {
		return types.Page{}, err
	}
	pages, ok := pagesObj.(types.PageTreeNode)
	if !ok {
		return types.Page{}, errors.New("pages invalid")
	}

	// fast path: check if root kids are pages
	if len(pages.Kids) >= pageNo {
		obj, err := q.file.GetObject(pages.Kids[pageNo-1])
		if err != nil {
			return types.Page{}, err
		}
		if p, ok := obj.(types.Page); ok {
			return p, nil
		}
	}

	// slow path: collect all pages in tree, then return page
	allPages, err := collectPages(q.file, pages)
	if err != nil {
		return types.Page{}, err
	}
	if pageNo > len(allPages) {
		return types.Page{}, errors.New("invalid page no")
	}
	return allPages[pageNo-1], nil
}