package parser

import (
	"errors"

	"github.com/raceresult/gopdf/pdffile"
	"github.com/raceresult/gopdf/types"
)

// Parser provides functions to extract objects from a pdf file
type Parser struct {
	file *pdffile.File
}

// New creates a new Parser object
func New(bts []byte) (*Parser, error) {
	f, err := readFile(bts)
	if err != nil {
		return nil, err
	}
	return &Parser{
		file: f,
	}, nil
}

// File returns the parsed File object
func (q *Parser) File() *pdffile.File {
	return q.file
}

// Info returns the InformationDictionary of the parsed file
func (q *Parser) Info() (types.InformationDictionary, error) {
	obj, err := q.file.GetObject(q.file.Info)
	if err != nil {
		return types.InformationDictionary{}, err
	}

	dict, ok := obj.(types.Dictionary)
	if !ok {
		return types.InformationDictionary{}, errors.New("info is not a dictionary")
	}

	var ID types.InformationDictionary
	if err := ID.Read(dict, q.file); err != nil {
		return types.InformationDictionary{}, err
	}

	return ID, nil
}
