package initialize

import (
	"fmt"
	"useropweb/global"
	proto_address "useropweb/proto/gen/address"
	proto_goods "useropweb/proto/gen/goods"
	proto_message "useropweb/proto/gen/message"
	proto_userfav "useropweb/proto/gen/userfav"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	userOpConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserOpSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[Init SrvConn] 连接 [用户操作服务失败]", err.Error())
	}
	global.UserFavSrvClient = proto_userfav.NewUserFavClient(userOpConn)
	global.MessageSrvClient = proto_message.NewMessageClient(userOpConn)
	global.AddressSrvClient = proto_address.NewAddressClient(userOpConn)

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
}

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
