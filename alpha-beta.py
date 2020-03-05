import numpy as np

inf = 1e9

"""
alpha-beta剪枝我自己想不出来，我只能想出beta剪枝
"""


def randomTree(depth):
    # 生成随机树
    if depth == 0:
        return {'sons': [], 'value': np.random.randint(1, 10)}
    # sons_count = np.random.randint(2, 4)
    sons_count = 2
    root = {'sons': [], 'value': 0}
    for i in range(sons_count):
        root['sons'].append(randomTree(depth - 1))
    return root


def Evaluate(node):
    return node['value']


def expand(node):
    return node['sons']


negmax_alphabeta_visit = 0


def NegaMax_AlphaBeta(node, alpha, beta, color):
    global negmax_alphabeta_visit
    negmax_alphabeta_visit += 1
    if len(node['sons']) == 0:
        print('leaf', node)
        return Evaluate(node)
    score = -inf
    for childNode in expand(node):
        value = -NegaMax_AlphaBeta(childNode, -beta, -alpha, -color)
        score = max(score, value)
        alpha = max(alpha, value)
        if alpha >= beta:
            break
    return score


bruteforce_visit = 0


def bruteForce(node):
    global bruteforce_visit
    bruteforce_visit += 1
    if not node['sons']: return node['value']
    me = inf
    for son in node['sons']:
        me = min(bruteForce(son), me)
    return -me


mine_visit = 0


def mine(node, fatherScore):
    # 此算法可以称为alpha剪枝，它没有beta剪枝
    global mine_visit
    mine_visit += 1
    if not node['sons']:
        print('leaf', node)
        return node['value']
    minSon = inf
    for son in expand(node):
        sonScore = mine(son, -minSon)
        minSon = min(minSon, sonScore)
        if minSon <= fatherScore:
            return inf
    return -minSon


def printTree(node, prefix):
    if not node['sons']:
        print(prefix + str(node['value']))
        return
    print(prefix + '$')
    for son in node['sons']:
        printTree(son, prefix + '  ')


def testMany():
    # 正确性测试
    for i in range(10):
        testOne()


def testOne():
    r = np.random.randint(0, 100)
    # r = 8
    print('seed', r)
    np.random.seed(r)
    # tree = randomTree(4)
    tree = randomTree(np.random.randint(2, 6))
    # printTree(tree, '')
    score = NegaMax_AlphaBeta(tree, alpha=-inf, beta=inf, color=1)
    print('knuth的alpha-beta剪枝算法', score, negmax_alphabeta_visit)
    print('暴力算法', bruteForce(tree), bruteforce_visit)
    mineScore = mine(tree, fatherScore=-inf)
    print('我的剪枝算法', mineScore, mine_visit)


def testBad():
    global mine_visit, negmax_alphabeta_visit
    badTree = {
        'sons': [
            {'sons': [], 'value': -4},
            {'sons': [{
                'sons': [
                    {'sons': [
                        {'sons': [], 'value': 6},
                        {'sons': [], 'value': 2}
                    ], 'value': 0},
                    {'sons': [
                        {'sons': [], 'value': 3.5},
                        {'sons': [], 'value': 3}
                    ], 'value': 0},
                ], 'value': 0
            }], 'value': 0
            }
        ], 'value': 0
    }
    print('=========')
    printTree(badTree, '')
    mine_visit = 0
    score = mine(badTree, -inf)
    print('mine visit score:', mine_visit, score)
    negmax_alphabeta_visit = 0
    score = NegaMax_AlphaBeta(badTree, -inf, inf, 1)
    print('his visit score', negmax_alphabeta_visit, score)


# testOne()
# testMany()
testBad()
