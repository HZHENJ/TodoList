package db

import (
	"fmt"
	"log"
	"time"
	"to-do-list/internal/repository/db/model"
	"to-do-list/pkg/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() {
	c := conf.Config.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.DbName,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)

	var dbLogger logger.Interface
	if conf.Config.Service.AppMode == "debug" {
		dbLogger = logger.Default.LogMode(logger.Info)
	} else {
		dbLogger = logger.Default
	}

	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: dbLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Panicln("数据库连接失败:", err)
	}

	sqlDB, err := DB.DB()
	sqlDB.SetMaxIdleConns(200)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	// 自动迁移
	DB.AutoMigrate(&model.User{}, &model.Task{})
}
