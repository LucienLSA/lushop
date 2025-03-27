package global

import (
	"context"
	"lushopapi/user_web/config"

	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator
	// var 声明全局的rdb变量
	Rdb  *redis.Client
	Rctx = context.Background()
)

// var ServerConfig = new(config.ServerConfig)
