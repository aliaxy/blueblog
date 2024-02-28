// Package mysql 数据库驱动
package mysql

import (
	"fmt"

	"blueblog/settings"

	_ "github.com/go-sql-driver/mysql" // 导入 mysql 驱动
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// 数据库
var db *sqlx.DB

// Init 初始化数据库连接
func Init(cfg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return
}

// Close 关闭连接
func Close() {
	_ = db.Close()
}
