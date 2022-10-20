package pdf

import (
	"errors"

	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
)

type CapturedPage struct {
	Contents  []types.Reference
	Resources struct {
		ExtGState  types.Dictionary
		ColorSpace types.Dictionary
		Pattern    types.Dictionary
		Shading    types.Dictionary
		XObject    types.Dictionary
		Font       types.Dictionary
		ProcSet    types.Array
		Properties types.Dictionary
	}
}

func (cp *CapturedPage) Build(page *Page) {
	page.AddCapturedPage(cp)
}

// NewCapturedPage is used to copy a page from another pdf
func (q *File) NewCapturedPage(sourcePage types.Page, sourceFile *pdffile.File) (*CapturedPage, error) {
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
		newRef := q.creator.AddObject(types.Copy(obj, copyRef))
		copiedMap[ref] = newRef
		return newRef
	}

	// ExtGState
	var res CapturedPage
	res.Resources.ExtGState, _ = types.Copy(sourcePage.Resources.ExtGState, copyRef).(types.Dictionary)
	res.Resources.ColorSpace, _ = types.Copy(sourcePage.Resources.ColorSpace, copyRef).(types.Dictionary)
	res.Resources.Pattern, _ = types.Copy(sourcePage.Resources.Pattern, copyRef).(types.Dictionary)
	res.Resources.Shading, _ = types.Copy(sourcePage.Resources.Shading, copyRef).(types.Dictionary)
	res.Resources.XObject, _ = types.Copy(sourcePage.Resources.XObject, copyRef).(types.Dictionary)
	res.Resources.Font, _ = types.Copy(sourcePage.Resources.Font, copyRef).(types.Dictionary)
	res.Resources.ProcSet, _ = types.Copy(sourcePage.Resources.ProcSet, copyRef).(types.Array)
	res.Resources.Properties, _ = types.Copy(sourcePage.Resources.Properties, copyRef).(types.Dictionary)

	// stream
	var cc []types.Reference
	if sourcePage.Contents == nil {
	} else if ref, ok := sourcePage.Contents.(types.Reference); ok {
		cc = append(cc, ref)
	} else if arr, ok := sourcePage.Contents.(types.Array); ok {
		for _, v := range arr {
			ref, ok := v.(types.Reference)
			if !ok {
				return nil, errors.New("content stream has unexpected type")
			}
			cc = append(cc, ref)
		}
	} else {
		return nil, errors.New("content stream has unexpected type")
	}
	for _, ref := range cc {
		cs, err := sourceFile.GetObject(ref)
		if err != nil {
			return nil, err
		}
		res.Contents = append(res.Contents, q.creator.AddObject(cs))
	}

	// return without error
	return &res, nil
}
