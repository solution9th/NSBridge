package app

import (
	"github.com/solution9th/NSBridge/config"

	"github.com/gin-gonic/gin"
)

var (
	version   = "default"
	buildDate = ""
)

func Ping(c *gin.Context) {
	// fmt.Println("Ping")
	c.Writer.Header().Set("Build-Version", version)
	c.Writer.Header().Set("Build-Date", buildDate)
	c.Writer.Header().Set("Build-Status", config.ServerConfig.Status)
	c.String(200, "%s", "dns ok")
	return
}

// Gin 中间件测试
// r.GET("/ping", TA(), TB(), Ping)
// TA
// TB
// Ping
// TB-
// TA-
// func TA() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Println("TA")
// 		c.Next()
// 		fmt.Println("TA-")
// 	}
// }

// func TB() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Println("TB")
// 		c.Next()
// 		fmt.Println("TB-")
// 	}
// }
