package web

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/solution9th/NSBridge/internal/utils"

	"github.com/gin-gonic/gin"
)

type DnsApp struct {
	App        string            `json:"app"`
	I1Host     string            `json:"i1_host"`
	AppKey     string            `json:"app_key"`
	AppSecrect string            `json:"app_secrect"`
	Domain     string            `json:"domain"`
	Data       map[string]string `json:"data"`
}

func InitInstance(c *gin.Context) {

	data, err := c.GetRawData()
	if err != nil {
		utils.Error("InitInstance err:", err)
	}
	var dns_app DnsApp
	json.Unmarshal(data, &dns_app)
	log.Println(string(data))

	c.JSON(200, gin.H{
		"login":     fmt.Sprintf("%s/login/%s", AppUrl, dns_app.AppKey),
		"callback":  CallBackUrl,
		"timestamp": time.Now().Format(time.RFC3339),
	})

}
