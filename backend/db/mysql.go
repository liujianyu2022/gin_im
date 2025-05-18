package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gin_im/config"
)

func NewMySQLDB(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		mysql.Open(config.GetMySQLDSN()), 
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
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

	return db, nil
}
