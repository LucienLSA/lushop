package test

import (
	"context"
	"lushopapi/user_web/global"
	"lushopapi/user_web/initialize"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestRedisGet(t *testing.T) {
	initialize.Redis()
	// 从redis中获取验证码
	value, err := global.Rdb.Get(context.Background(), "16666666").Result()
	if err == redis.Nil {
		t.Logf("key 不存在")
	}
	t.Logf(value)
}
