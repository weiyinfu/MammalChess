package solver

import (
	"math/rand"
)

//随机走子的judger，比较简单的写法
type RandSolver struct {
}

func (RandSolver) Solve(board []int, unknown []int, computerColor int) *SolveResult {
	//随机走子
	moves := GetAllMoves(board, computerColor)
	if len(moves) == 0 {
		return nil
	}
	ind := rand.Intn(len(moves))
	return &SolveResult{Strategy: moves[ind], Score: 0}
}
func GetAllMoves(board []int, computerColor int) [][]int {
	var moves [][]int
	for ind, chess := range board {
		if chess == CHESS_UNKNOWN {
			moves = append(moves, []int{ind})
			continue
		}
		if chess == CHESS_SPACE {
			continue
		}
		if chess < 8 && computerColor == 0 || chess >= 8 && computerColor == 1 {
			moves = append(moves, canGo(board, ind)...)
		}
	}
	if len(moves) == 0 {
		return nil
	}
	return moves
}

func canGo(board []int, pos int) [][]int {
	//获取pos处的棋子所有可去的位置
	x, y := pos/4, pos%4
	var tos [][]int
	for _, d := range DIRECTIONS {
		xx, yy := x+d[0], y+d[1]
		if !(xx >= 0 && yy >= 0 && xx < 4 && yy < 4) {
			continue
		}
		to := xx*4 + yy
		if board[to] == CHESS_UNKNOWN {
			continue
		}
		if CanEat(board[pos], board[to]) {
			tos = append(tos, []int{pos, to})
		}
	}
	return tos
}
