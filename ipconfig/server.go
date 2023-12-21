package ipconfig

import (
	"hello/plato_mjr/common/config"
	"hello/plato_mjr/ipconfig/domain"
	"hello/plato_mjr/ipconfig/source"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RunMain(path string) {
	config.InitConfig(path)
	source.Init()
	domain.Init()
	s := server.Default(server.WithHostPorts(":6777"))
	s.GET("ip/list", GetIpList)
	s.Spin()
}
