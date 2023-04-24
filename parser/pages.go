package parser

import (
	"errors"

	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
)

// collectPages is a helper function to GetAllPages
func collectPages(f *pdffile.File, node types.Dictionary) ([]types.Page, error) {
	// get node type
	typ, ok := node["Type"]
	if !ok {
		return nil, errors.New("page tree item does not have Type")
	}
	typName, ok := typ.(types.Name)
	if !ok {
		return nil, errors.New("page tree item Type is not Name")
	}

	// proceed depending on node type
	var res []types.Page
	switch typName {
	case "Pages":
		var ptn types.PageTreeNode
		if err := ptn.Read(node); err != nil {
			return nil, err
		}

		// iterate over kid nodes
		for _, kidRef := range ptn.Kids {
			// get kid dictionaries
			kidObj, err := f.GetObject(kidRef)
			if err != nil {
				return nil, err
			}
			kidDict, ok := kidObj.(types.Dictionary)
			if !ok {
				return nil, errors.New("kid object is not a dictionary")
			}
			for k, v := range node {
				if _, ok := kidDict[k]; !ok {
					kidDict[k] = v
				}
			}

			// collect pages from all kids
			pp, err := collectPages(f, kidDict)
			if err != nil {
				return nil, err
			}
			res = append(res, pp...)
		}

	case "Page":
		var p types.Page
		if err := p.Read(node, f); err != nil {
			return nil, err
		}
		res = append(res, p)

	default:
		return nil, errors.New("unknown page tree item type " + string(typName))
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

	// collect all pages in tree, then return page
	allPages, err := collectPages(q.file, ptn)
	if err != nil {
		return types.Page{}, err
	}
	if pageNo > len(allPages) {
		return types.Page{}, errors.New("invalid page no")
	}
	return allPages[pageNo-1], nil
}

// getCatalog returns the document catalog of the parsed pdf
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

// getPageTreeRoot returns the page tree root of the parsed pdf
func (q *Parser) getPageTreeRoot() (types.Dictionary, error) {
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
	return pagesDict, nil
}
