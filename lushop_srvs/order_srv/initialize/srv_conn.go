package initialize

import (
	"fmt"
	"ordersrv/global"
	proto_goods "ordersrv/proto/gen/goods"
	proto_inventory "ordersrv/proto/gen/inventory"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	// 商品服务连接
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[Init SrvConn] 连接 [商品服务失败]", err.Error())
	}
	GoodsSrvClient := proto_goods.NewGoodsClient(goodsConn)
	global.GoodsSrvClient = GoodsSrvClient
	// 库存服务连接
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[Init SrvConn] 连接 [库存服务失败]", err.Error())
	}
	global.InventorySrvClient = proto_inventory.NewInventoryClient(inventoryConn)
}
