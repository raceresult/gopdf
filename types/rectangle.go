package types

import "errors"

// PDF Reference 1.4, 3.8.3 Rectangles

type Rectangle struct {
	LLX Number
	LLY Number
	URX Number
	URY Number
}

func (q Rectangle) ToRawBytes() []byte {
	a := Array{q.LLX, q.LLY, q.URX, q.URY}
	return a.ToRawBytes()
}

func (q *Rectangle) Read(arr Array) error {
	// check length
	if len(arr) != 4 {
		return errors.New("rectangle requires array with 4 values")
	}

	// LLX
	var ok bool
	q.LLX, ok = arr[0].(Number)
	if !ok {
		v, ok := arr[0].(Int)
		if !ok {
			return errors.New("rectangle value invalid")
		}
		q.LLX = Number(v)
	}

	// LLY
	q.LLY, ok = arr[1].(Number)
	if !ok {
		v, ok := arr[1].(Int)
		if !ok {
			return errors.New("rectangle value invalid")
		}
		q.LLY = Number(v)
	}

	// URX
	q.URX, ok = arr[2].(Number)
	if !ok {
		v, ok := arr[2].(Int)
		if !ok {
			return errors.New("rectangle value invalid")
		}
		q.URX = Number(v)
	}

	// URY
	q.URY, ok = arr[3].(Number)
	if !ok {
		v, ok := arr[3].(Int)
		if !ok {
			return errors.New("rectangle value invalid")
		}
		q.URY = Number(v)
	}

	// return without error
	return nil
}

func (q Rectangle) Copy(copyRef func(reference Reference) Reference) Object {
	return Rectangle{
		LLX: q.LLX.Copy(copyRef).(Number),
		LLY: q.LLY.Copy(copyRef).(Number),
		URX: q.URX.Copy(copyRef).(Number),
		URY: q.URY.Copy(copyRef).(Number),
	}
}
