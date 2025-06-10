package types

type FileSpec struct {
	// Type // Required

	AFRelationship Name
	Desc           String
	UF             String
	EF             Dictionary
	F              String
}

func (q FileSpec) ToRawBytes() []byte {
	d := Dictionary{
		"Type": Name("Filespec"),
	}

	if q.AFRelationship != "" {
		d["AFRelationship"] = q.AFRelationship
	}
	if q.Desc != "" {
		d["Desc"] = q.Desc
	}
	if q.UF != "" {
		d["UF"] = q.UF
	}
	if len(q.EF) != 0 {
		d["EF"] = q.EF
	}
	if q.F != "" {
		d["F"] = q.F
	}

	return d.ToRawBytes()
}

func (q FileSpec) Copy(copyRef func(reference Reference) Reference) Object {
	return FileSpec{
		AFRelationship: q.AFRelationship.Copy(copyRef).(Name),
		Desc:           q.Desc.Copy(copyRef).(String),
		UF:             q.UF.Copy(copyRef).(String),
		EF:             q.EF.Copy(copyRef).(Dictionary),
		F:              q.F.Copy(copyRef).(String),
	}
}

func (q FileSpec) Equal(obj Object) bool {
	a, ok := obj.(FileSpec)
	if !ok {
		return false
	}
	if !Equal(q.AFRelationship, a.AFRelationship) {
		return false
	}
	if !Equal(q.Desc, a.Desc) {
		return false
	}
	if !Equal(q.UF, a.UF) {
		return false
	}
	if !Equal(q.EF, a.EF) {
		return false
	}
	if !Equal(q.F, a.F) {
		return false
	}
	return true
}
