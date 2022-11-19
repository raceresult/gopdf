package types

import "errors"

// PDF Reference 1.4, Table  3.17 Required entries in a page tree node

type PageTreeNode struct {
	// (Required) The type of PDF object that this dictionary describes; must be Pages for
	// a page tree node.
	// Type Name

	// (Required except in root node; must be an indirect reference) The page tree node that
	// is the immediate parent of this one.
	Parent Reference

	// (Required) An array of indirect references to the immediate children of this node.
	// The children may be page objects or other page tree nodes.
	Kids []Reference

	// (Required) The number of leaf nodes (page objects) that are descendants of this
	// node within the page tree.
	Count Int
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
	if q.Parent.Number > 0 {
		d["Parent"] = q.Parent
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

	// Count
	v, ok = dict["Count"]
	if !ok {
		return errors.New("pages field Count missing")
	}
	q.Count, ok = v.(Int)
	if !ok {
		return errors.New("pages field Count invalid")
	}

	// return without error
	return nil
}

func (q PageTreeNode) Copy(copyRef func(reference Reference) Reference) Object {
	n := PageTreeNode{
		Parent: Copy(q.Parent, copyRef).(Reference),
		Count:  Copy(q.Count, copyRef).(Int),
	}
	for _, v := range q.Kids {
		n.Kids = append(n.Kids, Copy(v, copyRef).(Reference))
	}
	return n
}

func (q PageTreeNode) Equal(obj Object) bool {
	a, ok := obj.(PageTreeNode)
	if !ok {
		return false
	}
	if !Equal(q.Parent, a.Parent) {
		return false
	}
	if len(q.Kids) != len(a.Kids) {
		return false
	}
	for i := range q.Kids {
		if !Equal(q.Kids[i], a.Kids[i]) {
			return false
		}
	}
	if !Equal(q.Count, a.Count) {
		return false
	}
	return true
}
