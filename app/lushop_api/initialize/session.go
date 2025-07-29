package initialize

import (
	"encoding/gob"
	"lushopapi/global"
	"net/url"
	"strconv"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"gopkg.in/boj/redistore.v1"
)

// Setup 初始化Redis会话存储
func Session() {
	// 注册可能存储在会话中的类型
	// gob.Register(map[string]interface{}{})
	gob.Register(url.Values{})
	var err error
	// 使用 NewRediStoreWithDB 初始化 Redis 存储
	global.Redistore, err = redistore.NewRediStoreWithDB(
		global.ServerConfig.RedisInfo.PoolSize,
		"tcp",
		global.ServerConfig.RedisInfo.Host+":"+global.ServerConfig.RedisInfo.Port,
		global.ServerConfig.RedisInfo.Password,
		strconv.Itoa(global.ServerConfig.RedisInfo.DB),
		[]byte(global.ServerConfig.SessionInfo.SecretKey),
	)
	if err != nil {
		zap.S().Fatalf("Redis会话存储初始化失败: %v", err)
	}

	// 配置会话选项
	global.Redistore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   global.ServerConfig.SessionInfo.MaxAge, // 会话有效期
		HttpOnly: true,                                   // 防止XSS攻击
		Secure:   false,                                  // HTTPS时设为true
	}
}
