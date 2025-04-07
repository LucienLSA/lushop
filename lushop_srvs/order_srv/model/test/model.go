package main

import (
	"log"
	"ordersrv/model"

	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// func NewDBClient(ctx context.Context) *gorm.DB {
// 	db := _db
// 	return db.WithContext(ctx)
// }

var (
	_db *gorm.DB
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/lushop_order_srv?charset=utf8mb4&parseTime=True&loc=Local"
	ormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  false, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    false, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   false, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: true,  // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		// To log
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)  // 设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) // 打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	_db = db
	_db = _db.Set("gorm:table_options", "charset=utf8mb4")
	err = _db.AutoMigrate(&model.OrderGoods{}, &model.OrderInfo{}, &model.ShoppingCart{})
	if err != nil {
		// Todo log
		panic(err)
	}
}
