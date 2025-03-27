package initialize

import (
	"fmt"
	"lushopapi/user_web/global"
	"lushopapi/user_web/proto"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func SrcConn() {
	// 从注册中心获取到底层groc服务的信息
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%s", consulInfo.Host, consulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`,
		global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	var userSrvHost string
	var userSrvPort int
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	if userSrvHost == "" {
		zap.S().Fatal("[Init SrvConn] 连接 [用户服务失败]")
		return
	}

	// ip := global.ServerConfig.UserSrvInfo.Host
	// port := global.ServerConfig.UserSrvInfo.Port
	// 拨号连接grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接【用户服务失败】",
			"msg", err.Error(),
		)
	}

	UserSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = UserSrvClient
}
