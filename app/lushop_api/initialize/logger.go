package initialize

import (
	"go.uber.org/zap"
)

// var (
// 	once sync.Once
// )

// func Logger() func() {
// 	// 创建 logger

// 	l, err := zap.NewDevelopment()
// 	if err != nil {
// 		panic(err)
// 	}
// 	logger := otelzap.New(
// 		l,
// 		// zap实例，按需配置
// 		otelzap.WithMinLevel(zap.InfoLevel), // 指定日志级别
// 		// otelzap.WithTraceIDField(true),      // 在日志中记录 traceID
// 	)
// 	// 替换全局的logger
// 	return otelzap.ReplaceGlobals(logger)
// }

func Logger() {
	logger, _ := zap.NewDevelopment()
	/*
		1. S()可以获取全局的sugar,可以自行设置一个全局的logger
		2. 日志分级别额，在生产环境中使用需要设置级别Debug, info,warn,error,fatal
		3. S()函数和L()函数集成了锁，提供了全局的安全访问的logger途径
	*/
	zap.ReplaceGlobals(logger)
}
