package alphabeta_test

import (
	"fmt"
	"github.com/weiyinfu/MammalChess/mammal4go/judger/alphabeta"
	"testing"
)

func TestIntSet_Add(t *testing.T) {
	x := alphabeta.NewIntSet(4, 4)
	x.Add(3)
	fmt.Println(x.Get(), x)
	x.Add(2)
	fmt.Println(x.Get(), x)
	x.Remove(3)
	x.Add(1)
	x.Remove(1)
	fmt.Println(x.Get(), x)
}
