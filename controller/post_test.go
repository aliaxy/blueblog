package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blueblog/controller"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, controller.CreatePostHandler)

	body := `{
		"title": "test",
		"content": "just a test"
	}`

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 方法一 判断响应内容是否包含指定字符串
	// assert.Contains(t, w.Body.String(), "需要登录")

	// 方法二 将响应的内容反序列化到 ResponseData 然后判断字段与预期是否一致
	res := new(controller.ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatal("json.Unmarshal w.Body failed, err: " + err.Error())
	}

	assert.Equal(t, controller.CodeInvalidParams, res.Code)
}
