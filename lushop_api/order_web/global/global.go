package global

import (
	"orderweb/config"
	proto_goods "orderweb/proto/gen/goods"
	proto_inventory "orderweb/proto/gen/inventory"
	proto_order "orderweb/proto/gen/order"

	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"
)

var (
	ServerConfig       *config.ServerConfig = &config.ServerConfig{}
	Trans              ut.Translator
	GoodsSrvClient     proto_goods.GoodsClient
	OrderSrvClient     proto_order.OrderClient
	InventorySrvClient proto_inventory.InventoryClient

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
