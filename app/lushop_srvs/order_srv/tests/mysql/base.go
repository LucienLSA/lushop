package main

import (
	"context"
	"fmt"
	proto_order "ordersrv/proto/gen/order"

	"google.golang.org/grpc"
)

var oClient proto_order.OrderClient
var conn *grpc.ClientConn

func InitClient() {
	var err error
	// conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
	conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	oClient = proto_order.NewOrderClient(conn)
}

func TestCreateCartItem(userId, nums, goodsId int32) {
	rsp, err := oClient.CreateCartItem(context.Background(), &proto_order.CartItemRequest{
		UserId:  userId,
		Nums:    nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func TestCartItemList(userId int32) {
	rsp, err := oClient.CartItemList(context.Background(), &proto_order.UserInfo{
		Id: userId,
	})
	if err != nil {
		panic(err)
	}
	for _, item := range rsp.Data {
		fmt.Println(item.Id, item.GoodsId, item.Nums)
	}
}

func TestUpdateCartItem(id int32) {
	_, err := oClient.UpdateCartItem(context.Background(), &proto_order.CartItemRequest{
		Id:      id,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateOrder() {
	_, err := oClient.CreateOrder(context.Background(), &proto_order.OrderRequest{
		UserId:  1,
		Address: "北京市",
		Name:    "bobby",
		Mobile:  "18787878787",
		Post:    "请尽快发货",
	})
	if err != nil {
		panic(err)
	}
}

func TestOrderDetail(orderId int32) {
	rsp, err := oClient.OrderDetail(context.Background(), &proto_order.OrderRequest{
		Id: orderId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.OrderInfo.OrderSn)
	for _, good := range rsp.Goods {
		fmt.Println(good.GoodsName)
	}
}

func TestOrderList() {
	rsp, err := oClient.OrderList(context.Background(), &proto_order.OrderFilterRequest{})
	if err != nil {
		panic(err)
	}
	for _, order := range rsp.Data {
		fmt.Println(order.OrderSn)
	}
}

func main() {
	InitClient()
	fmt.Println("init success")
	// TestCreateCartItem(1, 1, 422)
	// TestCartItemList(1)
	// TestUpdateCartItem(1)
	// TestCreateOrder()
	// TestOrderDetail(1)
	TestOrderList()
	conn.Close()
}
