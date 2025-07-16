package initialize

// func Redis() (err error) {
// 	cfg := global.ServerConfig.RedisInfo
// 	// 不可在返回error时，新定义变量 会出现空指针问题，因为已经定义了全局变量
// 	global.Rdb = redis.NewClient(&redis.Options{
// 		// Addr: fmt.Sprintf("%s:%d",
// 		// 	viper.GetString("redis.host"),
// 		// 	viper.GetInt("redis.port")),
// 		// Password: viper.GetString("redis.password"), // 密码
// 		// DB:       viper.GetInt("redis.db"),          // 数据库
// 		// PoolSize: viper.GetInt("redis.pool_size"),   // 连接池大小
// 		Addr:     strings.Join([]string{cfg.Host, ":", cfg.Port}, ""),
// 		Password: cfg.Password,
// 		DB:       cfg.DB,
// 		PoolSize: cfg.PoolSize,
// 	})

// 	_, err = global.Rdb.Ping(context.Background()).Result()
// 	if err != nil {
// 		zap.L().Error("connect redis failed, err:%v\n", zap.Error(err))
// 		return
// 	}
// 	return nil
// }

// func Close() {
// 	_ = global.Rdb.Close()
// }
