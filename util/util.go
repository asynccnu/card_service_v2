package util

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
)

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestID, ok := v.(string); ok {
		return requestID
	}
	return ""
}

// 获取当前时间，北京时间，东八区
func GetCurrentTime() *time.Time {
	// 集群上使用 UTC，需要手动加上时差
	t := time.Now().UTC().Add(8 * time.Hour)
	return &t
}

// 获取当前日期，格式 2006-01-02
func GetCurrentDate() string {
	t := GetCurrentTime()
	return t.Format("2006-01-02")
}
