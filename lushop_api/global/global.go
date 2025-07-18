package global

import (
	"lushopapi/config"
	v2goodsproto "lushopapi/proto/goods"
	v2inventoryproto "lushopapi/proto/inventory"
	v2orderproto "lushopapi/proto/order"
	v2userproto "lushopapi/proto/user"
	v2useropproto "lushopapi/proto/userop"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	Trans              ut.Translator
	ServerConfig       *config.ServerConfig = &config.ServerConfig{}
	UserSrvClient      v2userproto.UserClient
	UserOpSrvClient    v2useropproto.UserOpClient
	GoodsSrvClient     v2goodsproto.GoodsClient
	OrderSrvClient     v2orderproto.OrderClient
	InventorySrvClient v2inventoryproto.InventoryClient
	//NacosConfig   *config.NacosConfig = &config.NacosConfig{}
	RedisClient *redis.Client
	RedisSync   *redsync.Redsync
)

const Mode = "LUSHOP_DEBUG"

func GetEnvInfoBool(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func GetEnvInfoStr(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}
