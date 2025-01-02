package db

import (
	"os"
	"seed-sync/log"
	"sync"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var once sync.Once

var DB *gorm.DB

// 数据库初始化 拿到DB实例

func InitDb() {
	once.Do(func() {
		var err error
		// 连接 SQLite 数据库
		// 数据库文件会保存在 data.db
		DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("数据库连接失败", zap.Error(err))
		}

		// 获取底层的 sqlDB
		sqlDB, err := DB.DB()
		if err != nil {
			log.Fatal("获取数据库实例失败", zap.Error(err))
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(10)      // 设置空闲连接池中的最大连接数
		sqlDB.SetMaxOpenConns(100)     // 设置打开数据库连接的最大数量
		sqlDB.SetConnMaxLifetime(3600) // 设置连接可复用的最大时间（秒）

		// 执行初始化SQL
		initSql()
	})
}

func initSql() {
	// 执行初始化SQL
	initSQL, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatal("读取初始化SQL文件失败", zap.Error(err))
	}

	// 执行SQL语句
	if err := DB.Exec(string(initSQL)).Error; err != nil {
		log.Fatal("执行初始化SQL失败", zap.Error(err))
	}
}
