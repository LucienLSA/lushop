package main

import (
	"fmt"
	"lushopapi/user_web/global"
	"lushopapi/user_web/initialize"
	"lushopapi/user_web/utils/addr"

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
	mode := global.GetEnvInfoBool(global.ServerConfig.Mode)
	if !mode {
		port, err := addr.GetFreeport()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}
	zap.S().Info("init mode success")

	if err := Router.Run(fmt.Sprintf(":%v", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("用户web服务器启动失败", err.Error())
	}
	zap.S().Debugf("启动用户web服务器,端口:%d", global.ServerConfig.Port)
}
