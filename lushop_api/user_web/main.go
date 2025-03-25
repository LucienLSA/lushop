package main

import (
	"fmt"
	"lushopapi/user_web/initialize"

	"go.uber.org/zap"
)

func main() {
	port := 8023
	// 1. 初始化zap日志
	initialize.Logger()
	// 2. 初始化routers
	Router := initialize.Routers()
	zap.S().Debugf("启动服务器，端口:%d", 8023)
	if err := Router.Run(fmt.Sprintf(":%v", port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
