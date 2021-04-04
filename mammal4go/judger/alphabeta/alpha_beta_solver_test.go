package alphabeta

import (
	"fmt"
	"github.com/weiyinfu/MammalChess/mammal4go/judger"
	"math/rand"
	"testing"
)

func TestAlphaBetaSolver_Solve(t *testing.T) {
	ai := NewAlphaBetaSolver(3, 1)
	//把一个node从初始状态出发，一直生成着法，然后执行着法，最后到无法生成着法的时候
	var board []int
	var unknown []int
	for i := 0; i < 16; i++ {
		board = append(board, judger.CHESS_UNKNOWN)
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
				ty:     MOVE_NEW,
				who:    turn,
				src:    mo[0],
				moving: chess,
			}
		} else {
			move = &Move{
				ty:     MOVE_EAT,
				who:    turn,
				src:    mo[0],
				des:    mo[1],
				moving: board[mo[0]],
				eat:    board[mo[1]],
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
