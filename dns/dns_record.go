package dns

import (
	"errors"
	"math"

	"github.com/solution9th/NSBridge/sdk"
	"github.com/solution9th/NSBridge/sdk/fone"
	"github.com/solution9th/NSBridge/service/database"
	"github.com/solution9th/NSBridge/utils"
)

/*

	解析记录相关的hander

*/
var (
	ErrDomainNotExist       = errors.New("domain is not exist")
	ErrRecordNotExist       = errors.New("record is not exist")
	ErrRecordTypeNotSupport = errors.New("record type is not support")
	ErrRecordTypeMutex      = errors.New("record type is mutex")
	ErrRecordExist          = errors.New("record is duplicate, can't create it again")
	ErrRecordMissTxt        = errors.New("type txt miss txt content")
	ErrRecordIP             = errors.New("error ip")
	ErrRecordDomain         = errors.New("error domain value")
	ErrRecordSubDomain      = errors.New("error sub_domain")
)

// 解析记录互斥表 详见: /docs/img/*.png
var mutexMap = map[sdk.RecordType]map[sdk.RecordType]bool{
	sdk.RecordCNAMEType: {
		sdk.RecordAType:     true,
		sdk.RecordCNAMEType: true,
		sdk.RecordMXType:    true,
		sdk.RecordTXTType:   true,
		sdk.RecordNSType:    true,
	},
	sdk.RecordAType: {
		sdk.RecordAType:     false,
		sdk.RecordCNAMEType: true,
		sdk.RecordMXType:    false,
		sdk.RecordTXTType:   false,
		sdk.RecordNSType:    true,
	},
	sdk.RecordTXTType: {
		sdk.RecordAType:     false,
		sdk.RecordCNAMEType: true,
		sdk.RecordMXType:    false,
		sdk.RecordTXTType:   false,
		sdk.RecordNSType:    true,
	},
	sdk.RecordMXType: {
		sdk.RecordAType:     false,
		sdk.RecordCNAMEType: true,
		sdk.RecordMXType:    false,
		sdk.RecordTXTType:   false,
		sdk.RecordNSType:    true,
	},
	sdk.RecordNSType: {
		sdk.RecordAType:     true,
		sdk.RecordCNAMEType: true,
		sdk.RecordMXType:    true,
		sdk.RecordTXTType:   true,
		sdk.RecordNSType:    false,
	},
}

type RecordInfo struct {
	ID         int            `json:"id,omitempty"`
	DomainID   int            `json:"domain_id"`
	SubDomain  string         `json:"sub_domain"`
	RecordType sdk.RecordType `json:"record_type"`
	Value      string         `json:"value"`
	LineID     int            `json:"line_id"`
	TTL        int            `json:"ttl"`
	Unit       sdk.RecordUnit `json:"unit"`
	Priority   int            `json:"priority"`
}

// CheckRecord
// 1. 校验值是否合法
// 2. 校验互斥
// 3. 校验是否重复
// 不能相同的记录值:
// (1) CNAME 相同主机同一线路下CNAME和CNAME记录不能共存
// (1) 同一记录 && 相同sub_domain && 同一线路 时:  value 不能相同
func (d *DNS) CheckRecord(record RecordInfo) error {
	if record.SubDomain == "" || len(record.SubDomain) > 63 {
		return ErrRecordSubDomain
	}

	// 校验 value
	switch record.RecordType {
	case sdk.RecordAType:
		if !fone.IsOkIP(record.Value) {
			return ErrRecordIP
		}
	case sdk.RecordCNAMEType:
		if !utils.IsOkDomain(record.Value) {
			return ErrRecordDomain
		}
	case sdk.RecordMXType:
		if !utils.IsOkDomain(record.Value) {
			return ErrRecordDomain
		}
	case sdk.RecordTXTType:
		if record.Value == "" {
			return ErrRecordMissTxt
		}
	case sdk.RecordNSType:
		if !utils.IsOkDomain(record.Value) {
			return ErrRecordDomain
		}
	default:
		return ErrRecordTypeNotSupport

	}

	// 互斥
	m, err := d.db.GetAllRecordByDomainID(record.DomainID, 0, math.MaxUint64)
	// utils.Error(utils.GenJson(m))
	if err != nil {
		return err
	}

	for _, item := range m {
		if record.ID == item.ID {
			continue
		}
		if record.SubDomain == item.SubDomain {
			if record.RecordType == item.RecordType &&
				record.Value == item.Value &&
				record.LineID == item.LineID {
				return ErrRecordExist
			}
			if mutexMap[record.RecordType][item.RecordType] {
				return ErrRecordTypeMutex
			}
		}
	}
	return nil
}

