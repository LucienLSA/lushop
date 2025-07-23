package global

import (
	"useropsrv/config"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	// Rdb          *redis.Client
	ServerConfig config.ServerConfig
	// NacosConfig  config.NacosConfig
)

const Mode = "LUSHOP_DEBUG"

// 以上未初始化对象，在init中依赖注入方式初始化

func GetEnvInfoBool(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func GetEnvInfoStr(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}
