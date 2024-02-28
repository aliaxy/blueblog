// Package mysql 错误代码
package mysql

import "errors"

var (
	// ErrorUserExist 用户已存在
	ErrorUserExist = errors.New("用户已存在")
	// ErrorUserNotExist 用户不存在
	ErrorUserNotExist = errors.New("用户不存在")
	// ErrorInvalidPassword 用户名或密码错误
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	// ErrorInvalidID 无效 id
	ErrorInvalidID = errors.New("无效的ID")
)
