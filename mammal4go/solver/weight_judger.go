package solver

var CHESS_WEIGHT = []float64{1280, 640, 320, 160, 80, 40, 20, 10}

type WeightJudger struct {
}

func (WeightJudger) Judge(x *Node) float64 {
	//评价局面，总是表示0方的好坏，如果评价1方的好坏，添加符号即可
	smallScore := 0.0
	bigScore := 0.0
	for chess, status := range x.Status {
		if status == STATUS_UNKNOWN || status == STATUS_LIVE {
			if chess < 8 {
				smallScore += CHESS_WEIGHT[chess%8]
			} else {
				bigScore += CHESS_WEIGHT[chess%8]
			}
		}
	}
	if smallScore == 0 {
		//一个棋子都没了，必输
		return -INF
	} else if bigScore == 0 {
		//对面一个棋子都没有了，必胜
		return INF
	}
	score := smallScore - bigScore
	return score
}
