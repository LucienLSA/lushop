package main

// import (
// 	"context"
// 	"fmt"
// 	"goodssrv/proto"

// 	"google.golang.org/grpc"
// 	"google.golang.org/protobuf/types/known/emptypb"
// )

// var cClient proto.GoodsClient
// var conn *grpc.ClientConn

// func InitClient() {
// 	var err error
// 	// conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
// 	conn, err = grpc.Dial("127.0.0.1:8022", grpc.WithInsecure())
// 	if err != nil {
// 		panic(err)
// 	}
// 	// brandClient = proto.NewGoodsClient(conn)
// 	cClient = proto.NewGoodsClient(conn)
// }

// // 测试获取所有的商品分类列表
// func TestGetAllCategorysList() {
// 	rsp, err := cClient.GetAllCategorysList(context.Background(), &emptypb.Empty{})
// 	if err != nil {
// 		fmt.Println("查询失败")
// 		panic(err)
// 	}
// 	fmt.Println(rsp.JsonData)
// }

// // 测试获取品牌的列表
// func TestGetBrandList() {
// 	rsp, err := cClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
// 	if err != nil {
// 		fmt.Println("查询失败")
// 		panic(err)
// 	}
// 	fmt.Println(rsp.Total)
// 	for _, brand := range rsp.Data {
// 		fmt.Println(brand.Name)
// 	}
// }

// // 测试获取商品分类的子分类
// func TestGetSubCategory() {
// 	rsp, err := cClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
// 		Id: 136698,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(rsp.SubCategorys)
// }

// // 测试获取商品分类和品牌类别
// func TestCategoryBrandList() {
// 	rsp, err := cClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(rsp.Data)
// }

// // 测试通过商品分配获取品牌
// func TestGetCategoryBrandList() {
// 	rsp, err := cClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
// 		Id: 135200,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(rsp.Data)
// }

// // 测试商品查询
// func TestGoodsList() {
// 	rsp, err := cClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
// 		TopCategory: 136688,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(rsp.Total)
// 	fmt.Println(rsp.Data)
// }

// // 测试批量查询商品
// func TestBatchGetGoods() {
// 	rsp, err := cClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
// 		Id: []int32{421, 422, 423},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(rsp.Total)
// }

// // 测试获取商品的详情
// func TestGetGoodsDetail() {
// 	rsp, err := cClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
// 		Id: 421,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(rsp.Name)
// }

// func main() {
// 	InitClient()
// 	// TestGetBrandList()
// 	// TestGetAllCategorysList()
// 	// TestGetSubCategory()
// 	// TestGetCategoryBrandList()
// 	TestGoodsList()
// 	// TestBatchGetGoods()
// 	// TestGetGoodsDetail()
// 	conn.Close()
// }
