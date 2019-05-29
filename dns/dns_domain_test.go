package dns

import (
	"fmt"
	"testing"
	"time"

	"github.com/solution9th/NSBridge/config"
	"github.com/solution9th/NSBridge/service/mysql"
	"github.com/solution9th/NSBridge/utils"
)

func InitDB(t *testing.T) {

	err := mysql.InitDefaultDB("dns_one", "root", "666666", "10.20.63.173", 3306)
	if err != nil {
		t.Error("init db error:", err)
	}

	err = config.InitConfig("config")

	// 初始化日志
	utils.NewLogFile("/Users/bi/log/dns.log", time.RFC3339)
}

func TestGetDefaultNS(t *testing.T) {
	InitDB(t)
	s := GetDefaultNS()
	fmt.Printf("%#v\n", s)
}

func TestCreate(t *testing.T) {

	InitDB(t)

	d := New("fone")

	m, err := d.Create("bipenghao", "d123dd2233.newio.cc", "这是一个测试🌽")
	if err != nil {
		t.Error("create domain error:", err)
		return
	}

	t.Logf("%#v\n", m)
	fmt.Println(m.NameServer)
}

func TestDelete(t *testing.T) {

	InitDB(t)

	d := New("fone")

	err := d.DeleteDomain(28, 67)
	if err != nil {
		t.Error("delete domain error:", err)
	}
}

func TestGetRecordKS(t *testing.T) {

	for i := 0; i < 10; i++ {
		s := GetRecordSecret()
		fmt.Println(s)
		k := GetRecordKey(s)
		fmt.Println(k)
	}
}
