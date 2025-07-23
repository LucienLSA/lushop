package initialize

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

// var (
// 	once sync.Once
// )

func Logger() func() {
	// 创建 logger

	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger := otelzap.New(
		l,
		// zap实例，按需配置
		otelzap.WithMinLevel(zap.InfoLevel), // 指定日志级别
		// otelzap.WithTraceIDField(true),      // 在日志中记录 traceID
	)
	// 替换全局的logger
	return otelzap.ReplaceGlobals(logger)
}
