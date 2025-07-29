package global

import (
	"lushopapi/config"
	v2goodsproto "lushopapi/proto/goods"
	v2inventoryproto "lushopapi/proto/inventory"
	v2orderproto "lushopapi/proto/order"
	v2userproto "lushopapi/proto/user"
	v2useropproto "lushopapi/proto/userop"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gopkg.in/boj/redistore.v1"
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
	// redis
	RedisClient *redis.Client
	RedisSync   *redsync.Redsync
	Redistore   *redistore.RediStore

	// OAuth2相关
	Srv *server.Server
	Mgr *manage.Manager
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
