package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/solution9th/NSBridge/apiweb"
	"github.com/solution9th/NSBridge/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

// NewWebRouter new router
func NewWebRouter(r *gin.Engine) {

	// session
	store, err := redis.NewStore(10, "tcp", fmt.Sprintf("%s:%s", config.RedisConfig.Host, config.RedisConfig.Port), config.RedisConfig.Passwd, []byte("secret"))
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions("dns_session", store))

	samlRouter := r.Group("/saml")
	samlRouter.Use()
	{
		samlRouter.GET("/login", apiweb.SAMLLogin)
		samlRouter.POST("/acs", apiweb.SAMLAcs)
		samlRouter.PATCH("/acs", apiweb.GetI1Notice)
	}

	// web
	webapiRouter := r.Group("/web")
	webapiRouter.Use(apiweb.AuthMiddleware())
	{

		// 授权管理
		authRouter := webapiRouter.Group("/auth")
		{
			authRouter.GET("", apiweb.SearchAuthInfo)         // 获取auth列表
			authRouter.POST("", apiweb.CreateAuth)            // 新增授权
			authRouter.PUT("", apiweb.RemarkAuth)             // 更新remark
			authRouter.PUT("/:auth_id", apiweb.DisableAuth)   // 起禁用
			authRouter.DELETE("/:auth_id", apiweb.DeleteAuth) // 删除域名
		}
		// 域名管理
		domainRouter := webapiRouter.Group("/domain")
		{
			domainRouter.GET("", apiweb.GetDomainList) // 获取domain列表
		}

		// 解析记录
		recordRouter := webapiRouter.Group("/record")
		{
			recordRouter.GET("/:domain_id", apiweb.SearchRecords) // 获取某个域名下的解析记录列表
			recordRouter.GET("/:domain_id/types", apiweb.GetDomainTypes)
		}

		webapiRouter.GET("/logout", apiweb.Logout)
		webapiRouter.GET("/user", apiweb.GetUser)
	}

}
