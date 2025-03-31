package main

import (
	"context"
	"fmt"
	"lushopsrvs/goods_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var cClient proto.GoodsClient
var conn *grpc.ClientConn

func InitClient() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	// brandClient = proto.NewGoodsClient(conn)
	cClient = proto.NewGoodsClient(conn)
}

func TestGetAllCategorysList() {
	rsp, err := cClient.GetAllCategorysList(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("查询失败")
		panic(err)
	}
	fmt.Println(rsp.JsonData)
}

func TestGetBrandList() {
	rsp, err := cClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		fmt.Println("查询用户失败")
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

func main() {
	InitClient()
	// TestGetBrandList()
	TestGetAllCategorysList()
	conn.Close()
}
