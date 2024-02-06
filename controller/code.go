package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1_000 + iota
	CodeInvalidParams
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
)

var codeMessageMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务器繁忙",
}

func (c ResCode) Msg() string {
	msg, ok := codeMessageMap[c]
	if !ok {
		msg = codeMessageMap[CodeServerBusy]
	}
	return msg
}
