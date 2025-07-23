package global

import (
	"useropweb/config"
	proto_address "useropweb/proto/gen/address"
	proto_goods "useropweb/proto/gen/goods"
	proto_message "useropweb/proto/gen/message"
	proto_userfav "useropweb/proto/gen/userfav"

	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"
)

var (
	ServerConfig   *config.ServerConfig = &config.ServerConfig{}
	Trans          ut.Translator
	GoodsSrvClient proto_goods.GoodsClient

	AddressSrvClient proto_address.AddressClient
	MessageSrvClient proto_message.MessageClient
	UserFavSrvClient proto_userfav.UserFavClient

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
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
