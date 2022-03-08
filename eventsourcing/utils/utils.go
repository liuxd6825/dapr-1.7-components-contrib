package utils

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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
