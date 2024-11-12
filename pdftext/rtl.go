package pdftext

// reverseRTLString reverse the parts of the given string which are part of a right-to-left language
// only if the first character of the string is right-to-left
func reverseRTLString(s string) string {
	if s == "" {
		return s
	}

	isRTL := func(r rune) bool {
		return r >= 0x591 && r <= 0x6EF
	}

	rr := []rune(s)
	if !isRTL(rr[0]) {
		return s
	}

	b := true
	l := 0
	arr := make([]rune, 0, len(rr))
	for i, r := range rr {
		rtl := isRTL(r)
		if rtl && !b {
			arr = append(arr, rr[l:i]...)
			l = i
			b = !b
		} else if r >= 0x30 && r <= 0x39 && b || r >= 0x41 && r <= 0x5A && b || r >= 0x61 && r <= 0x7A && b {
			for x := i - 1; x >= l; x-- {
				arr = append(arr, rr[x])
			}
			l = i
			b = !b
		} else if i == len(rr)-1 {
			if rtl {
				for x := len(rr) - 1; x >= l; x-- {
					arr = append(arr, rr[x])
				}
			} else {
				arr = append(arr, rr[l:]...)
			}
		}
	}
	return string(arr)
}
