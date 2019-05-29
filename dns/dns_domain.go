package dns

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/solution9th/NSBridge/config"
	"github.com/solution9th/NSBridge/sdk"
	"github.com/solution9th/NSBridge/sdk/fone"
	"github.com/solution9th/NSBridge/service/database"
	"github.com/solution9th/NSBridge/utils"

	"github.com/haozibi/gendry/scanner"
)

var (
	ErrDomainHasExist = errors.New("domain exist")
	ErrDomainNotFound = errors.New("domain not found")
)

type DNS struct {
	engine sdk.DNSSDK
	db     *database.Tables
}

func New(engineName string) *DNS {

	r := new(DNS)

	switch engineName {
	case "fone":

		var (
			host    = config.FoneConfig.Host
			user    = config.FoneConfig.User
			passwd  = config.FoneConfig.Passwd
			timeout = config.FoneConfig.Timeout
		)

		r.engine, _ = fone.New(host, user, passwd, timeout)
	default:
		panic(errors.New("not found engine"))
	}

	r.db = database.New()

	return r
}

// Create 添加域名
//
// 添加流程
// 1. 检查域名是否存在
// 2. 如果不存在则创建事务
// 3. 在事务中添加域名记录
// 4. 调用 sdk 中添加域名的方法
// 5. 如果失败则回滚
// 6. 提交
func (d *DNS) Create(domainKey, domain, remark string) (database.MDNSDomain, error) {

	var m database.MDNSDomain

	if domain == "" {
		return m, fmt.Errorf("params error: domain is nil")
	}
	isFail := true

	count, err := d.db.CountDomainNum(domain)
	if err != nil {
		utils.Error("[dns] count domain num error:", err)
		return m, err
	}

	if count >= 1 {
		return m, ErrDomainHasExist
	}

	err = d.db.Begin()
	if err != nil {
		utils.Error("[dns] start ts error:", err)
		return m, err
	}

	defer func() {
		if isFail {
			err1 := d.db.Rollback()
			if err1 != nil {
				utils.Error("[dns] roll error:", err1)
			}
		}
	}()

	key := GetRecordKey(domain)
	secret := GetRecordSecret()

	data := map[string]interface{}{
		"domain_key":    domainKey,
		"domain":        domain,
		"name_server":   "",
		"soa_email":     "",
		"remark":        remark,
		"is_take_over":  false,
		"is_open_key":   true,
		"record_key":    key,
		"record_secret": secret,
	}

	id, err := d.db.InsertDomain(data)
	if err != nil {
		utils.Error("[dns] insert domain error:", err)
		return m, err
	}

	resp, err := d.engine.CreateDNSDomain(domain, GetDefaultNS(), "")
	if err != nil {
		if strings.Contains(err.Error(), "insert duplicate value") {
			return m, ErrDomainHasExist
		}
		utils.Error("[dns] sdk create domain error:", err)
		return m, err
	}

	where := map[string]interface{}{
		"id": id,
	}

	tmpBody, err := json.Marshal(resp.NS)
	if err != nil {
		utils.Error("[dns] json error:", err)
		return m, err
	}

	update := map[string]interface{}{
		"fone_domain_id": resp.ID,
		"name_server":    string(tmpBody),
	}

	err = d.db.UpdateDomain(where, update)
	if err != nil {
		utils.Error("[dns] update domain error:", err)
		return m, err
	}

	err = d.db.Commit()
	if err != nil {
		utils.Error("[dns] commit ts error:", err)
		return m, err
	}

	isFail = false

	return d.db.GetDomainByID(id)
}

// DeleteDomain 删除域名，查询操作有重复嫌疑
//
// 1. 查找域名
// 2. 开启事务
// 3. 数据库删除域名，删除记录
// 4. sdk 删除域名
// 5. 回滚或者提交
func (d *DNS) DeleteDomain(id, domainID int) error {

	m, err := d.db.GetDomainByID(id)
	if err != nil {
		if err == scanner.ErrEmptyResult {
			return ErrDomainNotFound
		}
		utils.Error("[dns] get domain error:", err)
		return err
	}

	if m.FoneDomainID != domainID {
		return fmt.Errorf("[dns] id and domainID not match")
	}

	isFail := true

	err = d.db.Begin()
	if err != nil {
		utils.Error("[dns] start ts error:", err)
		return err
	}

	defer func() {
		if isFail {
			err1 := d.db.Rollback()
			if err1 != nil {
				utils.Error("[dns] roll error:", err1)
			}
		}
	}()

	err = d.db.DeleteDomainByID(id)
	if err != nil {
		utils.Error("[dns] delete domain error:", err)
		return err
	}

	where := map[string]interface{}{
		"domain_id": id,
	}
	err = d.db.DeleteRecord(where)
	if err != nil {
		utils.Error("[dns] delete record error:", err)
		return err
	}

	err = d.engine.DeleteDNSDomain(domainID)
	if err != nil {
		utils.Error("[dns] sdk delete domain error:", err)
		return err
	}

	err = d.db.Commit()
	if err != nil {
		utils.Error("[dns] commit ts error:", err)
		return err
	}

	isFail = false

	return nil
}
