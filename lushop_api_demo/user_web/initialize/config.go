package initialize

import (
	"fmt"
	"userweb/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Config() {
	mode := global.GetEnvInfoBool(global.Mode)
	// mode := global.GetEnvInfoBool("LUSHOP_DEBUG")
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
