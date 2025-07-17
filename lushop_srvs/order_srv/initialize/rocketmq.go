package initialize

import (
	"fmt"
	"ordersrv/global"
	"ordersrv/handler"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
)

func InitMQ() {
	//test
	global.GroupInventory = producer.WithGroupName("mxshop-inventory")
	global.GroupOrder = producer.WithGroupName("mxshop-order")

	socket := fmt.Sprintf("%s:%d", global.ServerConfig.MqInfo.Host, global.ServerConfig.MqInfo.Port)
	fmt.Println(socket)
	global.MQInventory = InitMQNewProducer(global.GroupInventory, socket)
	global.MQOrder = InitMQNewProducer(global.GroupOrder, socket)
}

func InitMQNewProducer(Producer producer.Option, Socket string) rocketmq.Producer {
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

func InitRocketMQ() {
	// 生产者
	//发送事务消息 加入到事务组
	var err error
	//fmt.Println(global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)
	global.MQSendTranClient, err = rocketmq.NewTransactionProducer(
		&handler.OrderListener{},
		producer.WithNameServer([]string{fmt.Sprintf("%s:%s", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)}),
		producer.WithGroupName("transaction"),
	)
	if err != nil {
		fmt.Printf("【事务消息】生成producer失败: %s\n", err.Error())
		panic(err.Error())
		//return nil, err
		//fmt.Println("连接错误：", err)
	}
	//启动
	if err = global.MQSendTranClient.Start(); err != nil {
		fmt.Printf("【事务消息】启动producer失败: %s\n", err.Error())
		panic(err.Error())
		//return nil, err
	}
	// 生产者
	//发送延时消息 加入到延迟组
	global.MQSendClient, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{fmt.Sprintf("%s:%s", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)})),
		producer.WithGroupName("delay"),
	)
	if err != nil {
		fmt.Printf("【同步消息】生成producer失败: %s\n", err.Error())
		panic(err.Error())
	}
	//启动
	err = global.MQSendClient.Start()
	if err != nil {
		fmt.Printf("【同步消息】启动producer错误: %s\n", err.Error())
	}
	// 消费者
	//订阅消息 - 订单超时
	global.MQPushClient, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%s", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)}),
		consumer.WithGroupName(global.ServerConfig.RocketMQConfig.Group),
	)
	if err != nil {
		fmt.Printf("【订阅消息】生成producer失败: %s\n", err.Error())
		panic(err.Error())
	}
	if err = global.MQPushClient.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeout); err != nil {
		fmt.Printf("【订阅消息】失败：%s\n", err.Error())
		panic(err.Error())
	}
	//启动
	if err = global.MQPushClient.Start(); err != nil {
		fmt.Printf("【订阅消息】启动producer失败:%s\n", err.Error())
	}
}

// 注销消息生产者和消费者的客户端
func DeregisterMQ() {
	err := global.MQSendClient.Shutdown()
	if err != nil {
		zap.S().Error("【同步消息】注销失败")
	}
	err = global.MQSendTranClient.Shutdown()
	if err != nil {
		zap.S().Error("【事务消息】注销失败")
	}
	err = global.MQPushClient.Shutdown()
	if err != nil {
		zap.S().Error("【订阅消息】注销失败")
	}
}
