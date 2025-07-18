package initialize

import (
	"fmt"
	"lushopapi/global"

	v2goodsproto "lushopapi/proto/goods"
	v2inventoryproto "lushopapi/proto/inventory"
	v2orderproto "lushopapi/proto/order"
	v2userproto "lushopapi/proto/user"
	v2useropproto "lushopapi/proto/userop"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SrvConn() {
	InitUserSrv()   //用户srv
	InitUserOpSrv() //用户op srv
	InitGoodsSrv()  //商品srv
	InitOrderSrv()  //订单srv
	InitInvSrv()    //库存srv
}

func InitUserSrv() {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitUserSrv] 连接 【用户服务失败】")
		return
	}
	userSrcClient := v2userproto.NewUserClient(userConn)
	global.UserSrvClient = userSrcClient
}
func InitUserOpSrv() {
	consulInfo := global.ServerConfig.ConsulInfo
	userOpConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserOpSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitUserOpSrv] 连接 【用户op服务失败】")
		return
	}
	userOpSrcClient := v2useropproto.NewUserOpClient(userOpConn)
	global.UserOpSrvClient = userOpSrcClient
}
func InitGoodsSrv() {
	consulInfo := global.ServerConfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitGoodsSrv] 连接 【商品服务失败】")
		return
	}
	goodsSrcClient := v2goodsproto.NewGoodsClient(goodsConn)
	global.GoodsSrvClient = goodsSrcClient
}
func InitOrderSrv() {
	consulInfo := global.ServerConfig.ConsulInfo
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitOrderSrv] 连接 【订单服务失败】")
		return
	}
	orderSrcClient := v2orderproto.NewOrderClient(orderConn)
	global.OrderSrvClient = orderSrcClient
}
func InitInvSrv() {
	consulInfo := global.ServerConfig.ConsulInfo
	invConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitInvSrv] 连接 【库存服务失败】")
		return
	}
	invSrcClient := v2inventoryproto.NewInventoryClient(invConn)
	global.InventorySrvClient = invSrcClient
}

// func SrvConn() {
// 	consulInfo := global.ServerConfig.ConsulInfo
// 	var opts []grpc.DialOption
// 	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
// 		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
// 	)
// 	retryOpts := []grpc_retry.CallOption{
// 		grpc_retry.WithMax(3),                           // 最大重试次数
// 		grpc_retry.WithPerRetryTimeout(1 * time.Second), // 每次超时最大时间
// 		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
// 	}
// 	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
// 	goodsConn, err := grpc.Dial(
// 		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port,
// 			global.ServerConfig.GoodsSrvInfo.Name),
// 		opts...,
// 	)
// 	if err != nil {
// 		zap.S().Fatalf("[Init SrvConn] 连接 [商品服务失败]", err.Error())
// 	}
// 	GoodsSrvClient := proto.NewGoodsClient(goodsConn)
// 	global.GoodsSrvClient = GoodsSrvClient
// }

// func SrvConn2() {
// 	// 从注册中心获取到底层groc服务的信息
// 	cfg := api.DefaultConfig()
// 	consulInfo := global.ServerConfig.ConsulInfo
// 	cfg.Address = fmt.Sprintf("%s:%s", consulInfo.Host, consulInfo.Port)
// 	client, err := api.NewClient(cfg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`,
// 		global.ServerConfig.UserSrvInfo.Name))
// 	if err != nil {
// 		panic(err)
// 	}
// 	var userSrvHost string
// 	var userSrvPort int
// 	for _, value := range data {
// 		userSrvHost = value.Address
// 		userSrvPort = value.Port
// 		break
// 	}

// 	if userSrvHost == "" {
// 		zap.S().Fatal("[Init SrvConn] 连接 [用户服务失败]")
// 		return
// 	}

// 	// ip := global.ServerConfig.UserSrvInfo.Host
// 	// port := global.ServerConfig.UserSrvInfo.Port
// 	// 拨号连接grpc服务
// 	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
// 	if err != nil {
// 		zap.S().Errorw("[GetUserList] 连接【用户服务失败】",
// 			"msg", err.Error(),
// 		)
// 	}

// 	UserSrvClient := proto.NewUserClient(userConn)
// 	global.UserSrvClient = UserSrvClient
// }
