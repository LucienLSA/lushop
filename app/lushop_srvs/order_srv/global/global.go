package global

import (
	"ordersrv/config"
	v2goods "ordersrv/proto/goods"
	v2inventory "ordersrv/proto/inventory"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	Rdb          *redis.Client
	ServerConfig config.ServerConfig
	// NacosConfig  config.NacosConfig

	GoodsSrvClient     v2goods.GoodsClient
	InventorySrvClient v2inventory.InventoryClient

	// // MQ消费者
	// MQPushClient rocketmq.PushConsumer // 消息消费者，用于订阅并消费消息
	// // MQ生产者
	// MQSendTranClient rocketmq.TransactionProducer // 事务消息生产者，用于发送分布式事务消息。

	GroupInventory producer.Option
	GroupOrder     producer.Option

	MQOrder     rocketmq.Producer // 普通消息生产者，用于发送普通消息到消息队列。
	MQInventory rocketmq.Producer // 普通消息生产者，用于发送普通消息到消息队列。
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
