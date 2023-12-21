package domain

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Endpoint struct {
	IP           string
	Port         string
	ActiveSorce  float64
	StaticSorce  float64
	status       *Status
	statusWindow *StatWindow
	lock         sync.RWMutex
}

func NewEndpoint(ip, port string) *Endpoint {
	ep := &Endpoint{
		IP:   ip,
		Port: port,
	}
	ep.statusWindow = newStatWindow()
	ep.status = ep.statusWindow.GetStatus()
	go func() {
		// 持续监听状态队列，更新ep最新状态
		for stat := range ep.statusWindow.statChan {
			ep.statusWindow.appendStat(stat)
			newStat := ep.statusWindow.GetStatus()
			// 原子交换新老状态值
			atomic.SwapPointer((*unsafe.Pointer)((unsafe.Pointer)(ep.status)), unsafe.Pointer(newStat))
		}
	}()
	// 分值计算在dispatcher里维护
	return ep
}

func (ep *Endpoint) UpdateStatus(stat *Status) {
	ep.statusWindow.statChan <- stat
}

func (ep *Endpoint) CaculateSorce() {
	if ep.status != nil {
		ep.ActiveSorce = ep.status.CalculateActiveSorce()
		ep.StaticSorce = ep.status.CalculateStaticSorce()
	}
}
