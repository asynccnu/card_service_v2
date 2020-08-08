package user

import (
	"encoding/json"
	"github.com/asynccnu/card_service_v2/handler"
	"github.com/asynccnu/card_service_v2/pkg/errno"
	"github.com/asynccnu/card_service_v2/service"
	"github.com/gin-gonic/gin"
)

// 输入的表单
type param struct {
	Limit string `json:"limit"`
	Page  string `json:"page"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type tempAccount struct {
	Result results `json:"result"`
}

type results struct {
	Rows []row `json:"rows"`
}

type row struct {
	DealName   string  `json:"dealName"`
	OrgName    string  `json:"orgName"`
	TransMoney float32 `json:"transMoney"`
	WalletName string  `json:"walletName"`
	DealDate   string  `json:"dealDate"`
	OutMoney   float32 `json:"outMoney"`
	InMoney    float32 `json:"inMoney"`
}

// Get gets an account by userid and password
func Account(c *gin.Context) {
	// 声明payload变量，因为BindJSON方法需要接收一个指针进行操作
	var data param
	var tempLists tempAccount
	var lists []row
	var message loginPayload

	if err := c.BindJSON(&message); err != nil {
		handler.SendError(c, errno.ErrBind, nil, err.Error())
		return
	}

	data.Limit = c.Query("limit")
	data.Page = c.Query("page")
	data.Start = c.Query("start")
	data.End = c.Query("end")

	// 检查失败的情况
	if err := service.ConfirmUser(message.UserId, message.Password); err != nil {
		handler.SendError(c, errno.ErrPasswordIncorrect, nil, err.Error())
		return
	}

	// 获得string格式的account
	temp, err := service.DoList(message.UserId, message.Password, data.Limit, data.Page, data.Start, data.End)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	err = json.Unmarshal([]byte(temp), &tempLists)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	for _, list := range tempLists.Result.Rows {
		lists = append(lists, list)
	}

	handler.SendResponse(c, nil, lists)

	return
}
