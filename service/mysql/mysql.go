package mysql

import (
	"database/sql"
	"time"

	"github.com/solution9th/NSBridge/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/haozibi/gendry/manager"
	"github.com/haozibi/gendry/scanner"
)

// DefaultDB 默认 db，保存其他数据库信息
var DefaultDB *sql.DB

func init() {
	scanner.SetTagName("json")
}

// InitDefaultDB 初始化最基本的数据库
func InitDefaultDB(dbName, user, passwd, host string, port int) error {

	var err error

	DefaultDB, err = InitMySQL(dbName, user, passwd, host, port)
	if err != nil {
		utils.Errorf("init default error: %v", err)
		return err
	}
	return nil
}

// InitMySQL 初始化 mysql
func InitMySQL(dbName, user, passwd, host string, port int) (db *sql.DB, err error) {

	db, err = manager.New(dbName, user, passwd, host).Set(
		manager.SetCharset("utf8mb4"),
		manager.SetAllowCleartextPasswords(true),
		manager.SetInterpolateParams(true),
		manager.SetParseTime(true),
		manager.SetTimeout(1*time.Second),
		manager.SetReadTimeout(1*time.Second)).Port(port).Open(true)
	if err != nil {
		utils.Error("[mysql] init error:", err)
		return
	}

	db.SetConnMaxLifetime(60 * time.Second)

	utils.Infof("PING MySQL: %v OK", dbName)
	return db, nil
}
