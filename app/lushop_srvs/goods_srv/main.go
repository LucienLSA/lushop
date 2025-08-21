package main

import (
	"flag"
	"fmt"
	"goodssrv/global"
	"goodssrv/handler"
	"goodssrv/initialize"
	proto "goodssrv/proto"
	"goodssrv/utils/addr"
	"net/http"
	"strconv"

	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	zap.S().Info("init Config sucess")
	// 初始化Mysql
	initialize.MySQL()
	zap.S().Info("init MySQL sucess")
	// 初始化es
	initialize.Es()
	zap.S().Info("init ES sucess")

	zap.S().Info(global.ServerConfig)
	// IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", global.ServerConfig.Port, "端口号")
	flag.Parse()
	// zap.S().Info("ip:", *IP)
	if *Port == 0 {
		*Port, _ = addr.GetFreeport()
	}
	// zap.S().Info("port:", *Port)
	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	// 注册grpc服务健康检查
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(server, healthcheck)

	// 服务注册
	cfg := api.DefaultConfig()
	// "192.168.226.140:8500"
	cfg.Address = fmt.Sprintf("%s:%s", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Panic("【商品服务-srv】注册失败")
		panic(err)
	} else {
		zap.S().Info("ip:", global.ServerConfig.Host, ":", *Port)
		zap.S().Info("【商品服务-srv】注册成功")
	}
	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		// 后续从配置中心nacos中获取
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	// 生成注册对象
	registeration := new(api.AgentServiceRegistration)
	registeration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.New())
	registeration.ID = serviceID
	registeration.Port = *Port
	// 这里修改为配置中心所定义的Tags
	// registeration.Tags = []string{"goods_srv", "lushop_srv", "grpc", "lucien"}
	registeration.Tags = global.ServerConfig.Tags
	registeration.Address = global.ServerConfig.Host
	// 添加metrics元数据
	registeration.Meta = map[string]string{
		"metrics":      "true",
		"metrics_port": strconv.Itoa(global.ServerConfig.MetricsPort),
		"metrics_path": "/metrics",
	}
	registeration.Check = check
	err = client.Agent().ServiceRegister(registeration)
	if err != nil {
		panic(err)
	}

	// 启动服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			zap.S().Errorf("failed to start grpc:" + err.Error())
		}
	}()

	// 启动Prometheus metrics服务
	go func() {
		metricsAddr := fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.MetricsPort)
		http.Handle("/metrics", promhttp.Handler())
		zap.S().Infof("Prometheus metrics server started on %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			zap.S().Errorf("Failed to start metrics server: %v", err)
		}
	}()

	// 接收终止信号，优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Panic("【商品服务-srv】注销失败")
	} else {
		zap.S().Info("【商品服务-srv】注销成功")
	}
}
