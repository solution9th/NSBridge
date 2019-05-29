package apiweb

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/solution9th/NSBridge/utils"
	"log"
	"time"
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

	/* 来自i1的req.body
	{"app":"23lubrx5ffu0k233",
	"i1_host":"org.i1.dev.com",
	"app_key":"216hh0hgthu8",
	"app_secret":"4t6zwwnwtl7f9ap8rs0327fw",
	"domain":"org.i1.dev.com",
	"data":{"超级管理员":"mafeng"}}
	*/
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
