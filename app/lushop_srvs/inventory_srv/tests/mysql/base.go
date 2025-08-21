package main

// import (
// 	"context"
// 	"fmt"
// 	"inventorysrv/proto"

// 	"sync"

// 	"google.golang.org/grpc"
// )

// var vClient proto.InventoryClient
// var conn *grpc.ClientConn

// func InitClient() {
// 	var err error
// 	// conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
// 	conn, err = grpc.Dial("127.0.0.1:56709", grpc.WithInsecure())
// 	if err != nil {
// 		panic(err)
// 	}
// 	vClient = proto.NewInventoryClient(conn)
// }

// func TestSetInv(goodsId, num int32) {
// 	_, err := vClient.SetInv(context.Background(), &proto.GoodsInvInfo{
// 		GoodsId: goodsId,
// 		Num:     num,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("设置库存成功")
// }

// func TestInvDetail(goodsId int32) {
// 	resp, err := vClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
// 		GoodsId: goodsId,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(resp.Num)
// }

// func TestSell(wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	_, err := vClient.Sell(context.Background(), &proto.SellInfo{
// 		GoodsInfo: []*proto.GoodsInvInfo{
// 			{
// 				GoodsId: 421,
// 				Num:     1,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("商品库存扣减成功")
// }

// func TestReback() {
// 	_, err := vClient.Reback(context.Background(), &proto.SellInfo{
// 		GoodsInfo: []*proto.GoodsInvInfo{
// 			{
// 				GoodsId: 421,
// 				Num:     10,
// 			},
// 			{
// 				GoodsId: 423,
// 				Num:     10,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("商品库存归还成功")
// }

// func main() {
// 	InitClient()
// 	fmt.Println("init success")
// 	var i int32
// 	for i = 421; i < 840; i++ {
// 		TestSetInv(i, 100)
// 	}
// 	// TestInvDetail(421)
// 	// TestSell()
// 	// TestReback()
// 	// var wg sync.WaitGroup
// 	// wg.Add(10)
// 	// for i := 0; i < 10; i++ {
// 	// 	go TestSell(&wg)
// 	// }
// 	// wg.Wait()
// 	conn.Close()
// }
