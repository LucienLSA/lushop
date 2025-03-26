package main

import (
	"flag"
	"fmt"

	"lushopsrvs/user_srv/handler"
	"lushopsrvs/user_srv/initialize"
	"lushopsrvs/user_srv/proto"
	"net"

	"google.golang.org/grpc"
)

func main() {
	initialize.MySQL()
	IP := flag.String("ip", "127.0.0.1", "ip地址")
	Port := flag.Int("port", 8022, "端口号")
	flag.Parse()
	fmt.Println("ip:", *IP)
	fmt.Println("port:", *Port)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
