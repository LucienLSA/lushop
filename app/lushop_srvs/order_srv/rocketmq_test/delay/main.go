package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.226.140:9876"}))
	if err != nil {
		panic(err)
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	msg := primitive.NewMessage("lucien1", []byte("this is a delayed message"))
	msg.WithDelayTimeLevel(2)
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Printf("发送失败:%s\n", err)
	} else {
		fmt.Printf("发送成功:%s\n", res.String())
	}
	if err = p.Shutdown(); err != nil {
		panic(err)
	}
}
