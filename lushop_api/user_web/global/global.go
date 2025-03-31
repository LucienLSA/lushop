package global

import (
	"lushopapi/user_web/config"
	"lushopapi/user_web/proto"

	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator
	// var 声明全局的rdb变量
	Rdb           *redis.Client
	UserSrvClient proto.UserClient
	NacosConfig   *config.NacosConfig = &config.NacosConfig{}
)

// var (
// 	ServerConfig config.ServerConfig
// 	Trans        ut.Translator
// 	// var 声明全局的rdb变量
// 	Rdb           *redis.Client
// 	UserSrvClient proto.UserClient
// 	NacosConfig   config.NacosConfig
// )

const Mode = "LUSHOP_DEBUG"

// var ServerConfig = new(config.ServerConfig)

func GetEnvInfoBool(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func GetEnvInfoStr(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}
