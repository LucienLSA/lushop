package initialize

import (
	"fmt"
	"userweb/global"
	"userweb/proto"

	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

func SrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(3),                           // 最大重试次数
		grpc_retry.WithPerRetryTimeout(1 * time.Second), // 每次超时最大时间
		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
	}
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%s/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		opts...,
	)
	if err != nil {
		zap.S().Fatalf("[Init SrvConn] 连接 [用户服务失败]", err.Error())
	}
	UserSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = UserSrvClient
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
