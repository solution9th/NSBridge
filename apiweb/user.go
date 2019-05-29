package apiweb

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/solution9th/NSBridge/models"
	"github.com/solution9th/NSBridge/utils"
	"net/http"
	"strings"
)

func GetUser(c *gin.Context)  {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ParseResult(models.WebErrUserNotLogin, "", ""))
		return
	}

	userIdNames := strings.Split(user.(string), ":")
	if len(userIdNames) != 2 {
		c.AbortWithStatusJSON(http.StatusOK, utils.ParseResult(models.WebErrSrever, "", ""))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccessWithData(gin.H{
		"user_id": userIdNames[0],
		"user_name": userIdNames[1],
	}))
	return
}
