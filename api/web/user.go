package web

import (
	"net/http"
	"strings"

	"github.com/solution9th/NSBridge/internal/nserr"
	"github.com/solution9th/NSBridge/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(UserSessionKey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.ParseResult(nserr.WebErrUserNotLogin, "", ""))
		return
	}

	userIdNames := strings.Split(user.(string), ":")
	if len(userIdNames) != 2 {
		c.AbortWithStatusJSON(http.StatusOK, utils.ParseResult(nserr.WebErrSrever, "", ""))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccessWithData(gin.H{
		"user_id":   userIdNames[0],
		"user_name": userIdNames[1],
	}))
	return
}
