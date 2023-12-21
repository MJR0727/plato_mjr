package domain

const windowSize = 5

type StatWindow struct {
	statQueue []*Status
	statChan  chan *Status
	sumStat   *Status
	idx       int64
}

func newStatWindow() *StatWindow {
	return &StatWindow{
		statQueue: make([]*Status, windowSize),
		statChan:  make(chan *Status),
		sumStat:   &Status{},
	}
}

func (sw *StatWindow) GetStatus() *Status {
	return sw.sumStat.Avg(windowSize)
}

func (sw *StatWindow) appendStat(stat *Status) {
	// 1、将最久没有使用的stat淘汰
	sw.sumStat.Sub(sw.statQueue[sw.idx%windowSize])
	sw.statQueue[sw.idx%windowSize] = stat
	// 2、更新stat总和
	sw.sumStat.Add(stat)
	// 3、更新指针下标
	sw.idx++
}
