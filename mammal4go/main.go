package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weiyinfu/MammalChess/mammal4go/solver"
	"github.com/weiyinfu/MammalChess/mammal4go/util"
	"log"
	"net/http"
	"strconv"
)

/**
使用golang实现超级深的搜索，打造无敌AI
*/
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
		computorColor, has := context.GetQuery("computerColor")
		if !has {
			panic("lack computerColor")
		}
		unknown, has := context.GetQuery("unknown")
		if !has {
			panic("cannot find unKnown")
		}
		board := util.GetIntList(a)
		computerColorI, err := strconv.Atoi(computorColor)
		unKnown := util.GetIntList(unknown)
		if err != nil {
			panic("解析computerColor失败")
		}
		log.Print(fmt.Sprintf("board=%v,unknown=%v,computerColor=%v", board, unKnown, computerColorI))
		//所使用的AI
		ai := solver.NewAlphaBetaSolver(10, 1)
		resp := ai.Solve(board, unKnown, computerColorI)
		context.JSON(http.StatusOK, resp)
	})
	_ = x.Run("0.0.0.0:7788")
}
