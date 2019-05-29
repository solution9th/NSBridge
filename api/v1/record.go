package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/solution9th/NSBridge/dns"
	"github.com/solution9th/NSBridge/models"
	"github.com/solution9th/NSBridge/sdk"
	"github.com/solution9th/NSBridge/service/database"
	"github.com/solution9th/NSBridge/utils"

	"github.com/gin-gonic/gin"
	"github.com/haozibi/gendry/scanner"
)

// GetRecordList 获取解析记录列表
func GetRecordList(c *gin.Context) {

	var err error

	key := c.GetHeader(APIKeyName)
	if !strings.HasPrefix(key, dns.GetRecordKeyPrefix()) {
		// 说明不是 recordKey
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

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
	m, err := db.GetDomainByRecordKey(key)
	if err != nil {
		utils.Errorf("get domain by key: %s error: %v", key, err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "get dns error"))
		return
	}

	list, err := db.GetAllRecordByDomainID(m.ID, uint(start), uint(count))
	if err != nil {
		utils.Error("get all domains error:", err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "db get error"))
		return
	}

	totalNum, err := db.CountAllRecordNum(m.ID)
	if err != nil {
		utils.Error("get all domains num error:", err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "db get error"))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccessWithData(map[string]interface{}{
		"total": totalNum,
		"list":  list,
	}))
	return
}

// CreateRecord 新增记录
func CreateRecord(c *gin.Context) {

	key := c.GetHeader(APIKeyName)
	if !strings.HasPrefix(key, dns.GetRecordKeyPrefix()) {
		// 说明不是 recordKey
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	db := database.New()
	m, err := db.GetDomainByRecordKey(key)
	if err != nil {
		if err == scanner.ErrEmptyResult {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		utils.Errorf("get domain by key: %s error: %v", key, err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "get dns error"))
		return
	}

	p, resp := checkRecord(c)
	if resp.ErrCode != 0 {
		c.JSON(http.StatusOK, resp)
		return
	}

	p.DomainID = m.ID

	d := dns.New("fone")

	_, err = d.CreateRecord(p)
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "dns error"))
		return
	}
	c.JSON(http.StatusOK, utils.ParseSuccess())
	return
}

// UpdateRecord 更新记录
func UpdateRecord(c *gin.Context) {

	key := c.GetHeader(APIKeyName)
	if !strings.HasPrefix(key, dns.GetRecordKeyPrefix()) {
		// 说明不是 recordKey
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	idStr := c.Param("recordid")

	p, resp := checkRecord(c)
	if resp.ErrCode != 0 {
		c.JSON(http.StatusOK, resp)
		return
	}

	record, resp := getRecord(key, idStr)
	if resp.ErrCode != 0 {
		c.JSON(http.StatusOK, resp)
		return
	}

	p.DomainID = record.DomainID
	p.ID = record.ID

	d := dns.New("fone")

	err := d.UpdateRecord(p)
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", " update dns error"))
		return
	}
	c.JSON(http.StatusOK, utils.ParseSuccess())
	return
}

