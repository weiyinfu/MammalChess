package solver

import (
	"fmt"
	"github.com/weiyinfu/MammalChess/mammal4go/util"
	"strings"
)

const (
	INF = 1e9
)

//两种着法：吃子，动子
const (
	MOVE_EAT = 0
	MOVE_NEW = 1
)

const (
	CHESS_SPACE   = 17
	CHESS_UNKNOWN = 16
)

//棋子的三种状态：死了、活着、未知
const (
	STATUS_DIED    = 0
	STATUS_LIVE    = 1
	STATUS_UNKNOWN = 2
)

var DES_MAP [][]int //预先初始化每个src处的棋子可以去往的位置

//judger：所有AI的定义
var DIRECTIONS [][]int

func init() {
	DIRECTIONS = [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	DES_MAP = make([][]int, 16)
	for pos, _ := range DES_MAP {
		for _, d := range DIRECTIONS {
			x, y := pos/4, pos%4
			xx, yy := x+d[0], y+d[1]
			if !(xx >= 0 && yy >= 0 && xx < 4 && yy < 4) {
				continue
			}
			DES_MAP[pos] = append(DES_MAP[pos], xx*4+yy)
		}
	}
}

type SolveResult struct {
	Strategy []int   `json:"strategy"`
	Score    float64 `json:"score"`
}
type Solver interface {
	Solve(a []int, unknown []int, computerColor int) *SolveResult
}

type Node struct {
	Board      []int         //4*4的棋盘，用一维数组进行表示
	Status     []int         //16枚棋子的状态
	P          []util.IntSet //各方棋子所在的位置，避免每次扫描整个棋盘。P[0]表示一类棋子的，P[1]表示另一类棋子
	UnknownPos util.IntSet   //未知棋子所在的位置
	Unknown    util.IntSet   //未知棋子的形状
}

type Move struct {
	Type   int //移动的类型，可能是吃子
	Who    int //谁在操作
	Src    int //移动的棋子的初始位置
	Des    int //移动的棋子去往的位置
	Eat    int //移动的棋子是啥
	Moving int //吃掉的棋子是啥，如果是翻开新棋，则表明翻开的棋子是什么
}

func NewNode(board []int, unknown []int) *Node {
	p := []util.IntSet{*util.NewIntSet(16), *util.NewIntSet(16)}
	unknownPos := util.NewIntSet(16)
	unknownChess := util.NewIntSet(16)
	//初始化status
	status := make([]int, 16)
	for pos, chess := range board {
		if chess == CHESS_UNKNOWN {
			//为空
			unknownPos.Add(pos)
			continue
		}
		if chess == CHESS_SPACE {
			continue
		}
		status[chess] = STATUS_LIVE
		who := 0
		if chess >= 8 {
			who = 1
		}
		p[who].Add(pos)
	}
	for _, chess := range unknown {
		status[chess] = STATUS_UNKNOWN
		unknownChess.Add(chess)
	}
	return &Node{
		Board:      board,
		Status:     status,
		P:          p,
		Unknown:    *unknownChess,
		UnknownPos: *unknownPos,
	}
}
func NewNodeFromString(s string) *Node {
	a := strings.Split(s, "&")
	ma := map[string]string{}
	for _, i := range a {
		ind := strings.Index(i, "=")
		k := i[:ind]
		v := i[ind+1:]
		ma[k] = v
	}
	board := util.GetIntList(ma["a"])
	unknown := util.GetIntList(ma["unknown"])
	return NewNode(board, unknown)
}
func (x *Node) do(m *Move) {
	if m.Type == MOVE_EAT {
		//移动棋子
		x.Board[m.Src] = CHESS_SPACE
		x.Board[m.Des] = m.Moving
		x.P[m.Who].Remove(m.Src)
		x.P[m.Who].Add(m.Des)
		if m.Eat != CHESS_SPACE {
			//如果吃掉了棋子，需要改变status
			x.Status[m.Eat] = STATUS_DIED
			x.P[1-m.Who].Remove(m.Des)
		}
	} else {
		//翻开新棋子
		x.Board[m.Src] = m.Moving
		x.Status[m.Moving] = STATUS_LIVE
		whose := 0
		if m.Moving >= 8 {
			whose = 1
		}
		x.P[whose].Add(m.Src)
		x.UnknownPos.Remove(m.Src)
		x.Unknown.Remove(m.Moving)
	}
}

func (x *Node) undo(m *Move) {
	if m.Type == MOVE_EAT {
		x.Board[m.Des] = m.Eat
		x.Board[m.Src] = m.Moving
		x.P[m.Who].Remove(m.Des) //吃掉
		x.P[m.Who].Add(m.Src)
		if m.Eat != CHESS_SPACE {
			x.P[1-m.Who].Add(m.Des)
			x.Status[m.Eat] = STATUS_LIVE
		}
	} else {
		x.Board[m.Src] = CHESS_UNKNOWN
		x.Status[m.Moving] = STATUS_UNKNOWN
		whose := 0
		if m.Moving >= 8 {
			whose = 1
		}
		x.P[whose].Remove(m.Src)
		x.UnknownPos.Add(m.Src)
		x.Unknown.Add(m.Moving)
	}
}

func (self *Node) showMove(m *Move) string {
	if m.Type == MOVE_EAT {
		if m.Eat == CHESS_SPACE {
			return fmt.Sprintf(`user%v移动棋子%v从%v到%v`, m.Who, getName(m.Moving), getPos(m.Src), getPos(m.Des))
		} else {
			return fmt.Sprintf(`user%v移动%v%v吃掉%v%v`, m.Who, getName(m.Moving), getPos(m.Src), getName(m.Eat), getPos(m.Des))
		}
	}
	return fmt.Sprintf(`user%v翻开位于%v处的新棋%v`, m.Who, m.Src, getName(m.Moving))
}

//着法生成方式
const (
	GENERATE_ALL            = 1
	GENERATE_EAT            = 2 //不生成unknown着法
	GENERATE_SIMPLE_UNKNOWN = 3 //生成unnown，但是不展开
)

func CanEat(animal, food int) bool {
	//给定两枚棋子，判断是否可以吃子
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
		if food == 7 {
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

func (x *Node) generateMoves(who int, mode int) []*Move {
	//生成全部的着法
	var moveList []*Move
	for _, pos := range x.P[who].Get() {
		for _, des := range DES_MAP[pos] {
			if !CanEat(x.Board[pos], x.Board[des]) {
				continue
			}
			mo := Move{
				Type:   MOVE_EAT,
				Src:    pos,
				Des:    des,
				Who:    who,
				Eat:    x.Board[des],
				Moving: x.Board[pos],
			}
			moveList = append(moveList, &mo)
		}
	}
	if mode == GENERATE_EAT {
		//如果只生成吃子着法
		return moveList
	}
	if mode == GENERATE_SIMPLE_UNKNOWN {
		for _, pos := range x.UnknownPos.Get() {
			moveList = append(moveList, &Move{
				Type:   MOVE_NEW,
				Src:    pos,
				Who:    who,
				Moving: -1, //moving在运行时根据unknown填充，避免此处枚举太多可能性浪费空间
			})
		}
		return moveList
	}
	util.AssertEqual(mode, GENERATE_ALL, "mode错误")
	for _, pos := range x.UnknownPos.Get() {
		for _, maybe := range x.Unknown.Get() {
			moveList = append(moveList, &Move{
				Type:   MOVE_NEW,
				Src:    pos,
				Who:    who,
				Moving: maybe,
			})
		}
	}
	return moveList
}
func getPos(pos int) string {
	return fmt.Sprintf("(%v,%v)", pos/4, pos%4)
}
func getName(chess int) string {
	if chess == CHESS_SPACE {
		return "."
	}
	if chess == CHESS_UNKNOWN {
		return "X"
	}
	Name := strings.Split("象狮豹虎狼狗猫鼠", "")
	name := fmt.Sprintf("%v", Name[chess%8])
	if chess < 8 {
		name = fmt.Sprintf("%v", chess+1)
	}
	return name
}
func (x *Node) repr() string {
	s := fmt.Sprintf(`board=%v
status=%v
unknownPos=%v
unknown=%v
P[0]=%v
P[1]=%v
`, x.Board, x.Status, x.UnknownPos, x.Unknown, x.P[0], x.P[1])
	line := ""
	for i, v := range x.Board {
		line += getName(v)
		if i%4 == 3 {
			line += "\n"
		}
	}
	return fmt.Sprintf("%v\n%v", s, line)
}
