const getAi = require('./ai.js')
const ai = getAi()
const Game = require('./game')

function one() {
    /**
     * 在全局都已经知道的情况下，AI能否走出最优解
     * */
    const a = [
        6, 13, 3, 1,
        15, 0, 7, 8,
        11, 12, 9, 14,
        5, 2, 10, 4]
    Game.show(a)
    console.log(ai.solve(a, new Set(), 0, 10))
}

function two() {
    /***
     * 红方能否抓住？
     * 狮子和豹：狮子+老鼠+狼
     * */
    //15去14还是去11
    a = [17, 17, 17, 17,
        17, 17, 17, 17,
        17, 1, 4, 17,
        9, 10, 17, 17]
    Game.show(a)
    console.log('judgeScore  ' + ai.judge(a, new Set()))
    ai.verbose.strategyScore = true
    console.log(ai.solve(a, new Set(), 1, 5))
}


function three() {
    a = [17, 17, 17, 17,
        2, 17, 10, 17,
        17, 17, 9, 17,
        17, 8, 17, 1
    ]
    ai.verbose.strategyScore = true
    // ai.verbose.leafInfo = true
    // ai.verbose.strategyScore=true
    ai.show(a)
    ai.MAX_DEPTH = 10
    const best = ai.solve(a, new Set(), 1)
    console.log(best)
}

function testNoSolution() {
    //无解，返回null
    const a = [17, 0, 17, 17,
        17, 17, 17, 17,
        17, 17, 1, 16,
        2, 17, 17, 3
    ]
    ai.MAX_DEPTH = 3
    const best = ai.solve(a, new Set(), 1)
    console.log(best)
}

function five() {
    //敌方有1象2狮，却在这里来回踱步
    const a = [
        3, 17, 17, 9,
        17, 17, 1, 17,
        8, 17, 17, 17,
        17, 17, 17, 4]
    Game.show(a)
    ai.MAX_DEPTH = 10
    ai.verbose.strategyScore = true
    const best = ai.solve(a, new Set(), 1)
    // ai.verbose.leafInfo = true
    // ai.verbose.strategyScore=true
    console.log(best)
    const b = [
        3, 17, 17, 9,
        17, 17, 1, 17,
        17, 8, 17, 17,
        17, 17, 17, 4]
    console.log(ai.solve(b, new Set(), 1))
}

/**
 * 如果递归中用到set，调用子函数时把元素从set中删除，调用完子函数之后把元素添加到set中
 * */
function testSet() {
    var a = new Set()
    a.add(1)
    a.add(2)
    for (var i of a) {
        a.delete(i)
        console.log(i)
        a.add(i)
    }
}


function six() {
    var a = [
        17, 17, 17, 17,
        17, 17, 0, 15,
        17, 17, 17, 10,
        17, 17, 17, 4,
    ]
    const ans = ai.solve(a, new Set(), 1)
    console.log(ans)
}

function seven() {
    const a = [
        11, 2, 0, 8,
        6, 17, 17, 5,
        10, 1, 17, 17,
        16, 7, 4, 3]
    const ans = ai.solve(a, new Set(), 0)
    console.log(ans)
}

function noLongChase() {
    const a = [
        1, 17, 17, 3,
        17, 10, 17, 2,
        17, 17, 0, 17,
        11, 17, 17, 9]
    const ans = ai.solve(a, new Set(), 0)
    console.log(ans)
}

function whyPeace() {
    const a = [[
        17, 17, 17, 17,
        17, 17, 17, 17,
        17, 2, 14, 9,
        17, 13, 17, 8
    ], [`---2
        --C-
        --B-
        ----`]]

    //这局棋本来应该大棋必胜，可是最终结果却是和棋
    const ans = ai.solve(a, new Set(), 1)
    console.log(ans)
}

function whyTigerNotEat() {
    const a = [
        6, 16, 16, 14,
        7, 16, 10, 17,
        16, 16, 17, 11,
        16, 12, 16, 3
    ]
    ai.verbose.strategyScore = true
    // ai.verbose.leafInfo = true
    ai.MAX_DEPTH = 10
    console.log(ai.solve(a, new Set([0, 1, 2, 4, 9, 8, 10]), 1))
}

function WhyElephantNew() {
    //敌人的大象为何要翻新牌而不是吃掉我的老虎
    const a = [7, 6, 16, 16,
        17, 16, 8, 16,
        1, 16, 3, 16,
        16, 0, 16, 4]
    const unknown = [2, 5, 9, 10, 11, 12, 13, 14]
    ai.verbose.strategyScore = true
    console.log(ai.solve(a, new Set(unknown), 1))
}

function whyNew() {
    //敌人的大象为何要翻新牌而不是吃掉我的老虎
    const a = [16, 16, 17, 16,
        0, 17, 17, 5,
        16, 2, 8, 16,
        1, 11, 16, 16]
    const unknown = [3, 6, 4, 7, 9, 10, 15]
    ai.verbose.strategyScore = true
    // ai.verbose.leafInfo = true
    ai.MAX_DEPTH = 10
    console.log(ai.solve(a, new Set(unknown), 0))
}

// three()
// testNoSolution()
// five()
// one()
// two()
// six()
// seven()
// 不许长追()
// whyTigerNotEat()
// WhyElephantNew()
whyNew()