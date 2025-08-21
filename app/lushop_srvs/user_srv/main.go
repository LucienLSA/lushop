package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"usersrv/global"
	"usersrv/handler"
	"usersrv/initialize"
	proto "usersrv/proto"
	"usersrv/utils/addr"
	"usersrv/utils/register/consul"

	"net"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// 初始化Config
	initialize.Config()
	zap.S().Info("init Config sucess")
	// 初始化Mysql
	initialize.MySQL()
	zap.S().Info("init MySQL sucess")
	// 初始化日志
	initialize.Logger()
	zap.S().Info("init Logger sucess")
	zap.S().Info(global.ServerConfig)
	// IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", global.ServerConfig.Port, "端口号")
	flag.Parse()
	// zap.S().Info("ip:", *IP)
	if *Port == 0 {
		*Port, _ = addr.GetFreeport()
	}
	zap.S().Info("port:", *Port)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port))
	if err != nil {
		zap.S().Errorf("ip:", global.ServerConfig.Host)
		panic("failed to listen:" + err.Error())
	}
	// 注册grpc服务健康检查
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(server, healthcheck)

	// 服务注册
	// cfg := api.DefaultConfig()
	// // "192.168.226.140:8500"
	// cfg.Address = fmt.Sprintf("%s:%s", global.ServerConfig.ConsulInfo.Host,
	// 	global.ServerConfig.ConsulInfo.Port)
	// client, err := api.NewClient(cfg)
	// if err != nil {
	// 	panic(err)
	// }
	// // 生成对应的检查对象
	// check := &api.AgentServiceCheck{
	// 	// 后续从配置中心nacos中获取
	// 	GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
	// 	Timeout:                        "5s",
	// 	Interval:                       "5s",
	// 	DeregisterCriticalServiceAfter: "15s",
	// }
	// // 生成注册对象
	// registeration := new(api.AgentServiceRegistration)
	// registeration.Name = global.ServerConfig.Name
	// serviceID := fmt.Sprintf("%s", uuid.New())
	// registeration.ID = serviceID
	// registeration.Port = *Port
	// registeration.Tags = global.ServerConfig.Tags
	// registeration.Address = global.ServerConfig.Host
	// registeration.Check = check
	// err = client.Agent().ServiceRegister(registeration)
	// if err != nil {
	// 	panic(err)
	// }
	consulPortInt, _ := strconv.Atoi(global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.New())
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, consulPortInt)
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("【用户服务-srv】服务注册失败", err.Error())
	}
	zap.S().Debugf("【用户服务-srv】注册成功,port:%d", *Port)
	go func() {
		err = server.Serve(lis)
		if err != nil {
			zap.S().Errorf("failed to start grpc:" + err.Error())
		}
	}()

	// 接收终止信号，优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// err = client.Agent().ServiceDeregister(serviceID);
	err = register_client.DeRegister(serviceId)
	if err != nil {
		zap.S().Info("【用户服务-srv】服务注销失败")
	}
	zap.S().Info("【用户服务-srv】注销成功")
}
