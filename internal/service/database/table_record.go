package database

import (
	"fmt"
	"time"

	"github.com/solution9th/NSBridge/internal/sdk"
	"github.com/solution9th/NSBridge/internal/utils"

	"github.com/haozibi/gendry/builder"
	"github.com/haozibi/gendry/scanner"
)

var (
	recordTableName = "dns_record"
	domainTableName = "dns_domain"
)

type MDNSRecord struct {
	ID           int            `json:"id"`
	DomainID     int            `json:"domain_id"`
	FoneDomainID int            `json:"fone_domain_id"`
	FoneRecordID int            `json:"fone_record_id"`
	SubDomain    string         `json:"sub_domain"`
	RecordType   sdk.RecordType `json:"record_type"`
	Value        string         `json:"value"`
	LineID       int            `json:"line_id"`
	LineName     string         `json:"line_name"`
	TTL          int            `json:"ttl"`
	Unit         sdk.RecordUnit `json:"unit"`
	Priority     int            `json:"priority"`
	Disable      int            `json:"disable"`
	CreateAt     time.Time      `json:"create_at"`
	UpdateAt     time.Time      `json:"update_at"`
}

type RecordTypes struct {
	RecordType sdk.RecordType `json:"record_type"`
}

// GetFoneDomainIDByDomainID 根据domainID获取FoneDomainID
func (t *Tables) GetFoneDomainIDByDomainID(domainID int) (foneDomainID int, err error) {
	sql := fmt.Sprintf("SELECT fone_domain_id FROM %v WHERE id = {{id}}", domainTableName)
	cond, val, err := builder.NamedQuery(sql, map[string]interface{}{
		"id": domainID,
	})
	if err != nil {
		return 0, err
	}

	rows, err := t.DB.Query(cond, val...)
	if err != nil {
		return 0, err
	}

	tmp := struct {
		FoneDomainID int `json:"fone_domain_id"`
	}{}
	err = scanner.ScanClose(rows, &tmp)

	return tmp.FoneDomainID, err
}

// GetRecordByDomainID 根据域名id查找解析记录
func (t *Tables) GetAllRecordByDomainID(domainID int, start, count uint) (m []MDNSRecord, err error) {
	where := map[string]interface{}{
		"domain_id": domainID,
		"_orderby":  "create_at desc",
		"_limit":    []uint{start, count},
	}

	err = t.query(tableRecord, where, &m)
	return
}

// GetRecordTypesByDomainID 根据domainID 获取该domain记录的所有类型
func (t *Tables) GetRecordTypesByDomainID(domainID int) (r []RecordTypes, err error) {
	sql := fmt.Sprintf("SELECT DISTINCT(record_type) FROM %v WHERE domain_id= {{id}}", recordTableName)
	cond, val, err := builder.NamedQuery(sql, map[string]interface{}{
		"id": domainID,
	})
	if err != nil {
		return
	}

	rows, err := t.DB.Query(cond, val...)
	if err != nil {
		return
	}

	err = scanner.ScanClose(rows, &r)
	return
}

// SearchRecord 根据type subdomain value 查找解析记录
func (t *Tables) SearchRecord(domainID int, recordType, subDomainOrValue string, start, count uint) (m []MDNSRecord, err error) {

	sql := ""

	var where map[string]interface{}
	if recordType == "" { // 所有类型
		sql = "SELECT *  FROM " + tableRecord + " WHERE domain_id={{domain_id}} and (value like {{value}} or sub_domain like {{sub_domain}}) order by create_at desc limit {{start}}, {{count}};"
		where = map[string]interface{}{
			"domain_id":  domainID,
			"sub_domain": "%" + subDomainOrValue + "%",
			"value":      "%" + subDomainOrValue + "%",
			"start":      uint(start),
			"count":      uint(count),
		}
	} else {
		sql = "SELECT *  FROM " + tableRecord + " WHERE domain_id={{domain_id}} and record_type = {{record_type}} and (value like {{value}} or sub_domain like {{sub_domain}}) order by create_at desc limit {{start}}, {{count}};"
		where = map[string]interface{}{
			"domain_id":   domainID,
			"record_type": recordType,
			"sub_domain":  "%" + subDomainOrValue + "%",
			"value":       "%" + subDomainOrValue + "%",
			"start":       uint(start),
			"count":       uint(count),
		}
	}

	cond, val, err := builder.NamedQuery(sql, where)
	if err != nil {
		return
	}

	rows, err := t.DB.Query(cond, val...)
	if err != nil {
		return
	}

	err = scanner.ScanClose(rows, &m)
	if err != nil {
		return
	}

	lineMap := utils.GetLineMap()
	for i, v := range m {
		m[i].LineName = lineMap[v.LineID]
	}
	return
}

