package parser

import (
	"errors"

	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
)

// collectPages is a helper function to GetAllPAges
func collectPages(f *pdffile.File, node *types.PageTreeNode) ([]types.Page, error) {
	var res []types.Page
	for _, kidRef := range node.Kids {
		kidObj, err := f.GetObject(kidRef)
		if err != nil {
			return nil, err
		}

		kidDict, ok := kidObj.(types.Dictionary)
		if !ok {
			return nil, errors.New("kid object is not a dictionary")
		}
		typ, ok := kidDict["Type"]
		if !ok {
			return nil, errors.New("page tree item does not have Type")
		}
		typName, ok := typ.(types.Name)
		if !ok {
			return nil, errors.New("page tree item Type is not Name")
		}

		switch typName {
		case "Pages":
			var ptn types.PageTreeNode
			if err := ptn.Read(kidDict); err != nil {
				return nil, err
			}
			pp, err := collectPages(f, &ptn)
			if err != nil {
				return nil, err
			}
			res = append(res, pp...)

		case "Page":
			var p types.Page
			if err := p.Read(kidDict); err != nil {
				return nil, err
			}
			res = append(res, p)

		default:
			return nil, errors.New("unknown page tree item type " + string(typName))
		}
	}
	return res, nil
}

// GetAllPages returns all pages from the PageTree
func (q *Parser) GetAllPages() ([]types.Page, error) {
	// get page tree root
	ptn, err := q.getPageTreeRoot()
	if err != nil {
		return nil, err
	}

	// collect all pages recursively
	return collectPages(q.file, ptn)
}

// GetPage returns one page from the PageTree (first page = pageNo 1)
func (q *Parser) GetPage(pageNo int) (types.Page, error) {
	// check parameter
	if pageNo < 1 {
		return types.Page{}, errors.New("invalid page no")
	}

	// get page tree root
	ptn, err := q.getPageTreeRoot()
	if err != nil {
		return types.Page{}, err
	}

	// slow path: collect all pages in tree, then return page
	allPages, err := collectPages(q.file, ptn)
	if err != nil {
		return types.Page{}, err
	}
	if pageNo > len(allPages) {
		return types.Page{}, errors.New("invalid page no")
	}
	return allPages[pageNo-1], nil
}

func (q *Parser) getCatalog() (*types.DocumentCatalog, error) {
	catalogObj, err := q.file.GetObject(q.file.Root)
	if err != nil {
		return nil, err
	}
	catalogDict, ok := catalogObj.(types.Dictionary)
	if !ok {
		return nil, errors.New("catalog invalid")
	}
	var cat types.DocumentCatalog
	if err := cat.Read(catalogDict); err != nil {
		return nil, err
	}
	return &cat, nil
}

func (q *Parser) getPageTreeRoot() (*types.PageTreeNode, error) {
	cat, err := q.getCatalog()
	if err != nil {
		return nil, err
	}
	pagesObj, err := q.file.GetObject(cat.Pages)
	if err != nil {
		return nil, err
	}
	pagesDict, ok := pagesObj.(types.Dictionary)
	if !ok {
		return nil, errors.New("catalog pages is not a dictionary")
	}
	var ptn types.PageTreeNode
	if err := ptn.Read(pagesDict); err != nil {
		return nil, err
	}
	return &ptn, nil
}
