package alphabeta

import (
	"github.com/weiyinfu/MammalChess/mammal4go/judger"
	"math"
	"math/rand"
	"sort"
)

/**
TODO：在搜索中引入状态记忆，使用LRU记住最长遇到的若干状态，便于剪枝
这个空间越大，对时间提升就越大
*/
type AlphaBetaSolver struct {
	move       []Move
	eye        WeightJudger
	maxDepth   int //最大搜索深度
	maxRecover int //每条搜索路径上最大翻牌次数，翻牌次数太多会导致搜索变慢
}

func NewAlphaBetaSolver(maxDepth int, maxRecover int) *AlphaBetaSolver {
	return &AlphaBetaSolver{
		maxDepth:   maxDepth,
		maxRecover: maxRecover,
	}
}
func (self *AlphaBetaSolver) Solve(board []int, unknown []int, computerColor int) []int {
	no := NewNode(board, unknown)
	self.move = nil //运行之前，清空move记录
	self.dfs(no, 0, computerColor, -1e9, 1e9, 0, self.maxDepth)
	if len(self.move) == 0 {
		return nil
	}
	ind := rand.Intn(len(self.move))
	m := self.move[ind]
	if m.ty == MOVE_EAT {
		return []int{m.src, m.des}
	} else {
		return []int{m.src}
	}
}
func matai(score float64, depth int) float64 {
	/**
	 * 使用深度作为马太效应的因子：如果score是正数，表明局面对它有利，那么深度越浅越好
	 * 如果score是负数，表明局面对它不利，那么深度越深越好
	 * */
	if score < 0 {
		return score + float64(depth)
	}
	return score - float64(depth)
}
func (self *AlphaBetaSolver) dfs(x *Node, depth int, who int, lower float64, upper float64, recover int, maxDepth int) float64 {
	//recover表示翻开牌的个数
	//计算分数
	//fmt.Printf("depth=%v,maxDepth=%v recover=%v maxRecover=%v\n", depth, maxDepth, recover, self.maxRecover)
	if depth >= maxDepth {
		//搜索的深度到了，则直接使用eye进行评价
		score := self.eye.Judge(x)
		if who == 1 {
			//如果轮到1走了，那么返回的结果应该是1有多么的好
			score *= -1
		}
		score = matai(score, depth)
		return score
	}
	currentScore := self.eye.Judge(x)
	if math.Abs(currentScore) > INF {
		//如果局面已经是终局状态
		if who == 1 {
			currentScore *= -1
		}
		currentScore = matai(currentScore, depth)
		return currentScore
	}
	//如果已经翻了两张牌，则禁止生成翻牌，避免翻牌太多导致的指数爆炸
	generateMode := GENERATE_SIMPLE_UNKNOWN
	if recover == self.maxRecover {
		generateMode = GENERATE_EAT
	}
	moves := x.generateMoves(who, generateMode)
	if len(moves) == 0 {
		//无计可施
		return matai(-INF, depth)
	}
	sortMoves(x, moves)
	score := -INF
	for _, mo := range moves {
		//当前移动所产生的的子局面的分数
		var sonScore float64
		if mo.ty == MOVE_EAT {
			x.do(mo)
			//杀招裁剪，如果是吃子，则可以继续往下搜索
			nextMaxDepth := maxDepth
			if mo.eat != judger.CHESS_SPACE {
				//如果不是走向空位
				if depth+1 == nextMaxDepth {
					nextMaxDepth++
				}
			}
			sonScore = self.dfs(x, depth+1, 1-who, -upper, -score, recover, nextMaxDepth)
			x.undo(mo)
		} else {
			s := 0.0
			for _, be := range x.Unknown.Get() {
				mo.moving = be
				x.do(mo)
				//翻牌之后，继续向下搜索，不能立即停止
				nextMaxDepth := maxDepth
				if depth+1 == maxDepth {
					nextMaxDepth++
				}
				s += self.dfs(x, depth+1, 1-who, -upper, -score, recover+1, nextMaxDepth)
				x.undo(mo)
			}
			sonScore = s / float64(x.Unknown.Size())
		}
		if -sonScore >= score {
			if depth == 0 {
				//记录着法
				if -sonScore == score {
					//追平
					self.move = append(self.move, *mo)
				} else {
					//创建了新的着法
					self.move = []Move{*mo}
				}
			}
			if -sonScore > score {
				score = -sonScore
				//score总是追求越大越好
				if score >= upper {
					//执行剪枝，你的分数太高了，对手一定不会让你达到
					break
				}
			}
		}
	}
	return score
}

//TODO:实现迭代加深搜索，迭代加深搜索+大空间是完美的，迭代加深因为有记忆所以可以节省时间
func sortMoves(x *Node, moves []*Move) {
	//排序能够提升剪枝效率，优先吃大子，优先逃大子，其次翻牌
	sort.Slice(moves, func(i, j int) bool {
		a := moves[i]
		b := moves[j]
		if a.ty == b.ty {
			if a.ty == MOVE_NEW {
				//两个都是翻新牌
				return true
			} else {
				//两个都是吃子，谁吃的子大谁优先级高
				//谁吃掉的东西重要让谁先走
				desDif := a.eat - b.eat
				srcDif := a.moving - b.moving
				//如果同样是吃子，让子力强的子先吃子
				//如果同样是走空步，让子力强的子先走空步
				if desDif == 0 {
					return srcDif < 0
				} else {
					return desDif < 0
				}
			}
		} else {
			//优先吃子和动子，其次翻新牌
			if a.ty == MOVE_NEW {
				return false
			}
			return true
		}
	})
}
