package util

import (
	"fmt"
	"strconv"
	"strings"
)

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
func GetIntList(s string) []int {
	//s是一个逗号隔开的int列表组成的字符串
	var a []int
	for _, i := range strings.Fields(strings.ReplaceAll(s, ",", " ")) {
		v, err := strconv.Atoi(i)
		if err != nil {
			panic(fmt.Sprintf("invalid param err=%v i=%v", err, i))
		}
		a = append(a, v)
	}
	return a
}
