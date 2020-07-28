package user

import (
	"github.com/asynccnu/card_service_v2/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type LoginPayload struct { // 用于接收payload的结构体
	User_id  string `json:"user_id"`
	Password string `json:"password"`
}

type CardInfos struct {
    // No string			`json:"no"`
	// DeptName   string	`json:"deptName"`
	StatusDesc	string	`json:"statusDesc"`
	Balance float32	`json:"balance"`
	// Xm string		`json:"xm"`
	// ValidityDate	string	`json:"validityDate"`
	// Status	string	`json:"status"`
	// Username	string	`json:"username"`
}

type Card struct {
	CardInfo	CardInfos	`json:"cardInfo"`
}

func Status(c *gin.Context) {
	var data LoginPayload // 声明payload变量，因为BindJSON方法需要接收一个指针进行操作
	var s Card
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request!",
			"err":     err, // 接收过程中的错误视为Bad Request
		})
		return
	}

	if !service.ConfirmUser(data.User_id, data.Password) { // 检查失败的情况
		c.JSON(401, gin.H{
			"message": "Password or account wrong.",
		})
		return
	} else {
		ret := service.DoStatus(data.User_id, data.Password)
		json.Unmarshal([]byte(ret), &s)
		c.JSON(200, gin.H{
			"message": "Authentiaction Success.",
			//"ret":   ret,
			"status":	s.CardInfo.StatusDesc,
			 "money": s.CardInfo.Balance,
		})
	}
	return
}

type LSParm struct {
	User_id  string `json:"user_id"`
	Password string `json:"password"`
	Limit string
	Page string
	Start string
	End string
}

type Account struct{
	Result	Results	`json:"result"`
}

type Results struct{
	Rows	[]Row
}

type Row struct{
	DealName		string
	OrgName			string
	TransMoney		float32
	WalletName		string
	DealDate		string
	OutMoney		float32
	InMoney			float32
}

func List(c *gin.Context) {
	var data LSParm // 声明payload变量，因为BindJSON方法需要接收一个指针进行操作
	var s Account
	var b []Row
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request!",
			"err":     err, // 接收过程中的错误视为Bad Request
		})
		return
	}

	if !service.ConfirmUser(data.User_id, data.Password) { // 检查失败的情况
		c.JSON(401, gin.H{
			"message": "Password or account wrong.",
		})
		return
	} else {
		ret := service.DoList(data.User_id, data.Password, data.Limit, data.Page, data.Start, data.End)
		json.Unmarshal([]byte(ret), &s)
		for _,val := range s.Result.Rows{
			b = append (b,val)
		}
		c.JSON(200, gin.H{
			"message": "Authentiaction Success.",
			"list":	b,
			//"ret":   ret,
		})
		
	}

	return
}