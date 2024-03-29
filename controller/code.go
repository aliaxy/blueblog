// Package controller 状态码定义
package controller

// ResCode 返回类型定义
type ResCode int64

const (
	CodeSuccess         ResCode = 1_000 + iota // CodeSuccess 成功
	CodeInvalidParams                          // CodeInvalidParams 非法参数
	CodeUserExist                              // CodeUserExist 用户存在
	CodeUserNotExist                           // CodeUserNotExist 用户不存在
	CodeInvalidPassword                        // CodeInvalidPassword 密码错误
	CodeServerBusy                             // CodeServerBusy 服务繁忙

	CodeInvalidToken // CodeInvalidToken 认证失败
	CodeNeedLogin    // CodeNeedLogin 需要登录
)

var codeMessageMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务器繁忙",

	CodeInvalidToken: "无效的Token",
	CodeNeedLogin:    "需要登录",
}

// Msg 映射错误码
func (c ResCode) Msg() string {
	msg, ok := codeMessageMap[c]
	if !ok {
		msg = codeMessageMap[CodeServerBusy]
	}
	return msg
}
