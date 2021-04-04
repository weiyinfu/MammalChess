/**
 * AI对比器：比较AI的水平
 * js版的ai比较器，使用同步的方式调用AI
 */
const Game = require('./game')
const comparer = {
    randomInt(down, up) {
        return down + Math.floor((up - down) * Math.random())
    },
    illegal(game, ans, message) {
        Game.show(game.a)
        console.log(ans)
        throw message
    },
    async compare(ai, bi, compareCount) {
        /*比较两个AI，打印AI比较结果
        * ai和bi两个参数都表示AI，compareCount表示对局次数
         */
        const ais = [ai, bi]
        const statistics = {
            win: [0, 0],
            peace: 0
        }
        for (let i = 0; i < compareCount; i++) {
            const game = Game.newGame()
            var turn = this.randomInt(0, 2)//turn表示谁是先手
            //让先手翻开牌从而决定每个人的颜色
            const first = this.randomInt(0, 16)
            game.recover(first)
            var computerColor = game.a[first] < 8 ? 0 : 1
            var winner = -1//这场游戏最终的胜利者
            let round = 0;
            while (true) {
                round += 1;
                turn = 1 - turn
                computerColor = 1 - computerColor
                const ans = await ais[turn].solve(game.a, game.unknown, computerColor)
                if (ans.strategy === null) {
                    //如果返回无解，则表示AI主动认输
                    winner = 1 - turn
                    break
                } else if (ans.strategy.length === 1) {
                    //AI翻开新牌
                    if (game.canRecover(ans.strategy[0]))
                        game.recover(ans.strategy[0])
                    else {
                        this.illegal(game, ans, 'AI ' + turn + " illegal operation")
                    }
                } else if (ans.strategy.length === 2) {
                    //AI移动棋子
                    const [src, des] = ans.strategy
                    const srcColor = game.a[src] < 8 ? 0 : 1
                    if (srcColor !== computerColor) {
                        this.illegal(game, ans, '只能移动自己颜色的棋子')
                    }
                    if (game.canEat(src, des)) {
                        game.eat(src, des)
                    } else {
                        this.illegal(game, ans, 'AI ' + turn + " illegal operation")
                    }
                }
                //判断游戏是否结束
                const state = game.getState()
                if (state === Game.PEACE) {
                    winner = -1
                    break
                } else if (state === Game.FIRST_WIN || state === Game.SECOND_WIN) {
                    //如果游戏结束，谁走得最后一步谁就赢了
                    if (state === Game.FIRST_WIN) {
                        winner = computerColor === 0 ? turn : 1 - turn
                    } else {
                        winner = computerColor === 1 ? turn : 1 - turn
                    }
                    break
                }
            }
            if (winner !== -1) {
                statistics.win[winner]++
            } else {
                statistics.peace++
            }
            console.log(`第${i}局的结果：${winner} round=${round}`)
            Game.show(game.a)
        }
        console.log(`经过${compareCount}次对局，对局结果为:\n${JSON.stringify(statistics, null, 2)}`)
    }
}
module.exports = comparer