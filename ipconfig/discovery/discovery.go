package discovery

import (
	"log"

	"github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/clientv3"
)

type ServiceDiscovery struct {
	cli *clientv3.Client
}

func NewServiceDiscovery() {
	// 创建 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.getEndpointsDiscoveryPath()},
		DialTimeout: config.getEndpointsTimeOut,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
}
