package types

import "errors"

// PDF Reference 1.4, Table 3.23 Entries in a name tree node dictionary

type PageTreeNode struct {
	// Type
	// count  int

	Kids []Reference
}

func (q PageTreeNode) ToRawBytes() []byte {
	kids := make(Array, 0, len(q.Kids))
	for _, k := range q.Kids {
		kids = append(kids, k)
	}
	d := Dictionary{
		"Type":  Name("Pages"),
		"Count": Int(len(q.Kids)),
		"Kids":  kids,
	}

	return d.ToRawBytes()
}

func (q *PageTreeNode) Read(dict Dictionary) error {
	// Type
	v, ok := dict["Type"]
	if !ok {
		return errors.New("pages missing Type")
	}
	dtype, ok := v.(Name)
	if !ok {
		return errors.New("pages field Type invalid")
	}
	if dtype != "Pages" {
		return errors.New("unexpected value in pages field Type")
	}

	// Kids
	v, ok = dict["Kids"]
	if !ok {
		return errors.New("pages field Kids missing")
	}
	pages, ok := v.(Array)
	if !ok {
		return errors.New("pages field Kids invalid")
	}
	for _, a := range pages {
		va, ok := a.(Reference)
		if !ok {
			return errors.New("pages field Kids invalid")
		}
		q.Kids = append(q.Kids, va)
	}

	// return without error
	return nil
}
