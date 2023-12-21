package domain

import (
	"hello/plato_mjr/ipconfig/source"
	"sort"
	"sync"
)

// 计算分数，派发合适的IP、管理Endpoint状态

type Dispatcher struct {
	Endpoints map[string]*Endpoint
	lock      sync.RWMutex
}

var dp *Dispatcher

func Init() {
	// 持续监听事件队列，更新状态。
	dp = &Dispatcher{
		Endpoints: make(map[string]*Endpoint),
	}
	eventChan := source.EventChan()
	go func() {
		for event := range eventChan {
			switch event.Type {
			case source.AddNodeEvent:
				dp.AddNode(event)
			case source.DelNodeEvent:
				dp.DelNode(event)
			}

		}
	}()
}

func Dispatch(ctx *IpconfigContext) []*Endpoint {
	// 1、获取候选的Endpoint
	eps := dp.getCalculateEndpoints()
	// 2、计算Endpoint的分值
	for _, ep := range eps {
		ep.CaculateSorce()
	}
	// 3、根据动静分数进行排序
	sort.Slice(eps, func(i, j int) bool {
		if eps[i].ActiveSorce > eps[j].ActiveSorce {
			return true
		} else if eps[i].ActiveSorce == eps[j].ActiveSorce {
			if eps[i].StaticSorce > eps[j].StaticSorce {
				return true
			}
			return false
		}
		return false
	})
	// 4、返回top5
	return top5Endports(eps)
}

func top5Endports(eps []*Endpoint) []*Endpoint {
	if len(eps) < 5 {
		return eps
	}
	return eps[:5]
}

func (dp *Dispatcher) getCalculateEndpoints() []*Endpoint {
	dp.lock.RLock()
	defer dp.lock.RUnlock()
	eps := make([]*Endpoint, 0, len(dp.Endpoints))
	for _, ep := range dp.Endpoints {
		eps = append(eps, ep)
	}
	return eps
}

// Add or update the endpoint at the map.
func (dp *Dispatcher) AddNode(event *source.Event) {

	// 1、Map索引找到对应的Endpoint（修改/新增）
	var ep *Endpoint
	ep.lock.Lock()
	defer ep.lock.Unlock()

	var key = event.Key()
	if ep = dp.Endpoints[key]; ep == nil {
		ep = NewEndpoint(event.IP, event.Port)
	}

	// 2、设置节点当前状态、状态窗口
	ep.UpdateStatus(&Status{
		MessageBytes:  event.MessageBytes,
		ConnectionNum: event.ConnectionNum,
	})
	// 3、将没有添加的节点，添加到Map中
	dp.lock.Lock()
	dp.Endpoints[key] = ep
	dp.lock.Unlock()
}

// Delete the endpoint at the map.
func (dp *Dispatcher) DelNode(event *source.Event) {
	dp.lock.Lock()
	defer dp.lock.Unlock()
	delete(dp.Endpoints, event.Key())
}
