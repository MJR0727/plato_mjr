package discovery

import (
	"context"
	"hello/plato_mjr/common/config"
	"log"

	"github.com/bytedance/gopkg/util/logger"
	"go.etcd.io/etcd/clientv3"
)

type ServiceRegister struct {
	cli           *clientv3.Client
	ctx           *context.Context
	leaseId       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	value         string
}

func NewServiceRegister(key string, value *EndpointInfo, ctx *context.Context, lease int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndpointsDiscoveryPath(),
		DialTimeout: config.GtEndpointsTimeOut(),
	})
	if err != nil {
		log.Fatal(err)
	}
	ser := &ServiceRegister{
		cli:   cli,
		ctx:   ctx,
		key:   key,
		value: value.Marshal(),
	}
	if err = ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}
	return ser, nil
}

// 设置租期和续租
func (ser *ServiceRegister) putKeyWithLease(lease int64) error {
	resp, err := ser.cli.Grant(*ser.ctx, lease)
	if err != nil {
		return err
	}
	// 带有租期的put
	_, err = ser.cli.Put(*ser.ctx, ser.key, ser.value, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	// 续租
	leaseKeepAliveChan, err := ser.cli.KeepAlive(*ser.ctx, resp.ID)
	if err != nil {
		return err
	}
	ser.leaseId = resp.ID
	ser.keepAliveChan = leaseKeepAliveChan
	return nil
}

func (ser *ServiceRegister) UpdateValue(epi *EndpointInfo) error {
	value := epi.Marshal()
	// 带有租期的put
	_, err := ser.cli.Put(*ser.ctx, ser.key, value, clientv3.WithLease(ser.leaseId))
	if err != nil {
		return err
	}
	ser.value = value
	logger.CtxInfof(*ser.ctx, "ServiceRegister.updateValue leaseID=%d Put key=%s,val=%s, success!", ser.leaseId, ser.key, ser.value)
	return nil
}

func (ser *ServiceRegister) ListenLeaseChan() {
	for leaseKeepAliveRsp := range ser.keepAliveChan {
		logger.CtxInfof(*ser.ctx, "lease success leaseID:%d, Put key:%s,val:%s reps:+%v",
			ser.leaseId, ser.key, ser.value, leaseKeepAliveRsp)
	}
	logger.CtxErrorf(*ser.ctx, "lease keepAlive fail! LeaseID:%d, Put key:%s,val:%s",
		ser.leaseId, ser.key, ser.value)
}

func (ser *ServiceRegister) CancelLeaseClose() error {
	if _, err := ser.cli.Revoke(*ser.ctx, ser.leaseId); err != nil {
		return err
	}
	logger.CtxInfof(*ser.ctx, "lease close !!!  leaseID:%d, Put key:%s,val:%s  success!", ser.leaseId, ser.key, ser.value)
	return ser.cli.Close()
}