//
// 添加流程
// 1. 检查域名是否存在
// 2. 检查记录是否存在
// 3. 校验value是否合法
// 4. 查看记录是否互斥(参考互斥表)
// 5. 如果不存在则创建事务
// 6. 在事务中添加记录
// 7. 调用 sdk 中添加记录的方法
// 8. 如果失败则回滚
// 9. 提交
//
func (d *DNS) CreateRecord(record RecordInfo) (database.MDNSRecord, error) {

	var recordCreate database.MDNSRecord
	count, err := d.db.GetDomainCountByDomainID(record.DomainID)
	if err != nil {
		utils.Error("[dns] count record num error:", err, utils.GenJson(record))
		return recordCreate, err
	}

	if count <= 0 {
		return recordCreate, ErrDomainNotExist
	}

	//
	err = d.CheckRecord(record)
	if err != nil {
		return recordCreate, err
	}

	isFail := true
	err = d.db.Begin()
	if err != nil {
		utils.Error("[dns] start ts error:", err, utils.GenJson(record))
		return recordCreate, err
	}

	defer func() {
		if isFail {
			err1 := d.db.Rollback()
			if err1 != nil {
				utils.Error("[dns] insert record rollback error:", err1)
			}
		}
	}()

	recordMap := map[string]interface{}{
		"domain_id":   record.DomainID,
		"sub_domain":  record.SubDomain,
		"unit":        record.Unit,
		"ttl":         record.TTL,
		"priority":    record.Priority,
		"line_id":     record.LineID,
		"record_type": record.RecordType,
		"value":       record.Value,
	}
	id, err := d.db.InsertRecord(recordMap)
	if err != nil {
		utils.Error("[dns] insert record error:", err)
		return recordCreate, err
	}

	foneDomainID, err := d.db.GetFoneDomainIDByDomainID(record.DomainID)
	if err != nil {
		return recordCreate, err
	}

	recordID, err := d.engine.CreateDNSRecord(foneDomainID, record.SubDomain, record.RecordType, record.Value, record.LineID, record.Priority, record.TTL, record.Unit)
	if err != nil {
		utils.Error("[dns] sdk insert record error:", err, utils.GenJson(record))
		return recordCreate, err
	}

	where := map[string]interface{}{
		"id": id,
	}

	update := map[string]interface{}{
		"fone_domain_id": foneDomainID,
		"fone_record_id": recordID,
	}

	err = d.db.UpdateRecord(where, update)
	if err != nil {
		utils.Error("[dns] update record error:", err, utils.GenJson(record))
		return recordCreate, err
	}

	err = d.db.Commit()
	if err != nil {
		utils.Error("[dns] commit ts error:", err, utils.GenJson(record))
		return recordCreate, err
	}

	isFail = false

	recordCreate, err = d.db.GetRecordByReID(id)
	if err != nil {
		utils.Error("[dns] get record create error:", err, utils.GenJson(record))
		return recordCreate, err
	}
	return recordCreate, nil
}

func (d *DNS) GetRecordByDomainId(domainID int) ([]database.MDNSRecord, error) {
	var m []database.MDNSRecord
	count, err := d.db.GetDomainCountByDomainID(domainID)
	if err != nil {
		utils.Error("[dns] count record num error:", err)
		return m, err
	}

	if count <= 0 {
		return m, ErrDomainNotExist
	}

	m, err = d.db.GetAllRecordByDomainID(domainID, 0, math.MaxUint64)
	if err != nil {
		return m, err
	}
	return m, err
}

