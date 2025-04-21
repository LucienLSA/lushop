package main

import (
	"context"
	"fmt"
	"goodsweb/global"
	"goodsweb/initialize"
	"goodsweb/utils/addr"
	"goodsweb/utils/register/consul"
	"goodsweb/utils/track"
	"log"

	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func main() {
	// 1. 初始化zap日志
	logger := initialize.Logger()
	zap.S().Info("init Logger success")
	otelzap.L().Info("初始化logger成功")
	defer logger()

	// 2. 初始化配置文件
	initialize.Config()
	zap.S().Info("init Config success")
	otelzap.L().Info("初始化 Config 成功")

	// 10. 初始化trace
	ctx := context.Background()
	tp, err := track.InitTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("shutting down tracer provider failed, err:%v\n", err)
		}
	}()

	// 3. 初始化routers
	Router := initialize.Routers()
	zap.S().Info("init Routers success")

	// 4. 初始化validator翻译
	if err := initialize.Tran("zh"); err != nil {
		zap.S().Panic("初始化validator翻译失败", err.Error())
	}
	zap.S().Info("init validator trans success")

	// 5. 注册验证器
	initialize.SignUpMobile()
	zap.S().Info("init SignUpMobile success")

	// 7. 初始化srv的连接、包括nacos配置中心
	initialize.SrvConn()
	zap.S().Info("init SrcConn and nacos success")

	// 8. 初始化可用端口，debug模式则指定端口
	mode := global.GetEnvInfoBool(global.Mode)
	if !mode {
		port, err := addr.GetFreeport()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}
	zap.S().Info("init mode success")

	// 9. 初始化服务注册
	consulPortInt, _ := strconv.Atoi(global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.New())
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, consulPortInt)
	err = register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败", err.Error())
	}
	zap.S().Info("init gin service success")

	//  优雅运行退出
	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Errorf("启动商品web服务器失败,端口:%d", global.ServerConfig.Port)
		}
	}()
	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = register_client.DeRegister(serviceId)
	if err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}
}
