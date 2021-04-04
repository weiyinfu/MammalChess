package judger

import (
	"math/rand"
)

type RandSolver struct {
}

func (RandSolver) Solve(board []int, unknown []int, computerColor int) []int {
	//随机走子
	moves := GetAllMoves(board, computerColor)
	if len(moves) == 0 {
		return nil
	}
	ind := rand.Intn(len(moves))
	return moves[ind]
}
func GetAllMoves(board []int, computerColor int) [][]int {
	var moves [][]int
	for ind, chess := range board {
		if chess == 16 {
			moves = append(moves, []int{ind})
			continue
		}
		if chess == 17 {
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

func CanEat(animal, food int) bool {
	if food == CHESS_UNKNOWN {
		return false
	}
	if food == CHESS_SPACE {
		return true
	}
	if (animal < 8) == (food < 8) {
		//same color cannot eat
		return false
	}
	if food >= 8 {
		food -= 8
	}
	if animal >= 8 {
		animal -= 8
	}
	if animal == 7 {
		if food == 0 || food == 7 {
			return true
		}
		return false
	}
	if animal == 0 {
		if food == 8 {
			return false
		}
		//elephant can eat everything
		return true
	}
	//可以吃掉
	if animal <= food {
		return true
	}
	return false
}
