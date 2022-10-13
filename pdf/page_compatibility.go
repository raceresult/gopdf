package pdf

// PDF Reference 1.4, Table 3.20 Compatibility operators

// Compatibility_BX begins a compatibility section. Unrecognized operators (along with their
// operands) will be ignored without error until the balancing EX operator is encoun-
// tered.
func (q *Page) Compatibility_BX() {
	q.AddCommand("BX")
}

// Compatibility_EX ends a compatibility section begun by a balancing BX operator.
func (q *Page) Compatibility_EX() {
	q.AddCommand("EX")
}
