package initialize

import (
	"fmt"
	"usersrv/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Config() {
	// 从配置文件中读取配置
	mode := global.GetEnvInfoBool(global.Mode)
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	fmt.Println(mode)
	if mode {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息:%v", global.ServerConfig)
}
