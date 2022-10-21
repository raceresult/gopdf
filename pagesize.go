package gopdf

// PageSize is width x height
type PageSize [2]Length

// StandardPageSize is used for GetStandardPageSize to get a standard page size like DIN-A4
type StandardPageSize int

const (
	PageSizeA0 StandardPageSize = iota
	PageSizeA1
	PageSizeA2
	PageSizeA3
	PageSizeA4
	PageSizeA5
	PageSizeA6
	PageSizeA7
	PageSizeLetter
	PageSizeLegal
)

// GetStandardPageSize returns a standard page size like DIN-A4 in either portrait or landscape
func GetStandardPageSize(size StandardPageSize, landscape bool) PageSize {
	var x PageSize
	switch size {
	case PageSizeA0:
		x = PageSize{MM(841), MM(1189)}
	case PageSizeA1:
		x = PageSize{MM(594), MM(841)}
	case PageSizeA2:
		x = PageSize{MM(420), MM(594)}
	case PageSizeA3:
		x = PageSize{MM(297), MM(420)}
	case PageSizeA4:
		x = PageSize{MM(210), MM(297)}
	case PageSizeA5:
		x = PageSize{MM(148), MM(210)}
	case PageSizeA6:
		x = PageSize{MM(105), MM(148)}
	case PageSizeA7:
		x = PageSize{MM(74), MM(105)}
	case PageSizeLetter:
		x = PageSize{Inch(8.5), Inch(11)}
	case PageSizeLegal:
		x = PageSize{Inch(8.5), Inch(14)}
	}

	if landscape {
		x[0], x[1] = x[1], x[0]
	}
	return x
}
