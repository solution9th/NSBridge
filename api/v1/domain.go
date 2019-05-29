package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/solution9th/NSBridge/dns"
	"github.com/solution9th/NSBridge/models"
	"github.com/solution9th/NSBridge/sdk"
	"github.com/solution9th/NSBridge/service/cache"
	"github.com/solution9th/NSBridge/service/database"
	"github.com/solution9th/NSBridge/utils"

	"github.com/gin-gonic/gin"
	"github.com/haozibi/gendry/scanner"
)

// GetDomainLists 获得域名列表
func GetDomainLists(c *gin.Context) {

	var err error

	startStr := c.Query("start")
	countStr := c.Query("count")

	var (
		start = uint64(0)
		count = uint64(10)
	)

	if startStr != "" {
		start, err = strconv.ParseUint(startStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusOK, utils.ParseResult(models.ErrParams, "", "start params error"))
			return
		}
	}

	if countStr != "" {
		count, err = strconv.ParseUint(countStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusOK, utils.ParseResult(models.ErrParams, "", "count params error"))
			return
		}
	}

	db := database.New()
	list, err := db.GetAllDomains(uint(start), uint(count))
	if err != nil {
		utils.Error("get all domains error:", err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "db get error"))
		return
	}

	totalNum, err := db.CountAllDomainNum()
	if err != nil {
		utils.Error("get total num error:", err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "db get error"))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccessWithData(map[string]interface{}{
		"total": totalNum,
		"list":  list,
	}))
	return
}

// CreateDomain 添加域名
func CreateDomain(c *gin.Context) {

	var p struct {
		Domain string `json:"domain"`
		Remark string `json:"remark"`
	}

	var err error

	domanKey := c.GetHeader(APIKeyName)

	err = json.NewDecoder(c.Request.Body).Decode(&p)
	if err != nil {
		utils.Error("json error:", err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrParams, "", "post body error"))
		return
	}

	fmt.Println(p.Domain)

	if !utils.IsOkDomain(p.Domain) {
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrParams, "", "domain error"))
		return
	}

	d := dns.New("fone")
	m, err := d.Create(domanKey, p.Domain, p.Remark)
	if err != nil {
		if err == dns.ErrDomainHasExist {
			c.JSON(http.StatusOK, utils.ParseResult(models.ErrDomainExist, "", "domain has exist"))
			return
		}
		utils.Error("dns create domain error:", err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "create dns error"))
		return
	}
	c.JSON(http.StatusOK, utils.ParseSuccessWithData(m))
	return
}

// DeleteDomain 删除域名
func DeleteDomain(c *gin.Context) {

	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrParams, "", "miss id"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrParams, "", "params id error"))
		return
	}

	db := database.New()
	m, err := db.GetDomainByID(id)
	if err != nil {
		utils.Errorf("get domain by id: %d error: %v", id, err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "get dns error"))
		return
	}

	d := dns.New("fone")
	err = d.DeleteDomain(id, m.FoneDomainID)
	if err != nil {
		utils.Errorf("delete domain by id: %d,domainid: %d, error: %v", id, m.FoneDomainID, err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "delete dns error"))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccess())
	return
}

// IsTakeOver 检查域名是否托管
//
// 存在 30 秒的缓存
// 如果托管状态改变，则数据库记录的状态也会改变
// 直接拿分配给这个域名的 key 进行查询就行，不一定需要 id
func IsTakeOver(c *gin.Context) {

	key := c.GetHeader(APIKeyName)
	if !strings.HasPrefix(key, dns.GetRecordKeyPrefix()) {
		// 说明不是recordKey
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var result struct {
		IsTakeOver int `json:"is_take_over"`
	}

	cacheKey := utils.GetCacheTakeOverKey(key)

	err := cache.DefaultCache.Get(cacheKey, &result)
	if err == nil {
		utils.Infof("get is_take_over from cache: %v", cacheKey)
		c.JSON(http.StatusOK, utils.ParseSuccessWithData(result))
		return
	}

	db := database.New()
	m, err := db.GetDomainByRecordKey(key)
	if err != nil {
		if err == scanner.ErrEmptyResult {
			c.JSON(http.StatusOK, utils.ParseResult(models.ErrDomainNotFound, "", ""))
			return
		}
		utils.Errorf("get domain by id: %s error: %v", key, err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "get dns error"))
		return
	}

	ns, err := sdk.GetDomainNS(m.Domain)
	if err != nil {
		utils.Error("get ns error:", err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", ""))
		return
	}

	isTakeOver := 0

	for _, a := range ns {
		for _, b := range *m.NameServer {
			if a == b {
				isTakeOver = 1
				break
			}
		}
	}

	where := map[string]interface{}{
		"id": m.ID,
	}

	update := make(map[string]interface{})

	if isTakeOver == 1 && m.IsTakeOver != 1 {
		update["is_take_over"] = 1
	} else if isTakeOver == 0 && m.IsTakeOver != 0 {
		update["is_take_over"] = 0
	}

	if len(update) > 0 {
		err = db.UpdateDomain(where, update)
		if err != nil {
			utils.Error("update domain error:", m.ID, update, err)
			c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "update domain error"))
			return
		}
	}

	result.IsTakeOver = isTakeOver

	go func() {
		cache.DefaultCache.Set(cacheKey, result, 30*time.Second)
	}()

	c.JSON(http.StatusOK, utils.ParseSuccessWithData(result))
	return
}
