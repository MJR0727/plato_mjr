package source

import (
	"context"
	"fmt"
	"hello/plato_mjr/common/config"
	"hello/plato_mjr/common/discovery"
)

// 监听etcd数据
func Init() {
	ctx := context.Background()
	go DataHandler(ctx)
	// 注册mock数据进行测试
	if config.GetGlobleEnv() == "debug" {
		testServiceRegister(&ctx, "6767", "node1")
		testServiceRegister(&ctx, "6868", "node2")
		testServiceRegister(&ctx, "6869", "node3")
	}
}

func DataHandler(ctx context.Context) {
	serviceDiscovery := discovery.NewServiceDiscovery(&ctx)
	defer serviceDiscovery.Close()
	set := func(key, value string) {
		if epi, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(epi); epi != nil {
				event.Type = AddNodeEvent
				// 放入队列消费
				eventChan <- event
			}
		} else {
			fmt.Sprintf("DataHandler error.Delete EndPoint{IP:%s;Port:%s} fail.\n", epi.IP, epi.Port)
		}
	}
	del := func(key, value string) {
		if epi, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(epi); epi != nil {
				event.Type = DelNodeEvent
				// 放入队列消费
				eventChan <- event
			}
		} else {
			fmt.Sprintf("DataHandler error.Delete EndPoint{IP:%s;Port:%s} fail.\n", epi.IP, epi.Port)
		}
	}
	if err := serviceDiscovery.WatcherService(config.GetServiceIpDispatcherPath(), set, del); err != nil {
		panic(err)
	}
}
