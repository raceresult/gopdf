package pdf

// PDF Reference 1.4, Table 4.38 Inline image operators

// InlineImage_BI begins an inline image object
func (q *Page) InlineImage_BI() {
	q.AddCommand("BI")
}

// InlineImage_ID begins the image data for an inline image object.
func (q *Page) InlineImage_ID() {
	q.AddCommand("ID")
}

// InlineImage_EI ends an inline image object
func (q *Page) InlineImage_EI() {
	q.AddCommand("EI")
}
