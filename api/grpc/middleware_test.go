package grpc

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/solution9th/NSBridge/internal/config"
	"github.com/solution9th/NSBridge/internal/service/cache"
	"github.com/solution9th/NSBridge/internal/service/mysql"
	"github.com/solution9th/NSBridge/internal/utils"

	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
)

var (
	KSMap = map[string]string{
		"dodomain":                 "dodomainmain",
		"cmVjb3Jk2cab9d2edfd19539": "H1dvpIePyRivG0uR154M",
	}
)

// Init 初始化项目
func InitDB(t *testing.T) (err error) {

	fileName := "config"

	err = config.InitConfig(fileName)
	if err != nil {
		return err
	}

	var (
		redisHost   = config.RedisConfig.Host
		redisPort   = config.RedisConfig.Port
		redisPasswd = config.RedisConfig.Passwd
		redisDB     = config.RedisConfig.DB

		mysqlHost   = config.MySQLConfig.Host
		mysqlPort   = config.MySQLConfig.Port
		mysqlUser   = config.MySQLConfig.User
		mysqlPasswd = config.MySQLConfig.Passwd
		mysqlDB     = config.MySQLConfig.DBName
	)

	// 初始化 redis
	err = cache.InitRedis(redisHost, redisPort, redisPasswd, redisDB)
	if err != nil {
		utils.Error("[Init] initRedis error:", err)
		return err
	}

	// 初始化 默认MySQL
	err = mysql.InitDefaultDB(mysqlDB, mysqlUser, mysqlPasswd, mysqlHost, mysqlPort)
	if err != nil {
		utils.Error("[init] initMySQL error:", err)
		return err
	}

	fmt.Println("ooooooooooo")

	return nil
}

func Init(t *testing.T) {

	err := InitDB(t)
	if err != nil {
		t.Error("init error:", err)
	}

}

// config debug 为 false 时
func TestMiddleware(t *testing.T) {

	// init
	Init(t)

	r := gin.New()
	gin.SetMode(gin.TestMode)
	r.GET("/", AuthMiddleware())
	r.POST("/", AuthMiddleware())

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/").Expect().Status(403)

	// start test

	now := time.Now()

	tests := []struct {
		method    string
		headerMap map[string]string
		post      map[string]interface{}
		code      int
	}{
		// normal test
		{
			method: "GET",
			headerMap: map[string]string{
				APIKeyName:       "dodomain",
				APITimeStampName: fmt.Sprintf("%d", now.Unix()),
				APIHMACName:      testGetSign("dodomain", "GET", "/", nil, now),
			},
			post: nil,
			code: 200,
		},
		// normal test record
		{
			method: "GET",
			headerMap: map[string]string{
				APIKeyName:       "cmVjb3Jk2cab9d2edfd19539",
				APITimeStampName: fmt.Sprintf("%d", now.Unix()),
				APIHMACName:      testGetSign("cmVjb3Jk2cab9d2edfd19539", "GET", "/", nil, now),
			},
			post: nil,
			code: 200,
		},
		// test cache
		{
			method: "GET",
			headerMap: map[string]string{
				APIKeyName:       "dodomain",
				APITimeStampName: fmt.Sprintf("%d", now.Unix()),
				APIHMACName:      testGetSign("dodomain", "GET", "/", nil, now),
			},
			post: nil,
			code: 200,
		},
		// test post
		{
			method: "POST",
			headerMap: map[string]string{
				APIKeyName:       "dodomain",
				APITimeStampName: fmt.Sprintf("%d", now.Unix()),
				APIHMACName: testGetSign("dodomain", "POST", "/", map[string]interface{}{
					"k": "k",
				}, now),
			},
			post: map[string]interface{}{
				"k": "k",
			},
			code: 200,
		},
		// test post
		{
			method: "POST",
			headerMap: map[string]string{
				APIKeyName:       "dodomain",
				APITimeStampName: fmt.Sprintf("%d", now.Unix()),
				APIHMACName: testGetSign("dodomain", "POST", "/", map[string]interface{}{
					"k": "k",
				}, now),
			},
			post: map[string]interface{}{
				"k": "k11111111111",
			},
			code: 403,
		},
		{
			method: "GET",
			headerMap: map[string]string{
				APIKeyName:       "dodomain",
				APITimeStampName: fmt.Sprintf("%d", now.Unix()),
				APIHMACName:      testGetSign("dodomain", "GET", "/123123", nil, now),
			},
			post: nil,
			code: 403,
		},
	}

	for _, v := range tests {

		req := e.Request(v.method, "/").WithHeaders(v.headerMap)

		if v.post != nil {
			req.WithJSON(v.post)
		}

		req.Expect().Status(v.code)
	}

	cleanCache()
}

// ks 默认为 dodomian = dodomainmain
//
func testGetSign(key, method, url string, body map[string]interface{}, now time.Time) string {

	secret := KSMap[key]

	var bodyStr string

	if body != nil && method != "GET" {
		bb, _ := json.Marshal(body)
		bodyStr = string(bb)
	}

	s := fmt.Sprintf("%s%s%s%s%d%s", method, key, url, bodyStr, now.Unix(), secret)

	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func cleanCache() {
	cache.DefaultCache.Cache.FlushDB()
}
