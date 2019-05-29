package apiweb

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/solution9th/NSBridge/models"
	"github.com/solution9th/NSBridge/service/database"
	"github.com/solution9th/NSBridge/utils"
	"math"
	"net/http"
	"strconv"
)

var (
	authTable   = "auth"
	recordTable = "dns_record"
)

type Remark struct {
	ID     int    `json:"id,omitempty"`
	Remark string `json:"remark"`
}

// CreateAuth 创建授权
func CreateAuth(c *gin.Context) {
	/*
		{
		"remark":"这是cs域名"
		}
	*/
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrParams, err.Error(), nil))
		return
	}
	var remark Remark
	err = json.Unmarshal(data, &remark)
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrParamsFormat, err.Error(), nil))
		return
	}

	d := database.New()
	auth := map[string]interface{}{
		"domain_key":    utils.GenAuthKey(),
		"domain_secret": utils.GenAuthSecret(),
		"disable":       0,
		"remark":        remark.Remark,
	}
	_, err = d.InsertAuth(auth)
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrAuthInsert, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccess())
	return
}

// SearchAuthInfo 搜索及查询授权信息
func SearchAuthInfo(c *gin.Context) {
	domainKey := c.Query("domain_key")
	disable := c.Query("disable")
	offset := c.Query("offset")
	count := c.Query("count")
	disableInt, err := strconv.Atoi(disable)
	if disable != "" && err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrDisableParam, err.Error(), nil))
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if offset != "" && err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrOffsetParam, err.Error(), nil))
		return
	}
	countInt, err := strconv.Atoi(count)
	if count != "" && err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrCountParam, err.Error(), nil))
		return
	}

	if count == "" || countInt == 0 {
		countInt = math.MaxUint32
	}
	if offsetInt < 0 {
		offsetInt = 0
	}
	if disable == "" {
		disableInt = 2 // 获取所有的
	}

	d := database.New()
	m, err := d.SearchAllAuth(domainKey, disableInt, uint(offsetInt), uint(countInt))
	if err != nil {
		utils.Error(err.Error())
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrSrever, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccessWithData(m))
	return
}

// RemarkAuth 修改授权备注
func RemarkAuth(c *gin.Context) {
	/*
		{
		"id":2,
		"remark":"这是cs域名"
		}
	*/
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrParams, err.Error(), nil))
		return
	}
	var remark Remark
	err = json.Unmarshal(data, &remark)
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrParamsFormat, err.Error(), nil))
		return
	}
	d := database.New()
	where := map[string]interface{}{
		"id": remark.ID,
	}
	update := map[string]interface{}{
		"remark": remark.Remark,
	}
	err = d.UpdateAuth(where, update)
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrAuthUpdate, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccess())
	return
}

// DisableAuth 禁用授权
func DisableAuth(c *gin.Context) {
	id := c.Param("auth_id")
	idInt, err := strconv.Atoi(id)
	if id == "" || err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrParams, err.Error(), nil))
		return
	}

	d := database.New()

	exist, err := d.ExistAuthByAuthID(idInt)
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrSrever, err.Error(), nil))
		return
	}
	if !exist {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrAuthIDNotFound, "", nil))
		return
	}

	where := map[string]interface{}{
		"id": idInt,
	}
	err = d.DisableAuth(where)
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrAuthDisable, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccess())
	return
}

// DeleteAuth 删除授权
func DeleteAuth(c *gin.Context) {

	id := c.Param("auth_id")
	idInt, err := strconv.Atoi(id)
	if id == "" || err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrParams, err.Error(), nil))
		return
	}

	d := database.New()

	exist, err := d.ExistAuthByAuthID(idInt)
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrSrever, err.Error(), nil))
		return
	}
	if !exist {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrAuthIDNotFound, "", nil))
		return
	}

	err = d.DeleteAuthByID(idInt)
	if err != nil {
		if err == database.WebErrDelAuthNotDisable {
			c.JSON(http.StatusOK, utils.ParseResult(models.WebErrDelAuthNotDisable, "", nil))
			return
		}
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrAuthDelete, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.ParseSuccess())
	return
}
