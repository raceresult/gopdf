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
	commands [][]byte

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
	for _, ps := range pss {
		var found bool
		for _, ex := range q.Data.Resources.ProcSet {
			if ex == ps {
				found = true
				break
			}
		}
		if found {
			continue
		}
		q.Data.Resources.ProcSet = append(q.Data.Resources.ProcSet, ps)
	}
}

// AddFont adds a font to the list of resources of the page (unless already list ed) and return the font name
func (q *Page) AddFont(f FontHandler) types.Name {
	// create font Dictionary if not done yet
	if q.Data.Resources.Font == nil {
		q.Data.Resources.Font = types.Dictionary{}
	}

	// check if already listed
	ref := f.Reference()
	for k, v := range q.Data.Resources.Font {
		if v == ref {
			return k
		}
	}

	// create new name and add
	n := types.Name("F" + strconv.Itoa(len(q.Data.Resources.Font)+1))
	q.Data.Resources.Font[n] = ref
	return n
}

// AddXObject adds an XObject to the list of resources of the page (unless already listed) and return the resource name
func (q *Page) AddXObject(obj types.Reference) types.Name {
	// create XObject Dictionary if not done yet
	if q.Data.Resources.XObject == nil {
		q.Data.Resources.XObject = types.Dictionary{}
	}

	// check if already listed
	for k, v := range q.Data.Resources.XObject {
		if v == obj {
			return k
		}
	}

	// create new name and add
	n := types.Name("img" + strconv.Itoa(len(q.Data.Resources.XObject)+1))
	q.Data.Resources.XObject[n] = obj
	return n
}

// AddCommand adds any command/pdf operator to the content stream of the page
func (q *Page) AddCommand(operator string, args ...types.Object) {
	arr := make([][]byte, 0, len(args)+1)
	for _, v := range args {
		arr = append(arr, v.ToRawBytes())
	}
	arr = append(arr, []byte(operator))
	q.commands = append(q.commands, bytes.Join(arr, []byte{' '}))
}

// create is called when building the pdf file. It is supposed to add all objects to the creator and return a
// reference to the page object
func (q *Page) create(creator *pdffile.File, compress bool) (types.Reference, error) {
	// determine filter
	filter := types.Filter_NoFilter
	if compress {
		filter = types.Filter_FlateDecode
	}

	// create content stream
	stream, err := types.NewStream(bytes.Join(q.commands, []byte{'\n'}), filter)
	if err != nil {
		return types.Reference{}, err
	}
	streamRef := creator.AddObject(stream)

	// create page object and return reference to it
	q.Data.Contents = &streamRef
	return creator.AddObject(q.Data), nil
}
