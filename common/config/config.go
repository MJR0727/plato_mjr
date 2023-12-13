package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func initConfig(fullFileName string) {
	viper.SetConfigFile(fullFileName)
	viper.SetConfigType("yaml")
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

func getEndpointsDiscoveryPath() string {
	return viper.GetString("discovery.endpoints")
}

func getEndpointsTimeOut() string {
	return viper.GetString("discovery.timeout")
}

func getGlobleEnv() string {
	return viper.GetString("globle.env")
}

func getServiceIpDispatcherPath() string {
	return viper.GetString("ip_config.service_path")
}
