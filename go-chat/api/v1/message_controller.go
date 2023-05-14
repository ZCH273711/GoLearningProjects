package v1

import (
	"examples/go-chat/internal/service"
	"examples/go-chat/pkg/common/request"
	"examples/go-chat/pkg/common/response"
	"examples/go-chat/pkg/global/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetMessage(c *gin.Context) {
	log.Logger.Info(c.Query("uuid"))
	var messageRequest request.MessageRequest
	err := c.BindQuery(&messageRequest)
	if err != nil {
		log.Logger.Error("bindQueryError", log.Any("bindQueryError", err))
	}
	log.Logger.Info("messageRequest params: ", log.Any("messageRequest", messageRequest))

	message, err := service.MessageService.GetMessages(messageRequest)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(message))
}
