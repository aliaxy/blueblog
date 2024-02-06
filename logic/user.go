package logic

import (
	"blueblog/dao/mysql"
	"blueblog/models"
	"blueblog/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return
	}

	// 生成 UID
	userId := snowflake.GenID()

	// 构造 User 实例
	user := &models.User{
		UserID:   userId,
		Username: p.Username,
		Password: p.Password,
	}

	// 保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.Login(user)
}