// update 更新记录
//
// 更新流程
// 1. 检查记录是否存在
// 2. 如果存在则创建事务
// 3. 拿到fonedomainid 和fonerecordid
// 4. 调用数据库更新
// 5. 调用 sdk 更新记录的方法
// 6. 如果失败则回滚
// 7. 提交
func (d *DNS) UpdateRecord(record RecordInfo) error {

	exist, err := d.db.ExistRecordByRecordID(record.ID)
	if err != nil {
		utils.Error("[dns] update record err", err)
		return err
	}
	if !exist {
		return ErrRecordNotExist
	}

	// 校验
	err = d.CheckRecord(record)
	if err != nil {
		return err
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
				utils.Error("[dns] update record rollback error:", err1)
			}
		}
	}()

	foneDmainID, foneRecordID, err := d.db.GetFoneDomainIDAndFoneRecordIDByRecordID(record.ID)
	if err != nil {
		return err
	}

	where := map[string]interface{}{
		"id": record.ID,
	}
	recordMap := map[string]interface{}{
		"sub_domain":  record.SubDomain,
		"unit":        record.Unit,
		"ttl":         record.TTL,
		"priority":    record.Priority,
		"line_id":     record.LineID,
		"record_type": record.RecordType,
		"value":       record.Value,
	}
	err = d.db.UpdateRecord(where, recordMap)
	if err != nil {
		utils.Error("[dns] update record error:", err, utils.GenJson(record))
		return err
	}

	// var ip, domain, txt string
	// switch record.RecordType {
	// case sdk.RecordAType:
	// 	ip = record.Value
	// case sdk.RecordCNAMEType:
	// 	domain = record.Value
	// case sdk.RecordMXType:
	// 	domain = record.Value
	// case sdk.RecordTXTType:
	// 	txt = record.Value
	// case sdk.RecordNSType:
	// 	domain = record.Value
	// default:
	// 	return ErrRecordTypeNotSupport
	// }

	err = d.engine.UpdateDNSRecord(foneDmainID, foneRecordID, record.SubDomain, record.RecordType, record.Value, record.LineID, record.Priority, record.TTL, record.Unit)
	if err != nil {
		return err
	}

	err = d.db.Commit()
	if err != nil {
		utils.Error("[dns] commit ts error:", err, utils.GenJson(record))
		return err
	}

	isFail = false
	return err
}

// Delete disable记录
//
// 更新流程
// 1. 检查记录是否存在
// 2. 如果存在则创建事务
// 3. 拿到fonedomainid 和fonerecordid
// 4. 调用数据库修改
// 5. 调用 sdk 中disable记录的方法
// 6. 如果失败则回滚
// 7. 提交
func (d *DNS) DisableRecord(recordID int, disable bool) (err error) {
	exist, err := d.db.ExistRecordByRecordID(recordID)
	if err != nil {
		utils.Error("[dns] disable record err", err)
		return
	}
	if !exist {
		return ErrRecordNotExist
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
				utils.Error("[dns] disable record rollback error:", err1)
			}
		}
	}()

	foneDomainID, foneRecordID, err := d.db.GetFoneDomainIDAndFoneRecordIDByRecordID(recordID)
	if err != nil {
		return err
	}

	where := map[string]interface{}{
		"id": recordID,
	}
	recordMap := map[string]interface{}{
		"disable": disable,
	}

	err = d.db.UpdateRecord(where, recordMap)
	if err != nil {
		return err
	}

	err = d.engine.DisableDNSRecord(foneDomainID, foneRecordID, disable)
	if err != nil {
		return err
	}

	err = d.db.Commit()
	if err != nil {
		utils.Error("[dns] commit ts error:", err)
		return err
	}

	isFail = false
	return
}

// Delete 删除记录
//
// 删除流程
// 1. 检查记录是否存在
// 2. 如果存在则创建事务
// 3. 拿到fonedomainid 和fonerecordid
// 4. 调用数据库删除
// 5. 调用 sdk 中删除记录的方法
// 6. 如果失败则回滚
// 7. 提交
func (d *DNS) DeleteRecord(recordID int) error {

	exist, err := d.db.ExistRecordByRecordID(recordID)
	if err != nil {
		utils.Error("[dns] delete record err", err)
		return err
	}
	if !exist {
		return ErrRecordNotExist
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
				utils.Error("[dns] update rollback error:", err1)
			}
		}
	}()

	foneDoaminID, foneRecordID, err := d.db.GetFoneDomainIDAndFoneRecordIDByRecordID(recordID)
	if err != nil {
		return err
	}

	where := map[string]interface{}{
		"id": recordID,
	}
	err = d.db.DeleteRecord(where)
	if err != nil {
		return err
	}

	err = d.engine.DeleteDNSRecord(foneDoaminID, foneRecordID)
	if err != nil {
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
