package discovery

import (
	"context"
	"fmt"
	"log"
	"sync"

	"hello/plato_mjr/common/config"

	"go.etcd.io/etcd/clientv3"
)

type ServiceDiscovery struct {
	cli  *clientv3.Client
	ctx  *context.Context
	lock sync.Mutex
}

func NewServiceDiscovery(ctx *context.Context) *ServiceDiscovery {
	// 创建 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndpointsDiscoveryPath(),
		DialTimeout: config.GtEndpointsTimeOut(),
	})
	if err != nil {
		log.Fatal(err)
	}
	return &ServiceDiscovery{
		cli: cli,
		ctx: ctx,
	}
}

func (s *ServiceDiscovery) WatcherService(prefix string, set, del func(key, value string)) error {
	rsp, err := s.cli.Get(*s.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	// 初次触发set事件函数
	for _, kv := range rsp.Kvs {
		set(string(kv.Key), string(kv.Value))
	}
	// 持续监听往后的版本
	s.watch(prefix, rsp.Header.Revision+1, set, del)
	return nil
}

func (s *ServiceDiscovery) watch(prefix string, vision int64, set, del func(key, value string)) {
	watchChan := s.cli.Watch(*s.ctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(vision))
	for resp := range watchChan {
		for _, event := range resp.Events {
			switch event.Type {
			case clientv3.EventTypePut:
				set(string(event.Kv.Key), string(event.Kv.Value))
				fmt.Printf("Service registered: %s\n", event.Kv.Key)
			case clientv3.EventTypeDelete:
				del(string(event.Kv.Key), string(event.Kv.Value))
				fmt.Printf("Service unregistered: %s\n", event.Kv.Key)
			}
		}
	}
}

func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}
