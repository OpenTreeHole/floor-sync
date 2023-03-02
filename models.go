package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

type Map map[string]any

type Floor struct {
	ID        int       `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
}

type Floors []*Floor

type Hole struct {
	ID     int  `json:"id"`
	Hidden bool `json:"hidden"`
}

type Holes []*Hole

var DB *gorm.DB

var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
	},
	Logger: logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	),
}

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(Config.DbUrl), gormConfig)
	if err != nil {
		log.Fatalf("error failed to connect to database: %s", err)
	}
}
