package solver

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestAlphaBetaSolver_Solve(t *testing.T) {
	ai := NewAlphaBetaSolver(3, 1)
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

	for {
		mo := ai.Solve(board, unknown, turn)
		if len(mo) == 0 {
			break
		}
		//fmt.Println(mo)
		//fmt.Println(node.repr())
		//fmt.Println("before===")
		var move *Move
		if len(mo) == 1 {
			//new
			if len(unknown) == 0 {
				//已经没有棋可翻了
				panic("invalid move")
			}
			randInd := rand.Intn(len(unknown))
			chess := unknown[randInd]
			//从unknown中把翻出来的元素删除掉
			unknown = append(unknown[:randInd], unknown[randInd+1:]...)
			move = &Move{
				Type:   MOVE_NEW,
				Who:    turn,
				Src:    mo[0],
				Moving: chess,
			}
		} else {
			move = &Move{
				Type:   MOVE_EAT,
				Who:    turn,
				Src:    mo[0],
				Des:    mo[1],
				Moving: board[mo[0]],
				Eat:    board[mo[1]],
			}
		}
		node.do(move)
		//fmt.Println(node.repr())
		check(node)
		realMoves = append(realMoves, move)
		turn = 1 - turn
	}
	fmt.Println(node.repr())
	fmt.Println("总共进行的实验步数", len(realMoves))
	for i := len(realMoves) - 1; i >= 0; i-- {
		node.undo(realMoves[i])
	}
	fmt.Println(node.repr())
}
func TestNewAlphaBetaSolver(t *testing.T) {
	const s = `a=15,4,17,16,11,17,17,7,10,17,1,16,8,0,17,17&id=29779786&computerColor=0&unknown=3,13`
	node := NewNodeFromString(s)
	ai := NewAlphaBetaSolver(10, 1)
	ans := ai.SolveNode(node, 0)
	fmt.Println(node.repr())
	fmt.Println("bes solution is", node.showMove(ans))
	for move, score := range ai.moves {
		fmt.Println(node.showMove(&move), "=>", score)
	}
}

func TestNewAlphaBetaSolver2(t *testing.T) {
	const s = `a=15,1,13,9,17,5,0,17,8,17,17,3,17,17,17,17&id=850230011&computerColor=0&unknown=`
	node := NewNodeFromString(s)
	ai := NewAlphaBetaSolver(10, 1)
	ans := ai.SolveNode(node, 0)
	fmt.Println(node.repr())
	fmt.Println("bes solution is", node.showMove(ans))
	for move, score := range ai.moves {
		fmt.Println(node.showMove(&move), "=>", score)
	}
}
func TestAlphaBetaSolver_Solve2(t *testing.T) {
	const s = `a=15,1,13,9,17,5,0,17,8,17,17,3,17,17,17,17&id=850230011&computerColor=0&unknown=`
	node := NewNodeFromString(s)
	moves := node.generateMoves(0, GENERATE_ALL)
	sortMoves(node, moves)
	fmt.Println(node.repr())
	for _, m := range moves {
		fmt.Println(node.showMove(m))
	}
}
