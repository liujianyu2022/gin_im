package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gin_im/config"
	"gin_im/model"
)

func NewMySQLDB(config *config.Config) (*gorm.DB, error) {
	
	// 自定义日志模板，打印SQL语句
	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        // 彩色
		},
	)

	db, err := gorm.Open(
		mysql.Open(config.GetMySQLDSN()),
		&gorm.Config{
			Logger: customLogger,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MySQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MySQL.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.MySQL.ConnMaxLifetime) * time.Second)

	db.AutoMigrate(
		&model.User{},
		&model.Message{},
		&model.Contact{},
		&model.Group{},
	)

	return db, nil
}
