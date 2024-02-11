package mysql_test

import (
	"testing"

	"blueblog/dao/mysql"
	"blueblog/models"
	"blueblog/settings"
)

func init() {
	dbCfg := &settings.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "211010",
		DB:           "blueblog",
		Port:         13306,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := mysql.Init(dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := &models.Post{
		ID:          1,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a test",
	}
	err := mysql.CreatePost(post)
	if err != nil {
		t.Fatal("CreatePost insert record into mysql failed, err: " + err.Error())
	}
	t.Log("CreatePost insert record into mysql success")
}
