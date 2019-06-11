package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/solution9th/NSBridge/internal/nserr"
	"github.com/solution9th/NSBridge/internal/service/cache"
	"github.com/solution9th/NSBridge/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
			c.AbortWithStatusJSON(http.StatusOK, utils.ParseResult(nserr.WebErrUserNotLogin, "", ""))
			//c.Redirect(302, "/saml/login?RelayState=" + path)
			return
		}

		token, err := c.Cookie(TokenCookieKey)
		fmt.Println(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, utils.ParseResult(nserr.WebErrUserNotLogin, "", ""))
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
			c.AbortWithStatusJSON(http.StatusOK, utils.ParseResult(nserr.WebErrUserNotLogin, "", ""))
			//c.Redirect(302, "/saml/login?RelayState=" + path)
			return
		}

		c.Next()
	}
}
