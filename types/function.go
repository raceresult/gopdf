package types

// PDF Reference 1.4, Table 3.26 Entries common to all function dictionaries

type Function struct {
	// (Required) The function type:
	// 0 Sampled function
	// 2 Exponential interpolation function
	// 3 Stitching function
	// 4 PostScript calculator function
	FunctionType Int

	// (Required) An array of 2 ◊ m numbers, where m is the number of input
	// values. For each i from 0 to m − 1, Domain2i must be less than or equal to
	// Domain2i+1 , and the ith input value, x i, must lie in the interval
	// Domain2i ≤ x i ≤ Domain2i+1 . Input values outside the declared domain are
	// clipped to the nearest boundary value.
	Domain Array

	// (Required for type 0 and type 4 functions, optional otherwise; see below) An
	// array of 2 ◊ n numbers, where n is the number of output values. For each j
	// from 0 to n − 1, Range2j must be less than or equal to Range2j+1 , and the jth
	// output value, y j , must lie in the interval Range2j ≤ y j ≤ Range2j+1 . Output
	// values outside the declared range are clipped to the nearest boundary value.
	// If this entry is absent, no clipping is done
	Range Array
}

func (q Function) ToRawBytes() []byte {
	d := Dictionary{
		"FunctionType": q.FunctionType,
		"Domain":       q.Domain,
	}
	if len(q.Range) != 0 {
		d["Range"] = q.Range
	}
	return d.ToRawBytes()
}

func (q Function) Copy(copyRef func(reference Reference) Reference) Object {
	return Function{
		FunctionType: q.FunctionType.Copy(copyRef).(Int),
		Domain:       q.Domain.Copy(copyRef).(Array),
		Range:        q.Range.Copy(copyRef).(Array),
	}
}

func (q Function) Equal(obj Object) bool {
	a, ok := obj.(Function)
	if !ok {
		return false
	}
	if !Equal(q.FunctionType, a.FunctionType) {
		return false
	}
	if !Equal(q.Domain, a.Domain) {
		return false
	}
	if !Equal(q.Range, a.Range) {
		return false
	}
	return true
}
