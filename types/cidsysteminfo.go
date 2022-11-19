package types

// PDF Reference 1.4, Table 5.12 Entries in a CIDSystemInfo dictionary

type CIDSystemInfoDictionary struct {
	// (Required) A string identifying the issuer of the character collection—for exam-
	// ple, Adobe. For information about assigning a registry identifier, consult the ASN
	// Developer Program Web site or contact the Adobe Solutions Network (see the
	// Bibliography).
	Registry String

	// (Required) A string that uniquely names the character collection within the speci-
	// fied registry—for example, Japan1.
	Ordering String

	// (Required) The supplement number of the character collection. An original charac-
	// ter collection has a supplement number of 0. Whenever additional CIDs are
	// assigned in a character collection, the supplement number is increased. Supple-
	// ments do not alter the ordering of existing CIDs in the character collection. This
	// value is not used in determining compatibility between character collections.
	Supplement Int
}

func (q CIDSystemInfoDictionary) ToRawBytes() []byte {
	d := Dictionary{
		"Registry":   q.Registry,
		"Ordering":   q.Ordering,
		"Supplement": q.Supplement,
	}

	return d.ToRawBytes()
}

func (q CIDSystemInfoDictionary) Copy(copyRef func(reference Reference) Reference) Object {
	return CIDSystemInfoDictionary{
		Registry:   q.Registry.Copy(copyRef).(String),
		Ordering:   q.Ordering.Copy(copyRef).(String),
		Supplement: q.Supplement.Copy(copyRef).(Int),
	}
}

func (q CIDSystemInfoDictionary) Equal(obj Object) bool {
	a, ok := obj.(CIDSystemInfoDictionary)
	if !ok {
		return false
	}
	if !Equal(q.Registry, a.Registry) {
		return false
	}
	if !Equal(q.Ordering, a.Ordering) {
		return false
	}
	if !Equal(q.Supplement, a.Supplement) {
		return false
	}
	return true
}
