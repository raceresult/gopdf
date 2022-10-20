package types

// PDF Reference 1.4, 9.1 Predefined procedure sets

type ProcedureSet Name

const (
	ProcSetPDF    ProcedureSet = "PDF"
	ProcSetText   ProcedureSet = "Text"
	ProcSetImageB ProcedureSet = "ImageB"
	ProcSetImageC ProcedureSet = "ImageC"
	ProcSetImageI ProcedureSet = "ImageI"
)

func (q ProcedureSet) ToRawBytes() []byte {
	return Name(q).ToRawBytes()
}

func (q ProcedureSet) Copy(_ func(reference Reference) Reference) Object {
	return q
}
