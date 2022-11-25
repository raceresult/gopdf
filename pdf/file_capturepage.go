package pdf

import (
	"bytes"
	"errors"

	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
)

// NewCapturedPage is used to copy a page from another pdf
func (q *File) NewCapturedPage(sourcePage types.Page, sourceFile *pdffile.File) (types.Reference, error) {
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

	// collect content references
	var cc []types.Reference
	if sourcePage.Contents == nil {
	} else if ref, ok := sourcePage.Contents.(types.Reference); ok {
		cc = append(cc, ref)
	} else if arr, ok := sourcePage.Contents.(types.Array); ok {
		for _, v := range arr {
			ref, ok := v.(types.Reference)
			if !ok {
				return types.Reference{}, errors.New("content stream has unexpected type")
			}
			cc = append(cc, ref)
		}
	} else {
		return types.Reference{}, errors.New("content stream has unexpected type")
	}

	// copy content objects
	var data []byte
	for _, ref := range cc {
		cs, err := sourceFile.GetObject(ref)
		if err != nil {
			return types.Reference{}, err
		}
		if s, ok := cs.(types.StreamObject); ok {
			decoded, err := s.Decode()
			if err != nil {
				return types.Reference{}, err
			}
			data = append(data, decoded...)
			data = append(data, '\n')
		}
	}

	// create new content stream
	stream, err := types.NewStream(data, types.Filter_FlateDecode)
	if err != nil {
		return types.Reference{}, err
	}

	// add form object
	return q.creator.AddObject(types.Form{
		Stream:     stream.Stream,
		Dictionary: stream.Dictionary.(types.StreamDictionary),
		BBox:       sourcePage.MediaBox,
		Resources:  types.Copy(sourcePage.Resources, copyRef),
	}), nil
}

// NewFormFromPage creates a new form object from the given page
func (q *File) NewFormFromPage(page *Page) (types.Reference, error) {
	stream, err := types.NewStream(bytes.Join(page.contents, []byte{'\n'}), types.Filter_FlateDecode)
	if err != nil {
		return types.Reference{}, err
	}

	return q.creator.AddObject(types.Form{
		Stream:     stream.Stream,
		Dictionary: stream.Dictionary.(types.StreamDictionary),
		BBox:       page.Data.MediaBox,
		Resources:  page.Data.Resources,
	}), nil
}
