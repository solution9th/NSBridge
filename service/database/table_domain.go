package database

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/haozibi/gendry/builder"
	"github.com/haozibi/gendry/scanner"
)

type MDNSDomain struct {
	ID           int       `json:"id"`
	FoneDomainID int       `json:"fone_domain_id"`
	DomainKey    string    `json:"domain_key"`
	Domain       string    `json:"domain"`
	NameServer   *NS       `json:"name_server"`
	SOAEmail     string    `json:"soa_email"`
	Remark       string    `json:"remark"`
	IsTakeOver   int       `json:"is_take_over"` // 是否接管
	IsOpenKey    int       `json:"is_open_key"`  // key 是否可正常使用
	RecordKey    string    `json:"record_key"`   // 操作这个域名的key
	RecordSecret string    `json:"record_secret"`
	RecordCount  int       `json:"record_count"`
	CreateAt     time.Time `json:"create_at"`
	UpdateAt     time.Time `json:"update_at"`
}

type NS []string

func (d *NS) UnmarshalByte(data []byte) error {
	return json.Unmarshal(data, d)
}

func (t *Tables) query(tableName string, where map[string]interface{}, ptr interface{}) error {

	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("params error: query need ptr")
	}

	cond, vals, err := builder.BuildSelect(tableName, where, nil)
	if err != nil {
		return err
	}

	rows, err := t.DB.Query(cond, vals...)
	if err != nil {
		return err
	}

	err = scanner.ScanClose(rows, ptr)
	return err
}

func (t *Tables) count(cond string, val ...interface{}) (int, error) {

	rows, err := t.DB.Query(cond, val...)
	if err != nil {
		return 0, err
	}

	tmp := struct {
		Count int `json:"count"`
	}{}

	err = scanner.ScanClose(rows, &tmp)

	return tmp.Count, err
}

// GetDomainByID 根据 id 查找域名信息
func (t *Tables) GetDomainByID(id int) (m MDNSDomain, err error) {

	where := map[string]interface{}{
		"id": id,
	}

	err = t.query(tableDomain, where, &m)
	return
}

// GetDomainByDomainKey 根据 domain_key 查找域名
func (t *Tables) GetDomainByDomainKey(key string, start, count uint) (m []MDNSDomain, err error) {

	where := map[string]interface{}{
		"domain_key": key,
		"_limit":   []uint{start, count},
	}

	err = t.query(tableDomain, where, &m)
	return
}

// GetDomainByRecordKey 根据 key 查找域名
func (t *Tables) GetDomainByRecordKey(key string) (m MDNSDomain, err error) {

	where := map[string]interface{}{
		"record_key": key,
	}

	err = t.query(tableDomain, where, &m)
	return
}

// GetDomainByDomain 根据域名查找域名信息
func (t *Tables) GetDomainByDomain(domain string) (m MDNSDomain, err error) {

	where := map[string]interface{}{
		"domain":   domain,
		"_orderby": "create_at desc",
		"_limit":   []uint{1},
	}

	err = t.query(tableDomain, where, &m)
	return
}

// CountDomainNum 统计某个域名的个数
func (t *Tables) CountDomainNum(domain string) (int, error) {

	sql := fmt.Sprintf("SELECT COUNT(*) AS `count` FROM %v WHERE domain = {{domain}}", tableDomain)

	cond, val, err := builder.NamedQuery(sql, map[string]interface{}{
		"domain": domain,
	})
	if err != nil {
		return 0, err
	}

	return t.count(cond, val...)
}

// InsertDomain insert domain
func (t *Tables) InsertDomain(data map[string]interface{}) (id int, err error) {

	var datas []map[string]interface{}

	datas = append(datas, data)

	cond, vals, err := builder.BuildInsert(tableDomain, datas)
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

// UpdateDomain update domain
func (t *Tables) UpdateDomain(where, update map[string]interface{}) error {
	if where == nil {
		return fmt.Errorf("params error: where error")
	}

	cond, vals, err := builder.BuildUpdate(tableDomain, where, update)

	_, err = t.DB.Exec(cond, vals...)
	return err
}

// DeleteDomainByID 根据 id 删除域名
func (t *Tables) DeleteDomainByID(id int) error {

	if id <= 0 {
		return fmt.Errorf("params error: id error")
	}

	where := map[string]interface{}{
		"id": id,
	}

	cond, val, err := builder.BuildDelete(tableDomain, where)
	if err != nil {
		return err
	}

	_, err = t.DB.Exec(cond, val...)
	return err
}

// GetAllDomains 获得所有域名
func (t *Tables) GetAllDomains(start, count uint) (m []MDNSDomain, err error) {

	where := map[string]interface{}{
		"_orderby": "id desc",
		"_limit":   []uint{start, count},
	}

	err = t.query(tableDomain, where, &m)
	return
}

// GetAllDomains 获得所有域名
func (t *Tables) SearchDomainsByDomain(domain string, start, count uint) (m []MDNSDomain, err error) {

	where := map[string]interface{}{
		"domain like": "%" + domain + "%",
		"_orderby":    "id desc",
		"_limit":      []uint{start, count},
	}

	err = t.query(tableDomain, where, &m)
	return
}

// CountAllDomainNum 统计所有域名的个数
func (t *Tables) CountAllDomainNum() (int, error) {

	sql := fmt.Sprintf("SELECT COUNT(*) AS `count` FROM %v ", tableDomain)

	cond, val, err := builder.NamedQuery(sql, nil)
	if err != nil {
		return 0, err
	}

	return t.count(cond, val...)
}

// CountAllDomainNum 统计某个key下域名的个数
func (t *Tables) CountOwnDomainNum(domainKey string) (int, error) {

	sql := fmt.Sprintf("SELECT COUNT(*) AS `count` FROM %v WHERE domain_key={{domain_key}}", tableDomain)

	where := map[string]interface{}{
		"domain_key": domainKey,
	}

	cond, val, err := builder.NamedQuery(sql, where)
	if err != nil {
		return 0, err
	}

	return t.count(cond, val...)
}