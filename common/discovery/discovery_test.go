package discovery

import (
	"context"
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	ctx := context.Background()
	s := NewServiceDiscovery(&ctx)
	defer s.Close()
	s.WatcherService("/gRpc", func(key, value string) {}, func(key, value string) {})
	s.WatcherService("/http/", func(key, value string) {}, func(key, value string) {})
	for {
		select {
		case <-time.Tick(10 * time.Second):
		}
	}
}
