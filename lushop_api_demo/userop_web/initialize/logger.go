package initialize

import "go.uber.org/zap"

func Logger() {
	logger, _ := zap.NewDevelopment()
	/*
		1. S()可以获取全局的sugar,可以自行设置一个全局的logger
		2. 日志分级别额，在生产环境中使用需要设置级别Debug, info,warn,error,fatal
		3. S()函数和L()函数集成了锁，提供了全局的安全访问的logger途径
	*/
	zap.ReplaceGlobals(logger)
}
