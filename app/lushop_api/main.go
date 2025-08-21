package main

import (
	"context"
	"fmt"
	"lushopapi/global"
	"lushopapi/initialize"
	"lushopapi/utils"
	"lushopapi/utils/register/consul"
	"lushopapi/utils/track"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	// 1. 初始化zap日志
	initialize.Logger()
	zap.S().Info("init Logger success")

	// 2. 初始化配置文件
	initialize.Config()
	zap.S().Info("init Config success")

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

	// 6. 初始化redis
	if err := initialize.Redis(); err != nil {
		zap.S().Panic("初始化redis失败", err.Error())
	}
	zap.S().Info("init redis success")
	defer global.RedisClient.Close()

	// 初始化定时任务
	initialize.InitCron()
	zap.S().Info("init cron success")

	// 初始化OAuth2认证
	initialize.OAuth2()
	zap.S().Info("init OAuth2 success")

	// 初始化session
	initialize.Session()
	zap.S().Info("init session success")

	//  初始化trace
	ctx := context.Background()
	tp, err := track.Tracer(ctx)
	if err != nil {
		zap.S().Panic("初始化Tracer失败", err.Error())
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			zap.S().Fatalf("shutting down tracer provider failed, err:%v\n", err)
		}
	}()

	// 初始化sentinel
	initialize.Sentinel()
	zap.S().Info("init sentinel success")

	// 7. 初始化srv的连接
	initialize.SrvConn()
	zap.S().Info("init SrcConn success")

	// 8. 初始化可用端口，debug模式则指定端口
	mode := global.GetEnvInfoBool(global.Mode)
	if !mode {
		port, err := utils.GetFreePort()
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

	zap.S().Debugf("启动v2lushopshop-api【http】服务器,端口：%d", global.ServerConfig.Port)

	// 10. 优雅运行退出
	go func() {
		if err = Router.Run(fmt.Sprintf(":%v", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("v2lushopshop-api【http】服务器启动失败", err.Error())
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
