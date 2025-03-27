package initialize

import (
	"fmt"
	"lushopapi/user_web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Config() {
	mode := global.GetEnvInfoBool("LUSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s/%s-pro.yaml", configFilePrefix, configFilePrefix)
	if mode {
		configFileName = fmt.Sprintf("%s/%s-debug.yaml", configFilePrefix, configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 对象如何在其他文件中使用，因为它是一个局部变量，初始化全局变量
	// serverConfig := config.ServerConfig{}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Info("配置信息:%v", global.ServerConfig)
	// fmt.Printf("%v", v.Get("name"))
	// 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Info("配置文件产生变化:%s", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Info("配置信息:%v", global.ServerConfig)
	})
}