// GetRecord 获取解析记录详情
func GetRecord(c *gin.Context) {

	key := c.GetHeader(APIKeyName)
	if !strings.HasPrefix(key, dns.GetRecordKeyPrefix()) {
		// 说明不是 recordKey
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	idStr := c.Param("recordid")

	record, resp := getRecord(key, idStr)
	if resp.ErrCode != 0 {
		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccessWithData(record))
	return
}

// DeleteRecord 删除记录
func DeleteRecord(c *gin.Context) {

	key := c.GetHeader(APIKeyName)
	if !strings.HasPrefix(key, dns.GetRecordKeyPrefix()) {
		// 说明不是 recordKey
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	idStr := c.Param("recordid")

	record, resp := getRecord(key, idStr)
	if resp.ErrCode != 0 {
		c.JSON(http.StatusOK, resp)
		return
	}

	d := dns.New("fone")
	err := d.DeleteRecord(record.ID)
	if err != nil {
		utils.Errorf("delete record: %v error: %v", record.ID, err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "delete dns error"))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccess())
	return
}

// DisableRecord 启动暂停记录
func DisableRecord(c *gin.Context) {

	key := c.GetHeader(APIKeyName)
	if !strings.HasPrefix(key, dns.GetRecordKeyPrefix()) {
		// 说明不是 recordKey
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	idStr := c.Param("recordid")

	var p struct {
		Disabale bool `json:"disable"`
	}

	err := json.NewDecoder(c.Request.Body).Decode(&p)
	if err != nil {
		utils.Error("json error:", err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrParams, "", "post body error"))
		return
	}

	record, resp := getRecord(key, idStr)
	if resp.ErrCode != 0 {
		c.JSON(http.StatusOK, resp)
		return
	}

	if p.Disabale && record.Disable == 1 ||
		!p.Disabale && record.Disable == 0 {
		c.JSON(http.StatusOK, utils.ParseSuccess())
		return
	}

	d := dns.New("fone")
	err = d.DisableRecord(record.ID, p.Disabale)
	if err != nil {
		utils.Errorf("disable record: %v error: %v", record.ID, err)
		c.JSON(http.StatusOK, utils.ParseResult(models.ErrDB, "", "disable dns error"))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccess())
	return

}

// 检查参数
func checkRecord(c *gin.Context) (dns.RecordInfo, utils.Response) {

	var p dns.RecordInfo

	err := json.NewDecoder(c.Request.Body).Decode(&p)
	if err != nil {
		utils.Error("json error:", err)
		return p, utils.ParseResult(models.ErrParams, "", "post body error")
	}

	if !utils.IsExist(p.RecordType, sdk.AllowTypeList) {
		return p, utils.ParseResult(models.ErrParams, "", "record_type error")
	}

	if !utils.IsExist(p.Unit, sdk.AllowUnitList) {
		return p, utils.ParseResult(models.ErrParams, "", "unit param error")
	}

	if p.RecordType == sdk.RecordMXType && p.Priority <= 0 {
		// return utils.ParseResult(models.ErrParams, "", "mx need priority"))
		p.Priority = 5
	}

	if p.Priority > 100 {
		p.Priority = 100
	}

	for _, v := range strings.Split(p.SubDomain, ".") {
		if !check(v) {
			return p, utils.ParseResult(models.ErrParams, "", "sub_domain error")
		}
	}

	return p, utils.ParseSuccess()
}

func check(s string) bool {

	l := len([]rune(s))

	if l > 63 || l == 0 {
		return false
	}

	if l == 1 && (s == "*" || s == "@") {
		return true
	}

	if strings.HasPrefix(s, "-") || strings.HasSuffix(s, "-") {
		return false
	}

	re := `^[-a-zA-Z0-9]{0,62}$`

	return utils.Match(re, s)
}

func getRecord(key, idStr string) (database.MDNSRecord, utils.Response) {

	var record database.MDNSRecord

	if idStr == "" {
		return record, utils.ParseResult(models.ErrParams, "", "miss id")
	}

	recordID, err := strconv.Atoi(idStr)
	if err != nil || recordID <= 0 {
		return record, utils.ParseResult(models.ErrParams, "", "params id error")
	}

	db := database.New()
	m, err := db.GetDomainByRecordKey(key)
	if err != nil {
		if err == scanner.ErrEmptyResult {
			return record, utils.ParseResult(models.ErrPermissionDenied, "", "key not match")
		}
		utils.Errorf("get domain by key: %s error: %v", key, err)
		return record, utils.ParseResult(models.ErrDB, "", "get dns error")
	}

	record, err = db.GetRecordByRecordID(m.ID, recordID)
	if err != nil {
		if err == scanner.ErrEmptyResult {
			return record, utils.ParseResult(models.ErrRecordNotFound, "", "")
		}
		utils.Errorf("get record by id: %v,%v error: %v", m.ID, recordID, err)
		return record, utils.ParseResult(models.ErrDB, "", "get dns error")
	}

	return record, utils.ParseSuccess()
}
