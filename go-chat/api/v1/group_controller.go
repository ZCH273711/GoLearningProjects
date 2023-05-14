package v1

import (
	"examples/go-chat/config"
	"examples/go-chat/internal/service"
	"examples/go-chat/pkg/common/response"
	"examples/go-chat/pkg/global/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strings"
)

func GetFile(c *gin.Context) {
	fileName := c.Param("fileName")
	log.Logger.Info(fileName)
	data, _ := os.ReadFile(config.GetConfig().StaticPath.FilePath + fileName)
	c.Writer.Write(data)
}

func SaveFile(c *gin.Context) {
	namePreffix := uuid.New().String()
	userUuid := c.PostForm("uuid")
	file, _ := c.FormFile("file")
	fileName := file.Filename
	index := strings.LastIndex(fileName, ".")
	suffix := fileName[index:]

	newFileName := namePreffix + suffix
	log.Logger.Info("file", log.Any("file name", config.GetConfig().StaticPath.FilePath+newFileName))
	log.Logger.Info("userUuid", log.Any("userUuid name", userUuid))

	c.SaveUploadedFile(file, config.GetConfig().StaticPath.FilePath+newFileName)
	err := service.UserService.ModifyUserAvatar(newFileName, userUuid)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(newFileName))
}
