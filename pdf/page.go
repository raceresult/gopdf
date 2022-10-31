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
	contents [][]byte

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

// create is called when building the pdf file. It is supposed to add all objects to the creator and return a
// reference to the page object
func (q *Page) create(creator *pdffile.File, compress bool) (types.Reference, error) {
	// determine filter
	filter := types.Filter_NoFilter
	if compress {
		filter = types.Filter_FlateDecode
	}

	// create stream
	stream, err := types.NewStream(bytes.Join(q.contents, []byte{'\n'}), filter)
	if err != nil {
		return types.Reference{}, err
	}

	// create page object and return reference to it
	q.Data.Contents = creator.AddObject(stream)
	return creator.AddObject(q.Data), nil
}
