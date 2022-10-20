package pdf

import (
	"bytes"
	"strconv"

	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
)

// Page holds all information about a pdf page and the commands that will be put in the content stream
type Page struct {
	// additional page data, can also be modified from the outside
	Data types.Page

	// list of commands/operators already added to the page
	contents []interface{}

	// we need a reference to the current font for character encoding
	currFont FontHandler

	// internal text and graphics state to check if commands actually change the state
	textState          *types.TextState
	graphicsState      *graphicsState
	graphicsStateStack []*graphicsState
}

// newPage creates and returns a new page
func newPage(width, height float64, parent types.Reference) *Page {
	return &Page{
		Data: types.Page{
			MediaBox: types.Rectangle{LLX: 0, LLY: 0, URX: types.Number(width), URY: types.Number(height)},
			Parent:   parent,
		},
		textState:     types.NewTextState(),
		graphicsState: &graphicsState{},
	}
}

// AddProcSets adds ProcedureSets to the page resources unless already listed
func (q *Page) AddProcSets(pss ...types.ProcedureSet) {
	var d types.Array
	if q.Data.Resources.ProcSet != nil {
		d = q.Data.Resources.ProcSet.(types.Array)
	}

	for _, ps := range pss {
		var found bool
		for _, ex := range d {
			if ex == ps {
				found = true
				break
			}
		}
		if found {
			continue
		}
		d = append(d, ps)
	}
	q.Data.Resources.ProcSet = d
}

// AddFont adds a font to the list of resources of the page (unless already list ed) and return the font name
func (q *Page) AddFont(f FontHandler) types.Name {
	// create font Dictionary if not done yet
	var d types.Dictionary
	if q.Data.Resources.Font != nil {
		d = q.Data.Resources.Font.(types.Dictionary)
	} else {
		d = make(types.Dictionary)
	}

	// check if already listed
	ref := f.Reference()
	for k, v := range d {
		if v == ref {
			return k
		}
	}

	// create new name and add
	n := types.Name("F" + strconv.Itoa(len(d)+1))
	d[n] = ref
	q.Data.Resources.Font = d
	return n
}

// AddXObject adds an XObject to the list of resources of the page (unless already listed) and return the resource name
func (q *Page) AddXObject(obj types.Reference) types.Name {
	// create XObject Dictionary if not done yet
	var d types.Dictionary
	if q.Data.Resources.XObject != nil {
		d = q.Data.Resources.XObject.(types.Dictionary)
	} else {
		d = make(types.Dictionary)
	}

	// check if already listed
	for k, v := range d {
		if v == obj {
			return k
		}
	}

	// create new name and add
	n := types.Name("img" + strconv.Itoa(len(d)+1))
	d[n] = obj
	q.Data.Resources.XObject = d
	return n
}

// AddCommand adds any command/pdf operator to the content stream of the page
func (q *Page) AddCommand(operator string, args ...types.Object) {
	arr := make([][]byte, 0, len(args)+1)
	for _, v := range args {
		arr = append(arr, v.ToRawBytes())
	}
	arr = append(arr, []byte(operator))
	q.contents = append(q.contents, bytes.Join(arr, []byte{' '}))
}

// AddCapturedPage adds the content of a captured page to this page
func (q *Page) AddCapturedPage(cp *CapturedPage) {
	// ExtGState
	if cp.Resources.ExtGState != nil {
		d, ok := q.Data.Resources.ExtGState.(types.Dictionary)
		if !ok {
			d = make(types.Dictionary)
		}
		for k, v := range cp.Resources.ExtGState {
			d[k] = v
		}
		q.Data.Resources.ExtGState = d
	}

	// ColorSpace
	if cp.Resources.ColorSpace != nil {
		d, ok := q.Data.Resources.ColorSpace.(types.Dictionary)
		if !ok {
			d = make(types.Dictionary)
		}
		for k, v := range cp.Resources.ColorSpace {
			d[k] = v
		}
		q.Data.Resources.ColorSpace = d
	}

	// Pattern
	if cp.Resources.Pattern != nil {
		d, ok := q.Data.Resources.Pattern.(types.Dictionary)
		if !ok {
			d = make(types.Dictionary)
		}
		for k, v := range cp.Resources.Pattern {
			d[k] = v
		}
		q.Data.Resources.Pattern = d
	}

	// Shading
	if cp.Resources.Shading != nil {
		d, ok := q.Data.Resources.Shading.(types.Dictionary)
		if !ok {
			d = make(types.Dictionary)
		}
		for k, v := range cp.Resources.Shading {
			d[k] = v
		}
		q.Data.Resources.Shading = d
	}

	// xObjects
	if cp.Resources.XObject != nil {
		d, ok := q.Data.Resources.XObject.(types.Dictionary)
		if !ok {
			d = make(types.Dictionary)
		}
		for k, v := range cp.Resources.XObject {
			d[k] = v
		}
		q.Data.Resources.XObject = d
	}

	// fonts
	if cp.Resources.Font != nil {
		d, ok := q.Data.Resources.Font.(types.Dictionary)
		if !ok {
			d = make(types.Dictionary)
		}
		for k, v := range cp.Resources.Font {
			d[k] = v
		}
		q.Data.Resources.Font = d
	}

	// ProcSet
	if cp.Resources.ProcSet != nil {
		d, ok := q.Data.Resources.ProcSet.(types.Array)
		if !ok {
			d = make(types.Array, 0)
		}
		for _, v := range cp.Resources.ProcSet {
			d = append(d, v)
		}
		q.Data.Resources.ProcSet = d
	}

	// Properties
	if cp.Resources.Properties != nil {
		d, ok := q.Data.Resources.Properties.(types.Dictionary)
		if !ok {
			d = make(types.Dictionary)
		}
		for k, v := range cp.Resources.Properties {
			d[k] = v
		}
		q.Data.Resources.Properties = d
	}

	// stream
	for _, c := range cp.Contents {
		q.contents = append(q.contents, c)
	}
}

// create is called when building the pdf file. It is supposed to add all objects to the creator and return a
// reference to the page object
func (q *Page) create(creator *pdffile.File, compress bool) (types.Reference, error) {
	// determine filter
	filter := types.Filter_NoFilter
	if compress {
		filter = types.Filter_FlateDecode
	}

	// Contents can be a mixture of byte slices and references to existing content streams (when adding page from other document).
	// We need to join existing byte parts
	var allRefs types.Array
	var currBytes [][]byte
	flushBytes := func() error {
		if len(currBytes) == 0 {
			return nil
		}
		stream, err := types.NewStream(bytes.Join(currBytes, []byte{'\n'}), filter)
		if err != nil {
			return err
		}
		allRefs = append(allRefs, creator.AddObject(stream))
		currBytes = nil
		return nil
	}
	for _, c := range q.contents {
		switch v := c.(type) {
		case []byte:
			currBytes = append(currBytes, v)
		case types.Reference:
			if err := flushBytes(); err != nil {
				return types.Reference{}, err
			}
			allRefs = append(allRefs, v)
		}
	}
	if err := flushBytes(); err != nil {
		return types.Reference{}, err
	}

	// create page object and return reference to it
	if len(allRefs) == 1 {
		q.Data.Contents = allRefs[0]
	} else {
		q.Data.Contents = allRefs
	}
	return creator.AddObject(q.Data), nil
}
