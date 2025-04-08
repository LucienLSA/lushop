package main

import (
	"context"
	proto_address "useropsrv/proto/gen/address"
	proto_message "useropsrv/proto/gen/message"
	proto_userfav "useropsrv/proto/gen/userfav"

	"google.golang.org/grpc"
)

var mClient proto_message.MessageClient
var aClient proto_address.AddressClient
var uClient proto_userfav.UserFavClient
var conn *grpc.ClientConn

func InitClient() {
	var err error
	// conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
	conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	aClient = proto_address.NewAddressClient(conn)
	mClient = proto_message.NewMessageClient(conn)
	uClient = proto_userfav.NewUserFavClient(conn)
}

func TestAddressList() {
	_, err := aClient.GetAddressList(context.Background(), &proto_address.AddressRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
}

func TestMessageList() {
	_, err := mClient.MessageList(context.Background(), &proto_message.MessageRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
}

func TestUserFavList() {
	_, err := uClient.GetFavList(context.Background(), &proto_userfav.UserFavRequest{
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
