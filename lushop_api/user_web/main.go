package main

import (
	"fmt"
	"userweb/global"
	"userweb/initialize"
	"userweb/utils/addr"
	"userweb/utils/register/consul"

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
	defer global.Rdb.Close()

	// 7. 初始化srv的连接
	initialize.SrvConn()
	zap.S().Info("init SrcConn success")

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
	err := register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败", err.Error())
	}
	zap.S().Info("init gin user web success")

	// 10. 优雅运行退出
	go func() {
		if err = Router.Run(fmt.Sprintf(":%v", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("用户web服务器启动失败", err.Error())
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
