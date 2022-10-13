package types

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
