package v1

import (
	"examples/go-chat/internal/model"
	"examples/go-chat/internal/service"
	"examples/go-chat/pkg/common/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetGroup(c *gin.Context) {
	uuid := c.Param("uuid")
	groups, err := service.GroupService.GetGroups(uuid)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(groups))
}

func SaveGroup(c *gin.Context) {
	uuid := c.Param("uuid")
	var group model.Group
	c.ShouldBindJSON(&group)
	service.GroupService.SaveGroup(uuid, group)
	c.JSON(http.StatusOK, response.SuccessMsg(nil))
}

func JoinGroup(c *gin.Context) {
	userUuid := c.Param("userUuid")
	groupUuid := c.Param("groupUuid")
	err := service.GroupService.JoinGroup(groupUuid, userUuid)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(nil))
}

func GetGroupUsers(c *gin.Context) {
	groupUuid := c.Param("uuid")
	users := service.GroupService.GetUserIdByGroupId(groupUuid)
	c.JSON(http.StatusOK, response.SuccessMsg(users))
}
