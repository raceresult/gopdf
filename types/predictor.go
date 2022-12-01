package types

import (
	"bytes"
	"errors"
)

// PDF Reference 1.4, Table 3.8 Predictor values

const (
	PredictorNo      = 1  // No prediction.
	PredictorTIFF    = 2  // Use TIFF prediction for all rows.
	PredictorNone    = 10 // Use PNGNone for all rows.
	PredictorSub     = 11 // Use PNGSub for all rows.
	PredictorUp      = 12 // Use PNGUp for all rows.
	PredictorAverage = 13 // Use PNGAverage for all rows.
	PredictorPaeth   = 14 // Use PNGPaeth for all rows.
	PredictorOptimum = 15 // Use the optimum PNG prediction for each row.
)

// For predictor > 2 PNG filters (see RFC 2083) get applied and the first byte of each pixelrow defines
// the prediction algorithm used for all pixels of this row.
const (
	PNGNone    = 0x00
	PNGSub     = 0x01
	PNGUp      = 0x02
	PNGAverage = 0x03
	PNGPaeth   = 0x04
)

func decodePredictor(data []byte, params FilterParameters) ([]byte, error) {
	bytesPerPixel := (params.BitsPerComponent*params.Colors + 7) / 8
	rowSize := int(params.BitsPerComponent * params.Colors * params.Columns / 8)
	if params.Predictor != PredictorTIFF {
		// PNG prediction uses a row filter byte prefixing the pixelbytes of a row.
		rowSize++
	}

	// cr and pr are the bytes for the current and previous row.
	cr := make([]byte, rowSize)
	pr := make([]byte, rowSize)

	// Output buffer
	var b bytes.Buffer
	for len(data) != 0 {
		// Read decompressed bytes for one pixel row.
		if len(data) < rowSize {
			return nil, errors.New("EOF")
		}
		cr = data[:rowSize]
		data = data[rowSize:]

		d, err1 := processRow(pr, cr, int(params.Predictor), int(params.Colors), int(bytesPerPixel))
		if err1 != nil {
			return nil, err1
		}

		if _, err := b.Write(d); err != nil {
			return nil, err
		}

		// Swap byte slices.
		pr, cr = cr, pr
	}

	if b.Len()%int(params.BitsPerComponent*params.Colors*params.Columns/8) > 0 {
		return nil, errors.New("filter FlateDecode: postprocessing failed")
	}
	return b.Bytes(), nil
}

func processRow(pr, cr []byte, p, colors, bytesPerPixel int) ([]byte, error) {
	if p == PredictorTIFF {
		return applyHorDiff(cr, colors)
	}

	// Apply the filter.
	cdat := cr[1:]
	pdat := pr[1:]

	// Get row filter from 1st byte
	f := int(cr[0])

	// The value of Predictor supplied by the decoding filter need not match the value
	// used when the data was encoded if they are both greater than or equal to 10.

	switch f {

	case PNGNone:
		// No operation.

	case PNGSub:
		for i := bytesPerPixel; i < len(cdat); i++ {
			cdat[i] += cdat[i-bytesPerPixel]
		}

	case PNGUp:
		for i, p := range pdat {
			cdat[i] += p
		}

	case PNGAverage:
		// The average of the two neighboring pixels (left and above).
		// Raw(x) - floor((Raw(x-bpp)+Prior(x))/2)
		for i := 0; i < bytesPerPixel; i++ {
			cdat[i] += pdat[i] / 2
		}
		for i := bytesPerPixel; i < len(cdat); i++ {
			cdat[i] += uint8((int(cdat[i-bytesPerPixel]) + int(pdat[i])) / 2)
		}

	case PNGPaeth:
		filterPaeth(cdat, pdat, bytesPerPixel)

	}

	return cdat, nil
}

func applyHorDiff(row []byte, colors int) ([]byte, error) {
	// This works for 8 bits per color only.
	for i := 1; i < len(row)/colors; i++ {
		for j := 0; j < colors; j++ {
			row[i*colors+j] += row[(i-1)*colors+j]
		}
	}
	return row, nil
}

// intSize is either 32 or 64.
// Disabled intSize 64 for govet.
const intSize = 32 //<< (^uint(0) >> 63)

func abs(x int) int {
	// m := -1 if x < 0. m := 0 otherwise.
	m := x >> (intSize - 1)

	// In two's complement representation, the negative number
	// of any number (except the smallest one) can be computed
	// by flipping all the bits and add 1. This is faster than
	// code with a branch.
	// See Hacker's Delight, section 2-4.
	return (x ^ m) - m
}

// paeth implements the Paeth filter function, as per the PNG specification.
func paeth(a, b, c uint8) uint8 {
	// This is an optimized version of the sample code in the PNG spec.
	// For example, the sample code starts with:
	//	p := int(a) + int(b) - int(c)
	//	pa := abs(p - int(a))
	// but the optimized form uses fewer arithmetic operations:
	//	pa := int(b) - int(c)
	//	pa = abs(pa)
	pc := int(c)
	pa := int(b) - pc
	pb := int(a) - pc
	pc = abs(pa + pb)
	pa = abs(pa)
	pb = abs(pb)
	if pa <= pb && pa <= pc {
		return a
	} else if pb <= pc {
		return b
	}
	return c
}

// filterPaeth applies the Paeth filter to the cdat slice.
// cdat is the current row's data, pdat is the previous row's data.
func filterPaeth(cdat, pdat []byte, bytesPerPixel int) {
	var a, b, c, pa, pb, pc int
	for i := 0; i < bytesPerPixel; i++ {
		a, c = 0, 0
		for j := i; j < len(cdat); j += bytesPerPixel {
			b = int(pdat[j])
			pa = b - c
			pb = a - c
			pc = abs(pa + pb)
			pa = abs(pa)
			pb = abs(pb)
			if pa <= pb && pa <= pc {
				// No-op.
			} else if pb <= pc {
				a = b
			} else {
				a = c
			}
			a += int(cdat[j])
			a &= 0xff
			cdat[j] = uint8(a)
			c = b
		}
	}
}
