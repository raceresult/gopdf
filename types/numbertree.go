package types

// PDF Reference 1.4, Table 3.25 Entries in a number tree node dictionary

type NumberTree struct {
	// (Root and intermediate nodes only; required in intermediate nodes; present in the root node
	// if and only if Names is not present) An array of indirect references to the immediate chil-
	// dren of this node. The children may be intermediate or leaf nodes.
	Kids Array

	// (Root and leaf nodes only; required in leaf nodes; present in the root node if and only if Kids
	// is not present) An array of the form
	// [key1 value1 key2 value2 ... keyn valuen]
	// where each keyi is a string and the corresponding valuei is an indirect reference to the
	// object associated with that key. The keys are sorted in lexical order, as described below.
	Names Array

	// (Intermediate and leaf nodes only; required) An array of two strings, specifying the (lexi-
	// cally) least and greatest keys included in the Names array of a leaf node or in the Names
	// arrays of any leaf nodes that are descendants of an intermediate node.
	Limits Array
}

func (q NumberTree) ToRawBytes() []byte {
	d := Dictionary{}
	if len(q.Kids) != 0 {
		d["Kids"] = q.Kids
	}
	if len(q.Names) != 0 {
		d["Names"] = q.Names
	}
	if len(q.Limits) != 0 {
		d["Limits"] = q.Limits
	}
	return d.ToRawBytes()
}

func (q NumberTree) Copy(copyRef func(reference Reference) Reference) Object {
	return NumberTree{
		Kids:   q.Kids.Copy(copyRef).(Array),
		Names:  q.Names.Copy(copyRef).(Array),
		Limits: q.Limits.Copy(copyRef).(Array),
	}
}

func (q NumberTree) Equal(obj Object) bool {
	a, ok := obj.(NumberTree)
	if !ok {
		return false
	}
	if !Equal(q.Kids, a.Kids) {
		return false
	}
	if !Equal(q.Names, a.Names) {
		return false
	}
	if !Equal(q.Limits, a.Limits) {
		return false
	}
	return true
}
