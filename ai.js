if (typeof (importScripts) !== 'undefined') {
    importScripts('game.js')
    importScripts('https://cdn.bootcdn.net/ajax/libs/axios/0.21.1/axios.min.js')
}
if (typeof (require) !== 'undefined') {
    //global表示全局,把这个变量挂载到全局
    global.Game = require('./game.js')
    global.axios = require("axios")
}
if (!Game) {
    throw 'cannot import game.js'
}


function getAi() {
    return {
        //电脑自动走棋，电脑有三类决策：吃子，动子，翻牌
        INF: 10000,//局面评估函数的最大值
        MAX_DEPTH: 10,//ai搜索的最大深度
        STATIC_DEPTH: 2,//静态搜索步数
        /**
         * 猜测之后的搜索步数，猜测之后的最大深度不能超过5
         * 主要是防止翻完牌之后立马被吃掉，如果此值太大，会导致搜索变慢;如果此值太小，会导致计算机乱翻新棋
         * */
        GUESS_DEPTH: 3,
        bestStrategy: null,//搜索到的最佳决策
        visit: 0,//访问节点的个数
        verbose: {
            middleNode: false,
            strategyScore: true,
            leafInfo: false,
        },//是否冗余打印，用于调试

        judge(a, unkown) {
            //评价一个局面的好坏，0表示0-7,1代表8-15,只评价局面对0号玩家的好坏，如果0号玩家必胜，那么分数较高
            var chessWight = [1280, 640, 320, 160, 80, 40, 20, 10]
            var smallScore = 0, bigScore = 0
            for (const i of a) {
                if (i < 8) {
                    smallScore += chessWight[i % 8]
                } else if (i < 16) {
                    bigScore += chessWight[i % 8]
                }
            }

            for (const i of unkown) {
                if (i < 8) {
                    smallScore += chessWight[i % 8]
                } else {
                    bigScore += chessWight[i % 8]
                }
            }
            if (smallScore === 0) {
                //一个棋子都没了，必输
                return -this.INF
            } else if (bigScore === 0) {
                //对面一个棋子都没有了，必胜
                return this.INF
            }
            const score = smallScore - bigScore
            return score
        },
        solve(a, unknown, computerColor) {
            //AI的入口函数：unkown是Set，a是array
            this.bestStrategy = null
            this.visit = 0
            const unknownSet = new Set(unknown)
            //为了防止污染外部，此处必须复制一份，实际上这个工作应该在外部做，如果外部没有做，AI也应该自觉遵守规则
            const aa = a.slice()
            //校验一下unknown的正确性
            const unknownCount = aa.reduce((s, x) => s + (x === 16 ? 1 : 0), 0)
            if (unknownCount !== unknownSet.size) throw `invalid game state:${unknownCount}!=unknowSet.size(${unknownSet.size})`
            //因为inf还要加上depth，为了防止溢出alpha，beta，此处使用2×inf
            const score = this.go(aa, unknownSet, computerColor, -this.INF * 2, this.INF * 2, 0, this.MAX_DEPTH)
            //经过go之后，aa和unknownSet应该不变
            for (var i = 0; i < aa.length; i++) {
                if (aa[i] !== a[i]) throw 'aa changed after search'
            }
            if (unknown.length !== unknownSet.size) throw 'unknown set changed'
            const randInd = Math.floor(Math.random() * this.bestStrategy.length)
            const ans = this.bestStrategy[randInd]
            return {
                strategy: ans,
                score,
                bestStrategy: this.bestStrategy
            }
        },
        matai(score, depth) {
            /**
             * 使用深度作为马太效应的因子：如果score是正数，表明局面对它有利，那么深度越浅越好
             * 如果score是负数，表明局面对它不利，那么深度越深越好
             * */
            return score < 0 ? score + depth : score - depth
        },
        greedySorter(a, unknown, strategies) {
            //贪心法对着法进行排序
            return strategies.sort((x, y) => {
                if (x.length === y.length) {
                    if (x.length === 1) {
                        //两个都是翻新牌
                        return 0
                    } else if (x.length === 2) {
                        //两个都是吃子，谁吃的子大谁优先级高
                        //谁吃掉的东西重要让谁先走
                        let desDif = a[x[1]] - a[y[1]]
                        let srcDif = a[x[0]] - a[y[0]]
                        //如果同样是吃子，让子力强的子先吃子
                        //如果同样是走空步，让子力强的子先走空步
                        if (desDif === 0) return srcDif
                        else return desDif
                    } else {
                        throw 'error length'
                    }
                } else {
                    //优先吃子和动子，其次翻新牌
                    return -(x.length - y.length)
                }
            })
        },
        go(a, unknown, computerColor, alpha, beta, depth, maxDepth) {
            //a表示当前棋盘状态，computerColor表示电脑的花色（0表示0-7,1表示8-15）
            this.visit++
            if (depth >= maxDepth) {
                //如果达到了dfs深度
                let score = this.judge(a, unknown) * (computerColor === 0 ? 1 : -1)
                score = this.matai(score, depth)
                if (this.verbose.leafInfo)//如果冗余打印叶子节点的信息
                {
                    console.log(`depth=${depth} color=${computerColor} ${a} score=${score}`)
                    Game.show(a)
                }
                return score
            }
            //评判一下当前局面，如果当前局面已经是必胜或者必败状态，立即返回
            let currentScore = this.judge(a, unknown)
            if (Math.abs(currentScore) === this.INF) {
                //如果局面已经是终局状态
                if (computerColor === 1) {
                    currentScore *= -1
                }
                currentScore = this.matai(currentScore, depth)
                return currentScore
            }
            //生成着法
            var strategies = []
            for (var i = 0; i < 16; i++) {
                if (a[i] === 17) continue
                else if (a[i] === 16) {//翻牌
                    strategies.push([i])
                } else {
                    var chessColor = a[i] < 8 ? 0 : 1
                    if (chessColor === computerColor) {
                        //自己的颜色
                        var row = Math.floor(i / 4), col = i % 4
                        for (var d of Game.directions) {
                            var desRow = row + d[0], desCol = col + d[1]
                            if (!(desRow < 4 && desRow >= 0 && desCol < 4 && desCol >= 0)) continue
                            var des = desRow * 4 + desCol
                            if (a[des] === 16) continue//不能往未知元素上移动
                            else if (a[des] === 17) {
                                //向空白处移动棋子
                                strategies.push([i, des])
                            } else {
                                if (Game.chessCanEat(a[i], a[des])) {
                                    //如果颜色不同，那么可以吃
                                    strategies.push([i, des])
                                }
                            }
                        }
                    }
                }
            }
            //如果当前局面已经没有着法了，那么自己输了
            if (strategies.length === 0) return this.matai(-this.INF, depth)
            //对生成的着法进行排序，把胜率高的着法放在前面有利于快速剪枝
            strategies = this.greedySorter(a, unknown, strategies)
            //基于alpha-beta剪枝求最佳决策
            var score = -this.INF * 2//一开始对自己的评分非常低
            for (let stra of strategies) {
                var sonScore = null
                if (stra.length === 2) {
                    //如果是确定性的吃子操作
                    let [src, des] = stra
                    let eaten = a[des]
                    a[des] = a[src]
                    a[src] = 17
                    //如果是吃子着法而不是移动棋子，那么maxDepth应该保留足够大的空间。在关键的地方搜索的深度大一些
                    let nextMaxDepth = eaten === 17 ? maxDepth : Math.max(depth + this.STATIC_DEPTH, maxDepth)
                    sonScore = this.go(a, unknown, 1 - computerColor, -beta, -alpha, depth + 1, nextMaxDepth)
                    //恢复局面
                    a[src] = a[des]
                    a[des] = eaten
                } else {
                    //如果是翻新牌
                    //如果是不确定性的操作，那么枚举全部可能性，把未知的牌翻开
                    var allSonScore = 0
                    //注意：此处unkownList必须复制一份，否则导致下面因为频繁对unknown进行读写，永远遍历不完，产生死循环
                    var unknownList = Array.from(unknown)
                    for (var maybe of unknownList) {
                        a[stra[0]] = maybe
                        unknown.delete(maybe)
                        let nextMaxDepth = maxDepth
                        //此处很细腻
                        if (unknown.size > 0) {//如果未知元素太多，搜索不能太深
                            if (depth + this.GUESS_DEPTH < nextMaxDepth) {
                                nextMaxDepth = depth + this.GUESS_DEPTH
                            }
                        } else {
                            //如果已经没有未知元素了，那就正常搜索
                            nextMaxDepth = maxDepth + 1
                        }
                        const maybeScore = this.go(a, unknown, 1 - computerColor, -beta, -alpha, depth + 1, nextMaxDepth)
                        allSonScore += maybeScore
                        unknown.add(maybe)
                        a[stra[0]] = 16//此处应为16,重新遮住，而不应该是17。此处第一次写出了bug
                    }
                    //js中的Set没有length成员，只有size成员，如果使用length会导致sonScore=NaN
                    sonScore = allSonScore / unknown.size
                }
                if (sonScore == null) throw 'impossible to reach here! sonScore cannot be null'
                if (-sonScore >= score) {
                    if (depth === 0) {
                        //如果是第一层，那么记录着法
                        if (this.bestStrategy == null || -sonScore > score) {
                            this.bestStrategy = [stra]
                        } else {
                            this.bestStrategy.push(stra)
                        }
                    }
                    score = -sonScore
                }
                if (depth === 0 && this.verbose.strategyScore) {
                    console.log(`strategy:${stra}  score:${-sonScore}`)
                }
                alpha = Math.max(-sonScore, alpha)
                if (alpha > beta) {//一定不走这一枝，因为这一枝没有表现出足够的优势，父结点就不会选择这个儿子
                    break
                }
            }
            if (score === -this.INF * 2) {
                throw 'impossible'
            }
            if (this.verbose.middleNode) {
                console.log(depth + ' ' + a + ' ' + score)
            }
            return score
        }
    }
}

if (typeof (module) != 'undefined') {
    //如果是内部测试
    module.exports = getAi
} else {
    //this表示子线程中的全局
    this.addEventListener('message', function (e) {
        const {a, unknown, maxDepth, computerColor, id, aiUrl} = e.data
        if (aiUrl) {
            //如果提供了AiUrl，则发起HTTP请求获取结果
            axios.get(aiUrl, {
                params: {a: a.join(','), id, computerColor, unknown: unknown.join(',')}
            }).then(resp => {
                console.log(`got message from ${aiUrl} data=${resp.data}`)
                this.postMessage(resp.data);
            })
        } else {
            if (!(a && unknown && maxDepth && (computerColor === 0 || computerColor === 1))) {
                throw 'invalid arg'
            }
            const ai = getAi()
            ai.MAX_DEPTH = maxDepth
            const ans = ai.solve(a, unknown, computerColor)
            ans.id = id//把id放回去
            this.postMessage(ans)
        }
    }, false);
}
