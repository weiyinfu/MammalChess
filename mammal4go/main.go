package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weiyinfu/MammalChess/mammal4go/judger"
	"github.com/weiyinfu/MammalChess/mammal4go/judger/alphabeta"
	"log"
	"strconv"
	"strings"
)

/**
使用golang实现超级深的搜索，打造无敌AI
*/

var ai judger.Solver

func init() {
	//ai = judger.RandSolver{}
	ai = alphabeta.NewAlphaBetaSolver(10, 1)
}
func getIntList(s string) []int {
	var a []int
	for _, i := range strings.Fields(strings.ReplaceAll(s, ",", " ")) {
		v, err := strconv.Atoi(i)
		if err != nil {
			panic(fmt.Sprintf("invalid param err=%v i=%v", err, i))
		}
		a = append(a, v)
	}
	return a
}
func main() {
	x := gin.Default()
	x.GET("/test", func(context *gin.Context) {
		_, _ = context.Writer.WriteString("MammalChess AI")
	})
	x.GET("/solve", func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		a, has := context.GetQuery("a")
		if !has {
			panic("lack parameter a")
		}
		id, has := context.GetQuery("id")
		if !has {
			panic("lack parameter id")
		}
		computorColor, has := context.GetQuery("computerColor")
		if !has {
			panic("lack computerColor")
		}
		unknown, has := context.GetQuery("unknown")
		if !has {
			panic("cannot find unKnown")
		}
		board := getIntList(a)
		computerColorI, err := strconv.Atoi(computorColor)
		unKnown := getIntList(unknown)
		if err != nil {
			panic("解析computerColor失败")
		}
		log.Print(fmt.Sprintf("board=%v,unknown=%v,computerColor=%v", board, unKnown, computerColorI))
		strategy := ai.Solve(board, unKnown, computerColorI)
		ans := map[string]interface{}{
			"strategy": strategy,
			"id":       id,
		}
		ansString, err := json.Marshal(ans)
		if err != nil {
			panic("序列化答案失败")
		}
		_, _ = context.Writer.Write(ansString)
	})
	_ = x.Run("0.0.0.0:7788")
}
