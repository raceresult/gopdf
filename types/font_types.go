package types

// PDF Reference 1.4, Table 5.7 Font types

type FontType Name

const (
	FontType_Type0    FontType = "Type0"
	FontType_Type1    FontType = "Type1"
	FontType_Type3    FontType = "Type3"
	FontType_TrueType FontType = "TrueType"
	FontType_CIDFont  FontType = "CIDFont"
)

func (q FontType) ToRawBytes() []byte {
	return Name(q).ToRawBytes()
}

func (q FontType) Equal(obj Object) bool {
	a, ok := obj.(FontType)
	if !ok {
		return false
	}
	return q == a
}

type FontSubType Name

const (
	FontSub_Type0        FontSubType = "Type0"
	FontSub_Type1        FontSubType = "Type1"
	FontSub_MMType1      FontSubType = "MMType1"
	FontSub_Type3        FontSubType = "Type3"
	FontSub_TrueType     FontSubType = "TrueType"
	FontSub_CIDFontType0 FontSubType = "CIDFontType0"
	FontSub_CIDFontType2 FontSubType = "CIDFontType2"
)

func (q FontSubType) ToRawBytes() []byte {
	return Name(q).ToRawBytes()
}

func (q FontSubType) Copy(_ func(reference Reference) Reference) Object {
	return q
}
func (q FontSubType) Equal(obj Object) bool {
	a, ok := obj.(FontSubType)
	if !ok {
		return false
	}
	return q == a
}
