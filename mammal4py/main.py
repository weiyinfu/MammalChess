import json
from typing import List

import flask
from flask import request, make_response
import random
import logging

"""
这是一个随机走棋的DEMO，用于演示如何实现一个HTTP服务版的AI。
"""
app = flask.Flask(__name__)
logging.root.setLevel(logging.INFO)


def legal(x, y):
    return 4 > x >= 0 and 4 > y >= 0


def get_color(chess):
    return 0 if chess < 8 else 1


def can_eat(animal: int, food: int):
    if animal == 7:
        if food == 0 or food == 7:
            return True
        else:
            return False
    elif animal == 0:
        if food == 7:
            return False
        else:
            return True
    else:
        return animal <= food


def get_move(a: List[int], pos: int):
    # pos处的棋子的全部着法
    ans = []
    for dx, dy in ((0, 1), (0, -1), (-1, 0), (1, 0)):
        x, y = pos // 4, pos % 4
        tx, ty = x + dx, y + dy
        to = tx * 4 + ty
        if not legal(tx, ty):
            continue
        if a[to] == 16:
            continue
        if a[to] == 17:
            # 空白
            ans.append(to)
        if get_color(a[to]) == get_color(a[pos]):
            # 同色不能相吃
            continue
        if can_eat(a[pos] % 8, a[to] % 8):
            ans.append(to)
    return [[pos, i] for i in ans]


def go(a: List[int], unknown: List[int], computer_color: int):
    ans = []
    for pos, i in enumerate(a):
        if i == 16:
            # unknown
            ans.append([pos])
        elif i == 17:
            continue
        elif i < 8 and computer_color == 0:
            # try move it
            ans.extend(get_move(a, pos))
        elif i >= 8 and computer_color == 1:
            ans.extend(get_move(a, pos))
    return ans


@app.route("/solve")
def solve():
    computer_color = int(request.args['computerColor'])
    board = request.args['a']
    board = list(map(int, board.split(',')))
    unknown = request.args['unknown']
    unknown = [int(i) for i in unknown]
    req_id = request.args['id']
    logging.info(f"computer_color={computer_color} board={board} id={req_id}")
    moves = go(board, unknown, computer_color)
    logging.info(f"moves={moves}")
    resp = make_response(json.dumps({'strategy': random.choice(moves)}))
    resp.headers['Access-Control-Allow-Origin'] = '*'
    return resp


if __name__ == '__main__':
    app.run(debug=True, port=6677)
