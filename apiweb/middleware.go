package apiweb

import (
	"errors"
	"fmt"
	"github.com/solution9th/NSBridge/models"
	"github.com/solution9th/NSBridge/utils"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/solution9th/NSBridge/service/cache"

	"github.com/gin-contrib/sessions"
)

var (
	// dns:web:api:passport:user_id:token
	passportTokenFormat = "dns:web:api:passport:%s:%s"
)

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//path := c.Request.URL.Path
		session := sessions.Default(c)
		user := session.Get(UserSessionKey)
		fmt.Println(user)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusOK, utils.ParseResult(models.WebErrUserNotLogin, "", ""))
			//c.Redirect(302, "/saml/login?RelayState=" + path)
			return
		}

		token, err := c.Cookie(TokenCookieKey)
		fmt.Println(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, utils.ParseResult(models.WebErrUserNotLogin, "", ""))
			//c.Redirect(302, "/saml/login?RelayState=" + path)
			return
		}
		exist, err := cache.DefaultCache.Exist(fmt.Sprintf(passportTokenFormat, user, token))
		if err != nil {
			utils.Error("cache.DefaultCache.Exist", err.Error())
			c.AbortWithStatusJSON(500, gin.H{
				"msg": errors.New("cache.DefaultCache.Exist Err"),
			})
			return
		}

		if !exist {
			c.AbortWithStatusJSON(http.StatusOK, utils.ParseResult(models.WebErrUserNotLogin, "", ""))
			//c.Redirect(302, "/saml/login?RelayState=" + path)
			return
		}

		c.Next()
	}
}
