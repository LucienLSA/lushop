package global

import (
	"usersrv/config"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)

const Mode = "LUSHOP_DEBUG"

// 创建数据库客户端上下文
// 用于将Go的context.Context对象与数据库操作绑定在一起，主要作用是对控制请求的超时、取消、追踪等。
// func NewDBClient(ctx context.Context) *gorm.DB {
// 	db := _db
// 	return db.WithContext(ctx)
// }

func GetEnvInfoBool(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func GetEnvInfoStr(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}
