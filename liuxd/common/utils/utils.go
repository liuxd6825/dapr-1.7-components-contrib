package utils

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

func NewTimeNowString() string {
	return time.Now().Format("2006-01-02 15:04:05.505072")
}

func NewMongoNow() primitive.DateTime {
	return primitive.NewDateTimeFromTime(time.Now())
}

func SetHttpSuccess(ctx *fasthttp.RequestCtx, data interface{}, err error) {
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	var bytes []byte
	if data != nil {
		bytes, err = json.Marshal(data)
	}
	ctx.Success("application/json", bytes)
}

//
// AsMongoName
// @Description: 驼峰转蛇形
// @param s 要转换的字符串
// @return string
//
func AsMongoName(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	res := strings.ToLower(string(data[:]))
	if strings.HasPrefix(res, "_") {
		return res[1:]
	}
	return res
}
