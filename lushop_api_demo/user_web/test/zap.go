package test

import (
	"time"

	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./lushop.log",
		"stderr",
	}
	return cfg.Build()
}

func main() {
	// logger, _ := zap.NewProduction()
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	defer logger.Sync() // flushes buffer, if any

	url := "www.bilibili.com"
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}
