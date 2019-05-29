package database

import (
	"errors"
	"fmt"
	"github.com/haozibi/gendry/builder"
	"github.com/solution9th/NSBridge/utils"
	"time"
)

const authTable = "auth"

var (
	WebErrDelAuthNotDisable = errors.New("del not disable auth err")
)

type MAuth struct {
	ID           int       `json:"id"`
	DomainKey    string    `json:"domain_key"`
	DomainSecret string    `json:"domain_secret"`
	Remark       string    `json:"remark"`
	Disable      int       `json:"disable"`
	CreateAt     time.Time `json:"create_at"`
	UpdateAt     time.Time `json:"update_at"`
}

// ExistAuthByAuthID 查看某个authID是否存在
func (t *Tables) ExistAuthByAuthID(authID int) (bool, error) {
	sql := fmt.Sprintf("SELECT COUNT(*) AS `count` FROM %v WHERE id = {{id}}", authTable)
	cond, val, err := builder.NamedQuery(sql, map[string]interface{}{
		"id": authID,
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

// InsertAuth 新增Auth
func (t *Tables) InsertAuth(data map[string]interface{}) (id int, err error) {

	var datas []map[string]interface{}

	datas = append(datas, data)

	cond, vals, err := builder.BuildInsert(authTable, datas)
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

// GetAllAuth 获得所有域名
func (t *Tables) SearchAllAuth(domainKey string, disable int, start, count uint) (m []MAuth, err error) {

	where := map[string]interface{}{
		"domain_key like": "%" + domainKey + "%",
		"_orderby":        "id desc",
		"_limit":          []uint{start, count},
	}
	if disable != 2 {
		where = map[string]interface{}{
			"domain_key like": "%" + domainKey + "%",
			"disable":         disable,
			"_orderby":        "id desc",
			"_limit":          []uint{start, count},
		}
	}

	err = t.query(authTable, where, &m)
	return
}

// DisableAuth disable auth
func (t *Tables) DisableAuth(where map[string]interface{}) error {
	if where == nil {
		return fmt.Errorf("params error: where error")
	}

	sql := fmt.Sprintf("update %v set disable = 1 - disable where id = {{id}}", authTable)
	cond, vals, err := builder.NamedQuery(sql, where)

	_, err = t.DB.Exec(cond, vals...)
	return err
}

// UpdateAuth update auth
func (t *Tables) UpdateAuth(where, update map[string]interface{}) error {
	if where == nil {
		return fmt.Errorf("params error: where error")
	}

	cond, vals, err := builder.BuildUpdate(authTable, where, update)

	_, err = t.DB.Exec(cond, vals...)
	return err
}

// DeleteAuthByID 根据 id 删除auth
func (t *Tables) DeleteAuthByID(id int) (err error) {

	if id <= 0 {
		return fmt.Errorf("params error: id error")
	}

	where := map[string]interface{}{
		"id": id,
	}
	var m MAuth
	err = t.query(tableAuth, where, &m)
	if err != nil {
		utils.Error("DeleteAuthByID Err:", err.Error())
		return
	}
	if m.Disable == 0 {
		utils.Error("DeleteAuthByID Err: not disable of authId: ", id)
		err = WebErrDelAuthNotDisable
		return
	}
	cond, val, err := builder.BuildDelete(authTable, where)
	if err != nil {
		return
	}

	_, err = t.DB.Exec(cond, val...)
	return
}

// GetAuthByKey 根据 key 获得 auth
func (t *Tables) GetAuthByKey(domainKey string) (m MAuth, err error) {

	where := map[string]interface{}{
		"domain_key": domainKey,
		"_limit":     []uint{1},
	}

	err = t.query(tableAuth, where, &m)
	return
}
