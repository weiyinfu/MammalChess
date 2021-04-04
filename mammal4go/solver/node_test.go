package solver

import (
	"fmt"
	"github.com/weiyinfu/MammalChess/mammal4go/util"
	"math/rand"
	"testing"
)

func check(node *Node) {
	//校验结点是否正确
	x := NewNode(node.Board, node.Unknown.Get())
	util.Assert(x.P[0].Eq(&node.P[0]), fmt.Sprintf("P[0] should be equal %v %v", x.P[0], node.P[0]))
	util.Assert(x.P[1].Eq(&node.P[1]), fmt.Sprintf("P[1] should be equal %v %v", x.P[1], node.P[1]))
	util.Assert(util.ArrayEqual(x.Status, node.Status), fmt.Sprintf("status should be equal"))
	util.Assert(x.Unknown.Eq(&node.Unknown), "unknown should be equal")
	util.Assert(x.UnknownPos.Eq(&node.UnknownPos), "unknownPos should be equal")
}
func TestNewNode(t *testing.T) {
	//把一个node从初始状态出发，一直生成着法，然后执行着法，最后到无法生成着法的时候
	var board []int
	var unknown []int
	for i := 0; i < 16; i++ {
		board = append(board, CHESS_UNKNOWN)
		unknown = append(unknown, i)
	}
	node := NewNode(board, unknown)
	turn := 0
	var realMoves []*Move
	round := 0
	for {
		round += 1
		if round > 100000 {
			fmt.Println("无穷无尽")
			break
		}
		moveList := node.generateMoves(turn, GENERATE_ALL)
		if len(moveList) == 0 {
			break
		}
		ind := rand.Intn(len(moveList))
		mo := moveList[ind]
		//fmt.Println(node.repr())
		//fmt.Println("before===")
		//fmt.Println("doing move", mo.repr())
		node.do(mo)
		//fmt.Println(node.repr())
		check(node)
		realMoves = append(realMoves, mo)
		turn = 1 - turn
	}
	fmt.Println(node.repr())
	fmt.Println("总共进行的实验步数", len(realMoves))
	for i := len(realMoves) - 1; i >= 0; i-- {
		node.undo(realMoves[i])
	}
	fmt.Println(node.repr())
}
func TestNewNode2(t *testing.T) {
	var board []int
	var unknown []int
	for i := 0; i < 16; i++ {
		board = append(board, CHESS_UNKNOWN)
		unknown = append(unknown, i)
	}
	node := NewNode(board, unknown)
	allMoves := node.generateMoves(0, GENERATE_ALL)
	fmt.Println(len(allMoves))
}

func TestNewNode3(t *testing.T) {
	var board = []int{3, 17, 2, 17,
		17, 4, 17, 13,
		16, 17, 17, 15,
		14, 17, 11, 9}
	node := NewNode(board, []int{})
	allMoves := node.generateMoves(0, GENERATE_ALL)
	fmt.Println(node.repr())
	fmt.Println(len(allMoves))
	for _, m := range allMoves {
		fmt.Println(node.showMove(m))
	}
}

func TestNewNode4(t *testing.T) {
	for ind, v := range DES_MAP {
		fmt.Println(ind, v)
	}
}
