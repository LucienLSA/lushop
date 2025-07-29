package initialize

import (
	"fmt"
	"log"
	"ordersrv/global"
	"ordersrv/model"

	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

// 定义局部变量进行数据库初始化，最后再传给全局变量
// 理想的情况应该为如下，后续进行重构后改进
// 创建数据库客户端上下文
// 用于将Go的context.Context对象与数据库操作绑定在一起，主要作用是对控制请求的超时、取消、追踪等。
// func NewDBClient(ctx context.Context) *gorm.DB {
// 	db := _db
// 	return db.WithContext(ctx)
// }

// var (
// 	_db *gorm.DB
// )

func MySQL() {
	mysqlCfg := global.ServerConfig.MySQLInfo
	pathRead := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlCfg.User, mysqlCfg.PassWord, mysqlCfg.Host, mysqlCfg.Port, mysqlCfg.DbName)
	pathWrite := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlCfg.User, mysqlCfg.PassWord, mysqlCfg.Host, mysqlCfg.Port, mysqlCfg.DbName)
	ormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       pathWrite, // DSN data source name
		DefaultStringSize:         256,       // string 类型字段的默认长度
		DisableDatetimePrecision:  false,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    false,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   false,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: true,      // 根据版本自动配置
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
	sqlDB, _ := global.DB.DB()
	sqlDB.SetMaxIdleConns(20)  // 设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) // 打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	// _db = db
	_ = global.DB.Use(dbresolver.
		Register(dbresolver.Config{
			// `db2` 作为 sources，`db3`、`db4` 作为 replicas
			Sources:  []gorm.Dialector{mysql.Open(pathRead)},                         // 读操作
			Replicas: []gorm.Dialector{mysql.Open(pathWrite), mysql.Open(pathWrite)}, // 写操作
			Policy:   dbresolver.RandomPolicy{},                                      // sources/replicas 负载均衡策略
		}))
	// 没必要设置，在DSN就设置好了
	// _db = _db.Set("gorm:table_options", "charset=utf8mb4")
	err = global.DB.AutoMigrate(&model.OrderGoods{}, &model.OrderInfo{}, &model.ShoppingCart{})
	if err != nil {
		// Todo log
		panic(err)
	}
}
