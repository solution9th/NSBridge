package apiweb

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/solution9th/NSBridge/config"
	"github.com/solution9th/NSBridge/service/cache"
	"github.com/solution9th/NSBridge/utils"
	"net/http"
)

func Logout(c *gin.Context) {

	session := sessions.Default(c)
	userid := session.Get(UserSessionKey)
	session.Delete(UserSessionKey)
	session.Save()

	token, err := c.Cookie(TokenCookieKey)
	c.SetCookie(TokenCookieKey, "", -1, "/", config.SamlConfig.Domain, false, true)

	err = cache.DefaultCache.Delete(fmt.Sprintf(passportTokenFormat, userid, token))
	if err != nil {
		utils.Error("DefaultCache Del Err", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"msg": errors.New("DefaultCache Del Err"),
		})
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccess())
	return
}
