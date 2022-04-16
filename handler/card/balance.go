package card

import (
	"github.com/asynccnu/card_service_v2/handler"
	"github.com/asynccnu/card_service_v2/pkg/errno"
	"github.com/asynccnu/card_service_v2/service"

	"github.com/gin-gonic/gin"
)

// 获取余额
func Balance(c *gin.Context) {
	sid := c.MustGet("Sid").(string)
	password := c.MustGet("Password").(string)

	// 获取一卡通信息，包括余额信息
	cardInfo, err := service.GetCardInfo(sid, password)
	if err != nil {
		// 验证是否是账号密码错误
		if err == errno.ErrPasswordIncorrect {
			handler.SendResponse(c, errno.ErrPasswordIncorrect, nil)
			return
		}

		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, cardInfo)
}
