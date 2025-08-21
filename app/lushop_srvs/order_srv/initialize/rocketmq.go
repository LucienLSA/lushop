package initialize

import (
	"fmt"
	"ordersrv/global"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
)

// 初始化RocketMQ
func RocketMQ() {
	// 初始化生产者组
	global.GroupInventory = producer.WithGroupName(global.ServerConfig.RocketMQInfo.ProducerGroupInventory)
	global.GroupOrder = producer.WithGroupName(global.ServerConfig.RocketMQInfo.ProducerGroupOrder)

	// 初始化生产者
	socket := fmt.Sprintf("%s:%s", global.ServerConfig.RocketMQInfo.Host, global.ServerConfig.RocketMQInfo.Port)
	global.MQInventory = RocketMQNewProducer(global.GroupInventory, socket)
	global.MQOrder = RocketMQNewProducer(global.GroupOrder, socket)
}

func RocketMQNewProducer(Producer producer.Option, Socket string) rocketmq.Producer {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{Socket}),
		Producer)
	if err != nil {
		zap.S().Errorf("初始化生产者失败%s", err)
	}
	if err = p.Start(); err != nil {
		zap.S().Errorf("启动生产者失败%s", err)
	}
	return p
}

// func RocketMQ() {
// 	// 生产者
// 	//发送事务消息 加入到事务组
// 	var err error
// 	//fmt.Println(global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)
// 	global.MQSendTranClient, err = rocketmq.NewTransactionProducer(
// 		&handler.OrderListener{},
// 		producer.WithNameServer([]string{fmt.Sprintf("%s:%s", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)}),
// 		producer.WithGroupName(global.ServerConfig.RocketMQConfig.ProducerGroupTran),
// 	)
// 	if err != nil {
// 		zap.S().Errorf("【事务消息】生成producer失败: %s\n", err.Error())
// 		panic(err.Error())
// 		//return nil, err
// 		//fmt.Println("连接错误：", err)
// 	}
// 	//启动
// 	if err = global.MQSendTranClient.Start(); err != nil {
// 		zap.S().Errorf("【事务消息】启动producer失败: %s\n", err.Error())
// 		panic(err.Error())
// 		//return nil, err
// 	}
// 	zap.S().Debug("【事务消息】启动producer成功")
// 	// 生产者
// 	//发送延时消息 加入到延迟组
// 	global.MQSendClient, err = rocketmq.NewProducer(
// 		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{fmt.Sprintf("%s:%s", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)})),
// 		producer.WithGroupName(global.ServerConfig.RocketMQConfig.ProducerGroupDelay),
// 	)
// 	if err != nil {
// 		zap.S().Errorf("【延时消息】生成producer失败: %s\n", err.Error())
// 		panic(err.Error())
// 	}
// 	//启动
// 	err = global.MQSendClient.Start()
// 	if err != nil {
// 		zap.S().Errorf("【延时消息】启动producer错误: %s\n", err.Error())
// 	}
// 	zap.S().Debug("【延时消息】启动producer成功")

// 	// 消费者
// 	//订阅消息 - 订单超时
// 	global.MQPushClient, err = rocketmq.NewPushConsumer(
// 		consumer.WithNameServer([]string{fmt.Sprintf("%s:%s", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)}),
// 		consumer.WithGroupName(global.ServerConfig.RocketMQConfig.ConsumerGroup),
// 	)
// 	if err != nil {
// 		zap.S().Errorf("【订单超时】生成consumer失败: %s\n", err.Error())
// 		panic(err.Error())
// 	}
// 	if err = global.MQPushClient.Subscribe(global.ServerConfig.RocketMQConfig.TopicTimeOut, consumer.MessageSelector{}, handler.OrderTimeout); err != nil {
// 		zap.S().Errorf("【订单超时】生成consumer失败: %s\n", err.Error())
// 		panic(err.Error())
// 	}
// 	//启动
// 	if err = global.MQPushClient.Start(); err != nil {
// 		zap.S().Errorf("【订单超时】生成consumer失败: %s\n", err.Error())
// 	}
// 	zap.S().Debug("【订单超时】生成consumer成功")
// }

// // 注销消息生产者和消费者的客户端
// func DeregisterMQ() {
// 	err := global.MQSendClient.Shutdown()
// 	if err != nil {
// 		zap.S().Error("【延迟消息】注销失败")
// 	}
// 	err = global.MQSendTranClient.Shutdown()
// 	if err != nil {
// 		zap.S().Error("【事务消息】注销失败")
// 	}
// 	err = global.MQPushClient.Shutdown()
// 	if err != nil {
// 		zap.S().Error("【订单超时订阅消息】注销失败")
// 	}
// }
