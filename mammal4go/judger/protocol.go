package judger

var DIRECTIONS [][]int

const (
	CHESS_SPACE   = 17
	CHESS_UNKNOWN = 16
)

func init() {
	DIRECTIONS = [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
}

type Solver interface {
	Solve(a []int, unknown []int, computerColor int) []int
}
