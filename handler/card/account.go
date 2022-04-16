package card

import (
	"github.com/asynccnu/card_service_v2/handler"
	"github.com/asynccnu/card_service_v2/pkg/errno"
	"github.com/asynccnu/card_service_v2/service"
	"github.com/asynccnu/card_service_v2/util"

	"github.com/gin-gonic/gin"
)

type AccountResponseData struct {
	Count int                `json:"count"`
	List  []*service.DealRow `json:"list"`
}

// Get gets an account by userid and password
func Account(c *gin.Context) {
	sid := c.MustGet("Sid").(string)
	password := c.MustGet("Password").(string)

	// 获取当前日期
	currentDate := util.GetCurrentDate()

	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")
	start := c.DefaultQuery("start", currentDate)
	end := c.DefaultQuery("end", currentDate)

	// 获得string格式的account
	records, err := service.GetConsumeList(sid, password, limit, page, start, end)
	if err != nil {
		// 验证是否是账号密码错误
		if err == errno.ErrPasswordIncorrect {
			handler.SendResponse(c, errno.ErrPasswordIncorrect, nil)
			return
		}

		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, &AccountResponseData{
		Count: len(records),
		List:  records,
	})
}
