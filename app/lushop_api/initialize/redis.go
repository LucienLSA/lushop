package initialize

import (
	"context"
	"lushopapi/global"

	"strings"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 初始化连接
func Redis() (err error) {
	cfg := global.ServerConfig.RedisInfo
	// 不可在返回error时，新定义变量 会出现空指针问题，因为已经定义了全局变量
	global.RedisClient = redis.NewClient(&redis.Options{
		// Addr: fmt.Sprintf("%s:%d",
		// 	viper.GetString("redis.host"),
		// 	viper.GetInt("redis.port")),
		// Password: viper.GetString("redis.password"), // 密码
		// DB:       viper.GetInt("redis.db"),          // 数据库
		// PoolSize: viper.GetInt("redis.pool_size"),   // 连接池大小
		Addr:     strings.Join([]string{cfg.Host, ":", cfg.Port}, ""),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = global.RedisClient.Ping(context.Background()).Result()
	if err != nil {
		zap.L().Error("connect redis failed, err:%v\n", zap.Error(err))
		return
	}

	// 初始化 RedisSync 分布式锁
	pool := goredis.NewPool(global.RedisClient)
	global.RedisSync = redsync.New(pool)

	return nil
}

func Close() {
	_ = global.RedisClient.Close()
}
