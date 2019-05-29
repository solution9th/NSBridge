package apiweb

import (
	"github.com/gin-gonic/gin"
	"github.com/solution9th/NSBridge/models"
	"github.com/solution9th/NSBridge/service/database"
	"github.com/solution9th/NSBridge/utils"
	"math"
	"net/http"
	"strconv"
)

// GetDomainList 获取域名列表
func GetDomainList(c *gin.Context) {
	domain := c.Query("domain")
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
	d := database.New()
	m, err := d.SearchDomainsByDomain(domain, uint(offsetInt), uint(countInt))
	if err != nil {
		c.JSON(http.StatusOK, utils.ParseResult(models.WebErrSrever, err.Error(), nil))
		return
	}

	for i, v := range m {
		recordCount, err := d.GetRecordCountByDomainID(v.ID)
		if err != nil {
			utils.Error("GetRecordCountByDomainID Err: ", err.Error())
			c.JSON(http.StatusOK, utils.ParseResult(models.WebErrSrever, err.Error(), nil))
			return
		}
		m[i].RecordCount = recordCount
	}

	c.JSON(http.StatusOK, utils.ParseSuccessWithData(m))
}
