package global

import (
	"inventorysrv/config"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	Rdb          *redis.Client
	ServerConfig config.ServerConfig
	// NacosConfig  config.NacosConfig

	// MQ消费者
	MQPushClient rocketmq.PushConsumer // 消息消费者，用于订阅并消费消息
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
