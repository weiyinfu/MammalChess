const com = require("./aiComparer")
const getAi = require('./ai')
const ai = getAi()
const ai2 = getAi()
ai.MAX_DEPTH = 10
ai2.MAX_DEPTH = 10
com.compare(ai, ai2, 10)