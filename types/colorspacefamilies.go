package types

// PDF Reference 1.4, Table 4.12 Color space families

type ColorSpaceFamily Name

const (
	ColorSpace_DeviceGray = ColorSpaceFamily("DeviceGray")
	ColorSpace_DeviceRGB  = ColorSpaceFamily("DeviceRGB")
	ColorSpace_DeviceCMYK = ColorSpaceFamily("DeviceCMYK")

	ColorSpace_CalGray  = ColorSpaceFamily("CalGray")
	ColorSpace_CalRGB   = ColorSpaceFamily("CalRGB")
	ColorSpace_Lab      = ColorSpaceFamily("Lab")
	ColorSpace_ICCBased = ColorSpaceFamily("ICCBased")

	ColorSpace_Indexed    = ColorSpaceFamily("Indexed")
	ColorSpace_Pattern    = ColorSpaceFamily("Pattern")
	ColorSpace_Separation = ColorSpaceFamily("Separation")
	ColorSpace_DeviceN    = ColorSpaceFamily("DeviceN")
)

func (q ColorSpaceFamily) ToRawBytes() []byte {
	return Name(q).ToRawBytes()
}

func (q ColorSpaceFamily) Copy(_ func(reference Reference) Reference) Object {
	return q
}
