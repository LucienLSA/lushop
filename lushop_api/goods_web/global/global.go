package global

import (
	"goodsweb/config"
	"goodsweb/proto"

	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator

	GoodsSrvClient proto.GoodsClient
	NacosConfig    *config.NacosConfig = &config.NacosConfig{}
)

var Tracer = otel.Tracer("goods_web")

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
