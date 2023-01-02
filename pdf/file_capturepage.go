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
		copiedMap[ref] = types.Reference{} // to avoid endless recursion
		newRef := q.creator.AddObject(types.Copy(obj, copyRef))
		copiedMap[ref] = newRef
		return newRef
	}

	// collect content
	var data []byte
	var addContent func(obj types.Object) error
	addContent = func(obj types.Object) error {
		if obj == nil {
			return nil
		}
		switch item := obj.(type) {
		case types.Reference:
			newItem, err := sourceFile.GetObject(item)
			if err != nil {
				return err
			}
			addContent(newItem)
		case types.Array:
			for _, v := range item {
				addContent(v)
			}
		case types.StreamObject:
			decoded, err := item.Decode(sourceFile)
			if err != nil {
				return err
			}
			data = append(data, decoded...)
			data = append(data, '\n')
		default:
			return errors.New("content stream has unexpected type")
		}
		return nil
	}
	addContent(sourcePage.Contents)

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
