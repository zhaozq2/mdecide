Page({
  data: {
    currentPhase: 0,
    currentRoundIndex: 0,
    newOption: '',
    roundList: []
  },

  onLoad() {
    this.initRound()
  },

  get currentRound() {
    return this.data.roundList[this.data.currentRoundIndex]
  },

  initRound() {
    const newRound = {
      options: [],
      importanceScores: [],
      necessityScores: [],
      results: [],
      summary: '',
      importanceStatus: 'pending',
      necessityStatus: 'pending',
      status: 'pending'
    }
    this.setData({ roundList: [newRound] })
  },

  onInput(e) {
    this.setData({ newOption: e.detail.value })
  },

  addOption() {
    if (this.data.newOption.trim()) {
      const round = this.currentRound
      round.options.push(this.data.newOption.trim())
      round.importanceScores.push(0)
      round.necessityScores.push(0)
      this.setData({
        [`roundList[${this.data.currentRoundIndex}]`]: round,
        newOption: ''
      })
    }
  },

  removeOption(e) {
    const index = e.currentTarget.dataset.index
    const round = this.currentRound
    round.options.splice(index, 1)
    round.importanceScores.splice(index, 1)
    round.necessityScores.splice(index, 1)
    this.setData({
      [`roundList[${this.data.currentRoundIndex}]`]: round
    })
  },

  startImportance() {
    if (this.currentRound.options.length < 2) {
      wx.showToast({ title: '至少需要2个选项', icon: 'none' })
      return
    }
    this.setData({ currentPhase: 1 })
  },

  finishImportance() {
    const round = this.currentRound
    if (round.importanceScores.some(s => s === 0)) {
      wx.showToast({ title: '请完成所有评分', icon: 'none' })
      return
    }
    round.importanceStatus = 'completed'
    this.setData({
      [`roundList[${this.data.currentRoundIndex}]`]: round,
      currentPhase: 2
    })
  },

  setImportance(e) {
    const index = e.currentTarget.dataset.index
    const score = e.currentTarget.dataset.score
    const round = this.currentRound
    round.importanceScores[index] = score
    this.setData({
      [`roundList[${this.data.currentRoundIndex}]`]: round
    })
  },

  setNecessity(e) {
    const index = e.currentTarget.dataset.index
    const value = parseInt(e.currentTarget.dataset.value)
    const round = this.currentRound
    round.necessityScores[index] = value
    this.setData({
      [`roundList[${this.data.currentRoundIndex}]`]: round
    })
  },

  calculateResults() {
    const round = this.currentRound
    const resultsData = round.options.map((opt, idx) => {
      const importance = round.importanceScores[idx]
      const necessity = round.necessityScores[idx]
      return {
        option: opt,
        importanceScore: importance,
        necessityScore: necessity,
        totalScore: importance * necessity
      }
    })

    resultsData.sort((a, b) => b.totalScore - a.totalScore)

    round.results = resultsData.map((r, idx) => ({
      ...r,
      rank: idx + 1,
      isWinner: idx === 0 && r.totalScore > 0
    }))

    const winner = round.results.find(r => r.isWinner)
    if (winner) {
      round.summary = `经过重要性评分和必要性筛选，${winner.option}综合得分最高(${winner.totalScore}分)，是最佳选择。`
    } else {
      round.summary = '所有选项综合得分相同，请重新评估。'
    }

    round.necessityStatus = 'completed'
    this.setData({
      [`roundList[${this.data.currentRoundIndex}]`]: round,
      currentPhase: 3
    })
  },

  finishRound() {
    const round = this.currentRound
    round.status = 'completed'
    this.setData({
      [`roundList[${this.data.currentRoundIndex}]`]: round
    })
  },

  addRound() {
    const roundList = this.data.roundList
    const firstRound = roundList[0]
    const newRound = {
      options: [...firstRound.options],
      importanceScores: new Array(firstRound.options.length).fill(0),
      necessityScores: new Array(firstRound.options.length).fill(0),
      results: [],
      summary: '',
      importanceStatus: 'pending',
      necessityStatus: 'pending',
      status: 'pending'
    }
    roundList.push(newRound)
    this.setData({
      roundList,
      currentRoundIndex: roundList.length - 1,
      currentPhase: 0
    })
  },

  switchRound(e) {
    const index = e.currentTarget.dataset.index
    const round = this.data.roundList[index]
    let phase = 0
    if (round.status === 'completed') {
      phase = 3
    } else if (round.importanceStatus === 'completed') {
      phase = 2
    }
    this.setData({
      currentRoundIndex: index,
      currentPhase: phase
    })
  }
})