//游戏规则相关的描述
/**
 * 游戏状态编码：
 * 0：0～7胜利
 * 1：8～15胜利
 * 2：和棋
 * 3：还不确定
 * */
//游戏状态
const Game = {
    FIRST_WIN: 0,
    SECOND_WIN: 1,
    PEACE: 2,
    UNSURE: 3,
    //四个方向：动物可以移动的四个方向
    directions: [[0, 1], [1, 0], [0, -1], [-1, 0]],
    chessCanEat(srcValue, desValue) {
        //判断棋子srcValue是否可以吃掉desValue
        const srcColor = srcValue < 8, desColor = desValue < 8
        if (srcColor === desColor) return false
        srcValue %= 8
        desValue %= 8
        if (srcValue === 7 && desValue === 0) return true//老鼠吃大象
        if (srcValue === 0 && desValue === 7) return false//大象不能吃老鼠
        if (srcValue <= desValue) return true
        return false
    },
    show(a) {
        //向控制台打印棋盘
        for (let i = 0; i < 4; i++) {
            let line = ''
            for (let j = 0; j < 4; j++) {
                let chess = a[i * 4 + j]
                if (chess === 16) chess = '*'
                else if (chess === 17) chess = '-'
                else {
                    chess = chess >= 8 ? 'ABCDEFGH'[chess % 8] : (chess + 1)
                }
                line += chess
            }
            console.log(line)
        }
        console.log()
    }
}

Game.newGame = () => {
    const game = {
        a: [],//棋盘，使用一维包含16个元素的数组表示
        unknown: [],//还没有翻出来的牌
        PEACE_CLOCK: 25,//十步之内没有吃子判定和棋
        peaceClock: 0,//当前的peaceClock

        init() {
            //新游戏
            //获取游戏的初始状态
            var a = []
            var unknown = []
            for (let i = 0; i < 16; i++) {
                a.push(16)
                unknown.push(i)
            }
            this.a = a
            this.unknown = unknown
        },
        isNeibor(src, des) {
            //src和des是否相邻
            var [srcX, srcY] = [Math.floor(src / 4), src % 4]
            var [desX, desY] = [Math.floor(des / 4), des % 4]
            if (Math.abs(srcX - desX) + Math.abs(srcY - desY) === 1) {
                return true
            }
            return false
        },

        canEat(src, des) {
            //src处的棋子是否可以移动到des处
            if (!(this.a[src] >= 0 && this.a[src] < 16)) return false
            if (this.a[des] === 17) return true
            if (this.a[des] === 16) return false
            if (!this.isNeibor(src, des)) return false
            if (Game.chessCanEat(this.a[src], this.a[des])) return true
            return false
        },
        eat(src, des) {
            //src处的棋子吃掉des处的棋子
            const temp = this.a[des]
            this.a[des] = this.a[src]
            this.a[src] = 17
            if (temp !== 17) {//如果目标处有棋子，那么就要吃掉它
                this.peaceClock = 0//发生吃子，peace钟清零
            } else {
                //目标处没有棋子，走得是空步
                this.peaceClock++
            }
        },
        canRecover(ind) {
            return this.a[ind] === 16
        },
        recover(ind) {
            //翻开ind处的棋子
            const unknownInd = Math.floor(Math.random() * this.unknown.length)
            const chessId = this.unknown[unknownInd]
            this.unknown.splice(unknownInd, 1)
            this.a[ind] = chessId
            this.peaceClock = 0
        },
        haveStrategy(color) {
            //判断一方是否有着法，当一方没有着法时，该方为输
            for (var i = 0; i < this.a.length; i++) {
                if (this.a[i] >= 16) continue//空格或者未知
                const chessColor = this.a[i] < 8 ? 0 : 1
                if (chessColor !== color) continue
                for (var d of Game.directions) {
                    const x = Math.floor(i / 4) + d[0], y = i % 4 + d[1]
                    if (x >= 0 && x < 4 && y >= 0 && y < 4) {
                        if (this.canEat(i, x * 4 + y)) {
                            return true
                        }
                    }
                }
            }
            return false
        },
        getState() {
            //判断游戏是否结束，当一方棋子被全部吃光时，游戏结束
            //当一方棋子无可行着法时，游戏结束
            var count0 = 0, count1 = 0
            for (let i = 0; i < this.a.length; i++) {
                if (this.a[i] < 8) {
                    count0++
                } else if (this.a[i] < 16) {
                    count1++
                }
            }
            for (const i of this.unknown) {
                if (i < 8) count0++
                else if (i < 16) count1++
            }
            /**
             * 先判断胜负，再判断和棋，顺序不能错。
             * */
            if (count0 === 0) return Game.SECOND_WIN
            if (count1 === 0) return Game.FIRST_WIN
            if (this.peaceClock >= this.PEACE_CLOCK) {
                return Game.PEACE
            }
            if (this.unknown.length === 0) {
                //如果unknown为空，那么判断双方棋子是否可移动，无法移动的一方为输
                if (!this.haveStrategy(0)) {
                    return Game.SECOND_WIN
                }
                if (!this.haveStrategy(1)) {
                    return Game.FIRST_WIN
                }
            }
            return Game.UNSURE
        },
    }
    game.init()
    return game
}

if (typeof (module) != 'undefined') {
    //如果是内部测试
    module.exports = Game
}