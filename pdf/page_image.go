package pdf

import "git.rrdc.de/lib/gopdf/types"

// PDF Reference 1.4, Table 4.34 XObject operator

// XObject_Do paints the specified XObject. The operand name must appear as a key in the
// XObject subdictionary of the current resource dictionary (see Section 3.7.2,
// “Resource Dictionaries”); the associated value must be a stream whose Type
// entry, if present, is XObject. The effect of Do depends on the value of the
// XObject’s Subtype entry, which may be Image (see Section 4.8.4, “Image Dic-
// tionaries”), Form (Section 4.9, “Form XObjects”), or PS (Section 4.10, “Post-
// Script XObjects”).
func (q *Page) XObject_Do(img types.Reference) {
	q.AddProcSets(types.ProcSetImageB, types.ProcSetImageC, types.ProcSetImageI)

	n := q.AddXObject(img)
	q.AddCommand("Do", n)
}
