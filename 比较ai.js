const com = require("./aiComparer")
const ai = require('./ai')
const ai1 = ai.getAi()
const ai2 = ai.getUrlAi("http://localhost:7788/solve")
ai1.MAX_DEPTH = 6
ai2.MAX_DEPTH = 10
com.compare(ai1, ai2, 10)