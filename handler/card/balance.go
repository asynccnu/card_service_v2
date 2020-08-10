package card

import (
	"encoding/json"

	"github.com/asynccnu/card_service_v2/handler"
	"github.com/asynccnu/card_service_v2/pkg/errno"
	"github.com/asynccnu/card_service_v2/service"

	"github.com/gin-gonic/gin"
)

// 用于接收payload的结构体
// type loginPayload struct {
// 	UserId   string `json:"user_id"`
// 	Password string `json:"password"`
// }

// 用于解析网页json格式数据，所以json不能改下划线
type cardInfo struct {
	// No           string  `json:"no"`
	// DeptName     string  `json:"deptName"`
	StatusDesc string  `json:"statusDesc"`
	Balance    float32 `json:"balance"`
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

	sid := c.MustGet("Sid").(string)
	password := c.MustGet("Password").(string)

	// 检查失败的情况
	if err := service.ConfirmUser(sid, password); err != nil {
		handler.SendError(c, errno.ErrPasswordIncorrect, nil, err.Error())
		return
	}

	temp, err := service.GetCardInfo(sid, password)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	var tempCard card

	err = json.Unmarshal([]byte(temp), &tempCard)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, tempCard)
}
