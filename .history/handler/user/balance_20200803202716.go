package user

import (
	"github.com/asynccnu/card_service_v2/service"
	"github.com/asynccnu/card_service_v2/handler"
	"github.com/asynccnu/card_service_v2/pkg/errno"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// 用于接收payload的结构体
type LoginPayload struct { 
	User_id  		 string 	 `json:"user_id"`
	Password		 string		 `json:"password"`
}

type CardInfo struct {
    // No 			string		`json:"no"`
	// DeptName   	string		`json:"deptName"`
	StatusDesc		string		`json:"statusDesc"`
	Balance 		float32		`json:"balance"`
	// Xm 			string		`json:"xm"`
	// ValidityDate	string		`json:"validityDate"`
	// Status		string		`json:"status"`
	// Username		string		`json:"username"`
}

type Card struct {
	CardInfo	CardInfo	`json:"cardInfo"`
}


func Balance(c *gin.Context) {
	var data LoginPayload // 声明payload变量，因为BindJSON方法需要接收一个指针进行操作
	var s Card
	if err := c.BindJSON(&data); err != nil {
		handler.SendError(c,errno.ErrBind,nil,err.Error())
		return
	}

	if !service.ConfirmUser(data.User_id, data.Password) { // 检查失败的情况
		c.JSON(401, gin.H{
			"message": "Password or account wrong.",
		})
		return
	}
		
	ret := service.DoStatus(data.User_id, data.Password)
	json.Unmarshal([]byte(ret), &s)
	c.JSON(200, gin.H{
		"message": "Authentiaction Success.",
		"status":	s.CardInfo.StatusDesc,
		 "money": 	s.CardInfo.Balance,
				})
	}
	return
}

