package types

import "errors"

// PDF Reference 1.4, 3.5 Standard filters

type FilterParameters struct {
	// A code that selects the predictor algorithm, if any. If the value of this entry
	// is 1, the filter assumes that the normal algorithm was used to encode the data,
	// without prediction. If the value is greater than 1, the filter assumes that the
	// data was differenced before being encoded, and Predictor selects the predic-
	// tor algorithm. For more information regarding Predictor values greater
	// than 1, see “LZW and Flate Predictor Functions,” below. Default value: 1.
	Predictor Int

	// (Used only if Predictor is greater than 1) The number of interleaved color com-
	// ponents per sample. Valid values are 1 to 4 in PDF 1.2 or earlier and 1 or
	// greater in PDF 1.3 or later. Default value: 1.
	Colors Int

	// (Used only if Predictor is greater than 1) The number of bits used to represent
	// each color component in a sample. Valid values are 1, 2, 4, 8, and (in PDF 1.5)
	// 16. Default value: 8.
	BitsPerComponent Int

	// (Used only if Predictor is greater than 1) The number of samples in each row.
	// Default value: 1.
	Columns Int

	// (LZWDecode only) An indication of when to increase the code length. If the
	// value of this entry is 0, code length increases are postponed as long as pos-
	// sible. If the value is 1, code length increases occur one code early. This pa-
	// rameter is included because LZW sample code distributed by some vendors
	// increases the code length one code earlier than necessary. Default value: 1.
	EarlyChange Int
}

func (q *FilterParameters) Read(dict Dictionary) error {
	// set default values
	q.Predictor = 1
	q.Colors = 1
	q.BitsPerComponent = 8
	q.Columns = 1
	q.EarlyChange = 1

	// Predictor
	v, ok := dict["Predictor"]
	if ok {
		q.Predictor, ok = v.(Int)
		if !ok {
			return errors.New("DecodeParms value Predictor invalid")
		}
	}

	// Colors
	v, ok = dict["Colors"]
	if ok {
		q.Colors, ok = v.(Int)
		if !ok {
			return errors.New("DecodeParms value Colors invalid")
		}
	}

	// BitsPerComponent
	v, ok = dict["BitsPerComponent"]
	if ok {
		q.BitsPerComponent, ok = v.(Int)
		if !ok {
			return errors.New("DecodeParms value BitsPerComponent invalid")
		}
	}

	// Columns
	v, ok = dict["Columns"]
	if ok {
		q.Columns, ok = v.(Int)
		if !ok {
			return errors.New("DecodeParms value Columns invalid")
		}
	}

	// EarlyChange
	v, ok = dict["EarlyChange"]
	if ok {
		q.EarlyChange, ok = v.(Int)
		if !ok {
			return errors.New("DecodeParms value EarlyChange invalid")
		}
	}

	// return without error
	return nil
}
