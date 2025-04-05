package main

import (
	"fmt"
	"lushopsrvs/order_srv/proto"

	"google.golang.org/grpc"
)

var vClient proto.OrderClient
var conn *grpc.ClientConn

func InitClient() {
	var err error
	// conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
	conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	vClient = proto.NewOrderClient(conn)
}

func main() {
	InitClient()
	fmt.Println("init success")

	conn.Close()
}
