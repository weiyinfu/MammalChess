package solver

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestCanEat(t *testing.T) {
	a := []int{15, 16, 8, 16, 1, 3, 0, 10, 11, 2, 12, 16, 5, 16, 16, 14}
	fmt.Println(a[12], a[8])
	res := CanEat(a[12], a[8])
	can := CanEat(a[8], a[12])
	fmt.Println(res)
	assert.Equal(t, res, false)
	assert.Equal(t, can, true)
}
func TestCanEat2(t *testing.T) {
	fmt.Println(CanEat(0, 15))
}
