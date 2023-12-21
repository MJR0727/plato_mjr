package discovery

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()
	s, err := NewServiceRegister("/http/node1", &EndpointInfo{
		IP:   "127.0.0.1",
		Port: "6666",
	}, &ctx, 5)
	if err != nil {
		log.Fatalln(err)
	}
	go s.ListenLeaseChan()
	select {
	case <-time.After(30 * time.Second):
		s.CancelLeaseClose()
	}
}
