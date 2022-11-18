package types

// PDF Reference 1.4, 3.5 Standard filters

type Filter Name

const (
	Filter_ASCIIHexDecode  Filter = "ASCIIHexDecode"
	Filter_ASCII85Decode   Filter = "ASCII85Decode"
	Filter_LZWDecode       Filter = "LZWDecode"
	Filter_FlateDecode     Filter = "FlateDecode"
	Filter_RunLengthDecode Filter = "RunLengthDecode"
	Filter_CCITTFaxDecode  Filter = "CCITTFaxDecode"
	Filter_JBIG2Decode     Filter = "JBIG2Decode"
	Filter_DCTDecode       Filter = "DCTDecode"
)

func (q Filter) ToRawBytes() []byte {
	return Name(q).ToRawBytes()
}

func (q Filter) Copy(_ func(reference Reference) Reference) Object {
	return q
}
