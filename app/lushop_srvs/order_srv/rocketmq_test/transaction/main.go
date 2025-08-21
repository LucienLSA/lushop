package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type OrderListener struct{}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地逻辑")
	time.Sleep(time.Second * 3)
	fmt.Println("执行本地逻辑失败")
	return primitive.UnknowState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("rocketmq的消息回查")
	time.Sleep(time.Second * 15)
	return primitive.CommitMessageState
}

func main() {
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{},
		producer.WithNameServer([]string{"192.168.226.140:9876"}),
	)
	if err != nil {
		panic(err)
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	res, err := p.SendMessageInTransaction(context.Background(),
		primitive.NewMessage("TransTopic", []byte("This is a transation topic")))
	if err != nil {
		fmt.Printf("发送失败:%s\n", err)
	} else {
		fmt.Printf("发送成功:%s\n", res.String())
	}
	time.Sleep(time.Hour)
	if err = p.Shutdown(); err != nil {
		panic(err)
	}
}
