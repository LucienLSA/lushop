package initialize

// func RocketMQ() {
// 	// 消费者订阅库存归还的消息
// 	var err error
// 	global.MQPushClient, err = rocketmq.NewPushConsumer(
// 		consumer.WithNameServer([]string{fmt.Sprintf("%s:%s", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)}),
// 		consumer.WithGroupName(global.ServerConfig.RocketMQConfig.GroupName),
// 	)
// 	if err != nil {
// 		fmt.Printf("【订阅消息】生成producer失败: %s\n", err.Error())
// 		panic(err.Error())
// 	}
// 	if err = global.MQPushClient.Subscribe(global.ServerConfig.RocketMQConfig.TopicReback, consumer.MessageSelector{}, handler.AutoReback); err != nil {
// 		fmt.Printf("【订阅消息】失败：%s\n", err.Error())
// 		panic(err.Error())
// 	}
// 	//启动
// 	if err = global.MQPushClient.Start(); err != nil {
// 		fmt.Printf("【订阅消息】启动producer失败:%s\n", err.Error())
// 	}
// }

// // 注销消息生产者和消费者的客户端
// func DeregisterMQ() {
// 	err := global.MQPushClient.Shutdown()
// 	if err != nil {
// 		zap.S().Error("【订阅消息】注销失败")
// 	}
// }
