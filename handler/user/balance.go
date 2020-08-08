package user

import (
	"encoding/json"
	"github.com/asynccnu/card_service_v2/handler"
	"github.com/asynccnu/card_service_v2/pkg/errno"
	"github.com/asynccnu/card_service_v2/service"
	"github.com/gin-gonic/gin"
)

// 用于接收payload的结构体
type loginPayload struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

// 用于解析网页json格式数据，所以json不能改下划线
type cardInfo struct {
	// No           string  `json:"no"`
	// DeptName     string  `json:"deptName"`
	StatusDesc      string  `json:"statusDesc"`
	Balance         float32 `json:"balance"`
	// Xm           string  `json:"xm"`
	// ValidityDate string  `json:"validityDate"`
	// Status       string  `json:"status"`
	// Username     string  `json:"username"`
}

type card struct {
	CardInfo cardInfo `json:"cardInfo"`
}

// Get  gets status and money by userid and password
func Balance(c *gin.Context) {
	// 声明payload变量，因为BindJSON方法需要接收一个指针进行操作
	var data loginPayload
	var s card

	if err := c.BindJSON(&data); err != nil {
		handler.SendError(c, errno.ErrBind, nil, err.Error())
		return
	}

	// 检查失败的情况
	if err := service.ConfirmUser(data.UserId, data.Password); err != nil {
		handler.SendError(c, errno.ErrPasswordIncorrect, nil, err.Error())
		return
	}

	ret, err := service.DoStatus(data.UserId, data.Password)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	err = json.Unmarshal([]byte(ret), &s)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, s)

	return
}
