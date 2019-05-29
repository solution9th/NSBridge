package v1

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/solution9th/NSBridge/config"
	"github.com/solution9th/NSBridge/dns"
	"github.com/solution9th/NSBridge/service/database"
	"github.com/solution9th/NSBridge/utils"

	"github.com/gin-gonic/gin"
)

const (
	// APIKeyName 请求头中的 api key
	APIKeyName       = "X-API-KEY"
	APIHMACtName     = "X-API-HMAC"
	APITimeStampName = "X-API-TIMESTAMP"
)

var (
	InvalidError = errors.New("invalid api-key")
)

// AuthMiddleware api auth middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if config.IsDebug {

			c.Next()
			return
		}

		apiKey := c.GetHeader(APIKeyName)
		apiHMAC := c.GetHeader(APIHMACtName)
		apiTimestamp := c.GetHeader(APITimeStampName)

		if apiKey == "" || apiHMAC == "" || apiTimestamp == "" {
			utils.Debugf("miss params key>%v< hmac>%v< time>%v<", apiKey, apiHMAC, apiTimestamp)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		ts, err := strconv.ParseInt(apiTimestamp, 10, 64)
		if err != nil {
			utils.Errorf("ts error: %v", apiTimestamp)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		now := time.Now()

		timeout := 5 * time.Hour

		if ts > now.Add(timeout).Unix() || ts < now.Add(-1*timeout).Unix() {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		secret := getSecret(apiKey)
		if secret == "" {
			c.AbortWithError(http.StatusForbidden, InvalidError)
			return
		}

		apiHMAC = strings.ToLower(apiHMAC)

		if !isOkSign(c, secret, apiHMAC) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
		return
	}
}

// API-HMAC=md5(请求方法 + API-KEY + 接口请求 url + 请求参数体 + API-TIMESTAMP + SECRET-KEY);
func isOkSign(c *gin.Context, secret, hmac string) bool {

	apiKey := c.GetHeader(APIKeyName)
	apiTimestamp := c.GetHeader(APITimeStampName)
	uri := c.Request.RequestURI
	method := strings.ToUpper(c.Request.Method)

	body, err := c.GetRawData()
	if err != nil {
		utils.Errorf("get raw data error: %v", err)
		return false
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	s := fmt.Sprintf("%s%s%s%s%s%s", method, apiKey, uri, string(body), apiTimestamp, secret)

	sign := fmt.Sprintf("%x", md5.Sum([]byte(s)))

	if sign == hmac {
		utils.Infof("success: %v key: %s method: %s, uri: %s, body: %v", apiTimestamp, apiKey, method, uri, string(body))
		return true
	}

	utils.Infof("fail: %v key: %s method: %s, uri: %s, body: %v", apiTimestamp, apiKey, method, uri, string(body))
	utils.Errorf("sign: %v, hmac: %v", sign, hmac)

	return false
}

// 根据 key 的前缀区分 key 的类型
func getSecret(key string) string {

	db := database.New()
	if strings.HasPrefix(key, dns.GetRecordKeyPrefix()) {
		// 说明是操作 记录 的key
		m, err := db.GetDomainByRecordKey(key)
		if err != nil {
			utils.Error("get domain k: %v, error:", key, err)
			return ""
		}
		return m.RecordSecret
	}

	// 说明是添加域名的 key
	domain, err := db.GetAuthByKey(key)
	if err != nil {
		utils.Error("get domain k: %v, error:", key, err)
		return ""
	}

	return domain.DomainSecret
}
