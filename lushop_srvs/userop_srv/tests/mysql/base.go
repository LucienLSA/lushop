package main

import (
	"context"

	proto "useropsrv/proto"

	"google.golang.org/grpc"
)

// var mClient proto_message.MessageClient
// var aClient proto_address.AddressClient
// var uClient proto_userfav.UserFavClient
var client proto.UserOpClient
var conn *grpc.ClientConn

func InitClient() {
	var err error
	// conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
	conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	// aClient = proto_address.NewAddressClient(conn)
	// mClient = proto_message.NewMessageClient(conn)
	// uClient = proto_userfav.NewUserFavClient(conn)
	client = proto.NewUserOpClient(conn)
}

func TestAddressList() {
	_, err := client.GetAddressList(context.Background(), &proto.AddressRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
}

func TestMessageList() {
	_, err := client.MessageList(context.Background(), &proto.MessageRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
}

func TestUserFavList() {
	_, err := client.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	InitClient()
	TestAddressList()
	TestMessageList()
	TestUserFavList()
	conn.Close()
}
