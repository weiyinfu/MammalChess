package alphabeta

import (
	"fmt"
	"github.com/weiyinfu/MammalChess/mammal4go/judger"
	"github.com/weiyinfu/MammalChess/mammal4go/util"
)

const (
	INF = 1e9
)

var DES_MAP [][]int //预先初始化每个src处的棋子可以去往的位置

func init() {
	DES_MAP = make([][]int, 16)
	for pos, _ := range DES_MAP {
		for _, d := range judger.DIRECTIONS {
			x, y := pos/4, pos%4
			xx, yy := x+d[0], y+d[1]
			if !(xx >= 0 && yy >= 0 && xx < 4 && yy < 4) {
				continue
			}
			DES_MAP[pos] = append(DES_MAP[pos], xx*4+yy)
		}
	}
}

const (
	STATUS_DIED    = 0
	STATUS_LIVE    = 1
	STATUS_UNKNOWN = 2
)
const (
	MOVE_EAT = 0
	MOVE_NEW = 1
)

type Node struct {
	Board      []int
	Status     []int    //16枚棋子的状态
	P          []IntSet //a棋子所在的位置，P[0]表示一类棋子的，P[1]表示另一类棋子
	UnknownPos IntSet   //未知棋子所在的位置
	Unknown    IntSet   //未知棋子的形状
}

type Move struct {
	ty     int //移动的类型，可能是吃子
	who    int //谁在操作
	src    int
	des    int
	eat    int
	moving int //如果是翻开新棋，则表明翻开的棋子是什么
}

func (m Move) repr() interface{} {
	if m.ty == MOVE_EAT {
		return fmt.Sprintf(`%v移动棋子从%v到%v`, m.who, m.src, m.des)
	}
	return fmt.Sprintf(`%v翻开位于%v处的新棋%v`, m.who, m.src, m.moving)
}
func NewNode(board []int, unknown []int) *Node {
	p := []IntSet{*NewIntSet(16, 16), *NewIntSet(16, 16)}
	unknownPos := NewIntSet(16, 16)
	unknownChess := NewIntSet(16, 16)
	//初始化status
	status := make([]int, 16)
	for pos, chess := range board {
		if chess == judger.CHESS_UNKNOWN {
			//为空
			unknownPos.Add(pos)
			continue
		}
		if chess == judger.CHESS_SPACE {
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
func (x *Node) do(m *Move) {
	if m.ty == MOVE_EAT {
		//移动棋子
		util.AssertEqual(m.moving/8, m.who, "who is moving")
		x.Board[m.src] = judger.CHESS_SPACE
		x.Board[m.des] = m.moving
		x.P[m.who].Remove(m.src)
		x.P[m.who].Add(m.des)
		if m.eat != judger.CHESS_SPACE {
			//如果吃掉了棋子，需要改变status
			x.Status[m.eat] = STATUS_DIED
			x.P[1-m.who].Remove(m.des)
		}
	} else {
		//翻开新棋子
		x.Board[m.src] = m.moving
		x.Status[m.moving] = STATUS_LIVE
		whose := 0
		if m.moving >= 8 {
			whose = 1
		}
		x.P[whose].Add(m.src)
		x.UnknownPos.Remove(m.src)
		x.Unknown.Remove(m.moving)
	}
}

func (x *Node) undo(m *Move) {
	if m.ty == MOVE_EAT {
		x.Board[m.des] = m.eat
		x.Board[m.src] = m.moving
		x.P[m.who].Remove(m.des) //吃掉
		x.P[m.who].Add(m.src)
		if m.eat != judger.CHESS_SPACE {
			x.P[1-m.who].Add(m.des)
			x.Status[m.eat] = STATUS_LIVE
		}
	} else {
		x.Board[m.src] = judger.CHESS_UNKNOWN
		x.Status[m.moving] = STATUS_UNKNOWN
		whose := 0
		if m.moving >= 8 {
			whose = 1
		}
		x.P[whose].Remove(m.src)
		x.UnknownPos.Add(m.src)
		x.Unknown.Add(m.moving)
	}
}

//着法生成方式
const (
	GENERATE_ALL            = 1
	GENERATE_EAT            = 2 //不生成unknown着法
	GENERATE_SIMPLE_UNKNOWN = 3 //生成unnown，但是不展开
)

func (x *Node) generateMoves(who int, mode int) []*Move {
	//生成全部的着法
	var moveList []*Move
	for _, pos := range x.P[who].Get() {
		for _, des := range DES_MAP[pos] {
			if !judger.CanEat(x.Board[pos], x.Board[des]) {
				continue
			}
			mo := Move{
				ty:     MOVE_EAT,
				src:    pos,
				des:    des,
				who:    who,
				eat:    x.Board[des],
				moving: x.Board[pos],
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
				ty:     MOVE_NEW,
				src:    pos,
				who:    who,
				moving: -1, //moving在运行时根据unknown填充，避免此处枚举太多可能性浪费空间
			})
		}
		return moveList
	}
	util.AssertEqual(mode, GENERATE_ALL, "mode错误")
	for _, pos := range x.UnknownPos.Get() {
		for _, maybe := range x.Unknown.Get() {
			moveList = append(moveList, &Move{
				ty:     MOVE_NEW,
				src:    pos,
				who:    who,
				moving: maybe,
			})
		}
	}
	return moveList
}

func (x *Node) repr() string {
	//TODO:使用反射实现
	return fmt.Sprintf(`board=%v
status=%v
unknownPos=%v
unknown=%v
P[0]=%v
P[1]=%v
`, x.Board, x.Status, x.UnknownPos, x.Unknown, x.P[0], x.P[1])
}

type Judger interface {
	Judge(x *Node) float64
}
