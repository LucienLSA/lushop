package initialize

import (
	"fmt"
	"log"
	"usersrv/global"
	"usersrv/model"

	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func MySQL() {
	mysqlCfg := global.ServerConfig.MySQLInfo
	pathRead := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlCfg.User, mysqlCfg.PassWord, mysqlCfg.Host, mysqlCfg.Port, mysqlCfg.DbName)
	// pathWrite := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	mysqlCfg.User, mysqlCfg.PassWord, mysqlCfg.Host, mysqlCfg.Port, mysqlCfg.DbName)
	ormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       pathRead, // DSN data source name
		DefaultStringSize:         256,      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,    // 根据版本自动配置
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
	global.DB = db
	// _ = global.DB.Use(dbresolver.
	// 	Register(dbresolver.Config{
	// 		// `db2` 作为 sources，`db3`、`db4` 作为 replicas
	// 		Sources:  []gorm.Dialector{mysql.Open(pathRead)},                         // 读操作
	// 		Replicas: []gorm.Dialector{mysql.Open(pathWrite), mysql.Open(pathWrite)}, // 写操作
	// 		Policy:   dbresolver.RandomPolicy{},                                      // sources/replicas 负载均衡策略
	// 	}))
	// global.DB = global.DB.Set("gorm:table_options", "charset=utf8mb4")
	err = global.DB.AutoMigrate(&model.User{})
	if err != nil {
		// Todo log
		panic(err)
	}
}
