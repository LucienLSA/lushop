package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"useropsrv/global"
	"useropsrv/handler"

	"useropsrv/initialize"
	proto "useropsrv/proto"

	"useropsrv/utils/addr"
	"useropsrv/utils/register/consul"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// 初始化日志
	initialize.Logger()
	zap.S().Info("init Logger sucess")
	// 初始化Config
	initialize.Config()
	zap.S().Info("init config sucess")
	// 初始化Mysql
	initialize.MySQL()
	zap.S().Info("init MySQL sucess")
	// // 初始化redis
	// if err := initialize.Redis(); err != nil {
	// 	zap.S().Panic("初始化redis失败", err.Error())
	// }
	// zap.S().Info("init Redis sucess")

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
	// proto_message.RegisterMessageServer(server, &handler.UserOpServer{})
	// proto_address.RegisterAddressServer(server, &handler.UserOpServer{})
	// proto_userfav.RegisterUserFavServer(server, &handler.UserOpServer{})
	proto.RegisterUserOpServer(server, &handler.UserOpServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port))
	if err != nil {
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
	// // 这里修改为配置中心所定义的Tags
	// // registeration.Tags = []string{"goods_srv", "lushop_srv", "grpc", "lucien"}
	// registeration.Tags = global.ServerConfig.Tags
	// registeration.Address = global.ServerConfig.Host
	// registeration.Check = check
	// err = client.Agent().ServiceRegister(registeration)
	// if err != nil {
	// 	panic(err)
	// }

	// 9. 初始化服务注册
	consulPortInt, _ := strconv.Atoi(global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.New())
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, consulPortInt)
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("【用户操作服务-srv】注册失败:", err.Error())
	} else {
		zap.S().Info("ip:", global.ServerConfig.Host, ":", *Port)
		zap.S().Info("【用户操作服务-srv】注册成功")
	}

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
	err = register_client.DeRegister(serviceId)
	// client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		zap.S().Panic("【用户操作服务-srv】注销失败:", err.Error())
	} else {
		zap.S().Info("【用户操作服务-srv】注销成功")
	}
}
