package user

import (
	"strconv"

	. "github.com/asynccnu/card_service_v2/handler"
	"github.com/asynccnu/card_service_v2/log"
	"github.com/asynccnu/card_service_v2/model"
	"github.com/asynccnu/card_service_v2/pkg/errno"
	"github.com/asynccnu/card_service_v2/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Update update a exist user account info.
func Update(c *gin.Context) {
	var (
		postData model.UserModel
	)

	log.Info("Update function called.",
		zap.String("X-Request-Id", util.GetReqID(c)))
	// Get the user id from the url parameter.
	userID, _ := strconv.Atoi(c.Param("id"))

	if user, err := model.GetUserById(uint64(userID)); err != nil {
		SendError(c, errno.ErrUserNotFound, nil, err.Error())
		return
	} else {
		if err := c.Bind(&postData); err != nil {
			SendBadRequest(c, errno.ErrBind, nil, err.Error())
			return
		}

		// 更新数据
		user.Username = postData.Username
		user.Password = postData.Password

		// Validate the data.
		if err := user.Validate(); err != nil {
			SendError(c, errno.ErrValidation, nil, err.Error())
			return
		}

		// Encrypt the user password.
		if err := user.Encrypt(); err != nil {
			SendError(c, errno.ErrEncrypt, nil, err.Error())
			return
		}

		// Save changed fields.
		if err := user.Update(); err != nil {
			SendError(c, errno.ErrDatabase, nil, err.Error())
			return
		}

		SendResponse(c, nil, nil)
	}

}