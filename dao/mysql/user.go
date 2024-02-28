// Package mysql 用户相关
package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"blueblog/models"
)

const secret = "aliaxyblueblog"

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(*) from user where username = ?"
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 插入用户
func InsertUser(user *models.User) (err error) {
	user.Password = encryptPassword(user.Password)

	sqlStr := "insert into user(user_id, username, password) values(?, ?, ?)"
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// Login 登录
func Login(user *models.User) (err error) {
	oPassword := user.Password
	sqlStr := "select user_id, username, password from user where username = ?"
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return
	}
	if encryptPassword(oPassword) != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// GetUserByID 根据用户 ID 获取信息
func GetUserByID(id int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := "select user_id, username from user where user_id = ?"

	err = db.Get(user, sqlStr, id)
	return
}
