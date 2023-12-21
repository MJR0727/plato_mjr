package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func InitConfig(fullFileName string) {
	viper.SetConfigFile(fullFileName)
	viper.SetConfigType("yaml")
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

func GetEndpointsDiscoveryPath() []string {
	return viper.GetStringSlice("discovery.endpoints")
}

func GtEndpointsTimeOut() time.Duration {
	return viper.GetDuration("discovery.timeout") * time.Second
}

func GetGlobleEnv() string {
	return viper.GetString("globle.env")
}

func GetServiceIpDispatcherPath() string {
	return viper.GetString("ipconfig.service_path")
}
