package app

import (
	"fmt"

	"github.com/solution9th/NSBridge/api/web"
	"github.com/solution9th/NSBridge/internal/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
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
		samlRouter.GET("/login", web.SAMLLogin)
		samlRouter.POST("/acs", web.SAMLAcs)
		samlRouter.PATCH("/acs", web.GetI1Notice)
	}

	// web
	webapiRouter := r.Group("/web")
	webapiRouter.Use(web.AuthMiddleware())
	{

		// 授权管理
		authRouter := webapiRouter.Group("/auth")
		{
			authRouter.GET("", web.SearchAuthInfo)         // 获取auth列表
			authRouter.POST("", web.CreateAuth)            // 新增授权
			authRouter.PUT("", web.RemarkAuth)             // 更新remark
			authRouter.PUT("/:auth_id", web.DisableAuth)   // 起禁用
			authRouter.DELETE("/:auth_id", web.DeleteAuth) // 删除域名
		}
		// 域名管理
		domainRouter := webapiRouter.Group("/domain")
		{
			domainRouter.GET("", web.GetDomainList) // 获取domain列表
		}

		// 解析记录
		recordRouter := webapiRouter.Group("/record")
		{
			recordRouter.GET("/:domain_id", web.SearchRecords) // 获取某个域名下的解析记录列表
			recordRouter.GET("/:domain_id/types", web.GetDomainTypes)
		}

		webapiRouter.GET("/logout", web.Logout)
		webapiRouter.GET("/user", web.GetUser)
	}

}
