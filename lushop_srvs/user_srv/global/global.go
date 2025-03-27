package global

import (
	"lushopsrvs/user_srv/config"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
)

// 创建数据库客户端上下文
// 用于将Go的context.Context对象与数据库操作绑定在一起，主要作用是对控制请求的超时、取消、追踪等。
// func NewDBClient(ctx context.Context) *gorm.DB {
// 	db := _db
// 	return db.WithContext(ctx)
// }