// CountAllRecordNum 统计某个域名中所有记录的个数
func (t *Tables) CountAllRecordNum(domainID int) (int, error) {

	sql := fmt.Sprintf("SELECT COUNT(*) AS `count` FROM %v  WHERE domain_id = {{domain_id}} ", tableRecord)

	cond, val, err := builder.NamedQuery(sql, map[string]interface{}{
		"domain_id": domainID,
	})
	if err != nil {
		return 0, err
	}

	return t.count(cond, val...)

}

// GetFoneDomainIDAndFoneRecordIDByRecordID 根据recordID获取foneDoaminID和foneRecordID
func (t *Tables) GetFoneDomainIDAndFoneRecordIDByRecordID(recordID int) (foneDomainID, foneRecordID int, err error) {
	sql := fmt.Sprintf("SELECT fone_domain_id, fone_record_id FROM %v WHERE id = {{id}}", recordTableName)
	cond, val, err := builder.NamedQuery(sql, map[string]interface{}{
		"id": recordID,
	})
	if err != nil {
		return 0, 0, err
	}

	rows, err := t.DB.Query(cond, val...)
	if err != nil {
		return 0, 0, err
	}

	tmp := struct {
		FoneDomainID int `json:"fone_domain_id"`
		FoneRecordID int `json:"fone_record_id"`
	}{}
	err = scanner.ScanClose(rows, &tmp)
	return tmp.FoneDomainID, tmp.FoneRecordID, err
}

// GetDomainCountByDomainID 根据doaminID获取
func (t *Tables) GetDomainCountByDomainID(domainID int) (int, error) {
	sql := fmt.Sprintf("SELECT COUNT(*) AS `count` FROM %v WHERE id = {{id}}", domainTableName)
	cond, val, err := builder.NamedQuery(sql, map[string]interface{}{
		"id": domainID,
	})
	if err != nil {
		return 0, err
	}

	return t.count(cond, val...)
}

// GetRecordCountByDomainID 统计某个域名的解析记录
func (t *Tables) GetRecordCountByDomainID(domainID int) (int, error) {
	sql := fmt.Sprintf("SELECT COUNT(*) AS `count` FROM %v WHERE domain_id = {{domain_id}}", recordTableName)
	cond, val, err := builder.NamedQuery(sql, map[string]interface{}{
		"domain_id": domainID,
	})
	if err != nil {
		return 0, err
	}

	return t.count(cond, val...)
}

// ExistRecordByRecordID 查看某个recordID是否存在
func (t *Tables) ExistRecordByRecordID(recordID int) (bool, error) {
	sql := fmt.Sprintf("SELECT COUNT(*) AS `count` FROM %v WHERE id = {{id}}", recordTableName)
	cond, val, err := builder.NamedQuery(sql, map[string]interface{}{
		"id": recordID,
	})
	if err != nil {
		return false, err
	}

	count, err := t.count(cond, val...)
	if count > 0 {
		return true, nil
	}

	return false, err
}

// InsertRecord 新增解析记录
func (t *Tables) InsertRecord(data map[string]interface{}) (id int, err error) {

	var datas []map[string]interface{}

	datas = append(datas, data)

	cond, vals, err := builder.BuildInsert(recordTableName, datas)
	r, err := t.DB.Exec(cond, vals...)
	if err != nil {
		return 0, err
	}

	ids, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(ids), nil
}

// UpdateRecord 更新解析记录
func (t *Tables) UpdateRecord(where, update map[string]interface{}) error {
	if where == nil {
		return fmt.Errorf("params error: where error")
	}

	cond, vals, err := builder.BuildUpdate(recordTableName, where, update)

	_, err = t.DB.Exec(cond, vals...)
	return err
}

// DeleteRecord 删除解析记录
func (t *Tables) DeleteRecord(where map[string]interface{}) error {
	if where == nil {
		return fmt.Errorf("params error: where error")
	}
	cond, val, err := builder.BuildDelete(recordTableName, where)

	_, err = t.DB.Exec(cond, val...)
	return err
}

// GetRecordByRecordID 根据记录 domainID id 获取记录所有信息
func (t *Tables) GetRecordByRecordID(domainID, recordID int) (m MDNSRecord, err error) {

	where := map[string]interface{}{
		"domain_id": domainID,
		"id":        recordID,
	}

	err = t.query(tableRecord, where, &m)
	return
}

// GetRecordByReID 根据记录 id 获取记录所有信息
func (t *Tables) GetRecordByReID(recordID int) (m MDNSRecord, err error) {

	where := map[string]interface{}{
		"id": recordID,
	}

	err = t.query(tableRecord, where, &m)
	return
}
