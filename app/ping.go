package app

import (
	"github.com/solution9th/NSBridge/internal/config"

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
