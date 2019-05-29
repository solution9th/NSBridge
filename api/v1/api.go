package v1

import (
	"net/http"
	"time"

	"github.com/solution9th/NSBridge/models"
	"github.com/solution9th/NSBridge/utils"

	"github.com/gin-gonic/gin"
)

// Ping test router
func Ping(c *gin.Context) {

	c.JSON(http.StatusOK, utils.ParseResult(models.ErrErr, "", map[string]interface{}{
		"msg":       "ok",
		"create_at": time.Now().Unix(),
	}))
	return
}
