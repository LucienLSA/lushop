package main

import (
	"fmt"
	"lushopapi/user_web/global"
	"lushopapi/user_web/initialize"

	"go.uber.org/zap"
)

func main() {
	// 1. 初始化zap日志
	initialize.Logger()
	// 2. 初始化配置文件
	initialize.Config()
	// 3. 初始化routers
	Router := initialize.Routers()
	// 4. 初始化validator翻译
	if err := initialize.Tran("zh"); err != nil {
		zap.S().Panic("初始化validator翻译失败", err.Error())
	}
	// 5. 注册验证器
	initialize.SignUpMobile()

	zap.S().Debugf("启动用户web服务器,端口:%d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%v", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("用户web服务器启动失败", err.Error())
	}
}
