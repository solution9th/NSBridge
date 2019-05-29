package apiweb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/solution9th/NSBridge/models"
	"github.com/solution9th/NSBridge/service/database"
	"github.com/solution9th/NSBridge/utils"
	"math"
	"net/http"
	"strconv"
)

// SearchRecords 根据domainID sub_domain record_type value offset count 获取记录列表
func SearchRecords(c *gin.Context) {
	domainID := c.Param("domain_id")
	domainIDInt, err := strconv.Atoi(domainID)
	if domainID != "" && err != nil || domainID == "" {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrDomainIDParam, err.Error(), nil))
		return
	}

	d := database.New()
	recordCount, err := d.GetRecordCountByDomainID(domainIDInt)
	if recordCount <= 0 {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrEmptyData, "", nil))
		return
	}

	recordType := c.Query("record_type")
	subDomainOrValue := c.Query("sub_or_val")
	offset := c.Query("offset")
	count := c.Query("count")
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

	m, err := d.SearchRecord(domainIDInt, recordType, subDomainOrValue, uint(offsetInt), uint(countInt))
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrDomainIDParam, err.Error(), nil))
		return
	}
	fmt.Println(utils.GenJson(m))
	c.JSON(http.StatusOK, utils.ParseSuccessWithData(m))
	return
}

// GetDomainTypes 根据domainID 获取 该domain 有哪些类型的record
func GetDomainTypes(c *gin.Context)  {
	domainID := c.Param("domain_id")
	domainIDInt, err := strconv.Atoi(domainID)
	if domainID != "" && err != nil || domainID == "" {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrDomainIDParam, err.Error(), nil))
		return
	}

	d := database.New()
	types, err := d.GetRecordTypesByDomainID(domainIDInt)
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrSrever, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.ParseSuccessWithData(types))
	return
}