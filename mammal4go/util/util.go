package util

func Assert(x bool, desc string) {
	if x {

	} else {
		panic(desc)
	}
}

func AssertEqual(x interface{}, y interface{}, desc string) {
	Assert(x == y, desc)
}

func ArrayEqual(a, b []int) bool {
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func ArrayCopy(a []int) []int {
	var b []int
	for _, v := range a {
		b = append(b, v)
	}
	return b
}
