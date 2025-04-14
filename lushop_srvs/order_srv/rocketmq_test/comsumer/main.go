package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.226.140:9876"}),
		consumer.WithGroupName("lushop"),
	)
	if err != nil {
		panic(err)
	}
	err = c.Subscribe("lucien1", consumer.MessageSelector{},
		func(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range me {
				fmt.Printf("获取到值: %v \n", me[i])
			}
			return consumer.ConsumeSuccess, nil
		})
	if err != nil {
		fmt.Println("读取消息失败")
	}
	err = c.Start()
	if err != nil {
		panic(err)
	}
	// 不能让主协程退出
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		panic(err)
	}
}
