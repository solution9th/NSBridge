package dns

import (
	"fmt"
	"testing"
	"time"

	"github.com/solution9th/NSBridge/sdk"
	"github.com/solution9th/NSBridge/service/mysql"
	"github.com/solution9th/NSBridge/utils"
)

func initDB(t *testing.T) {

	err := mysql.InitDefaultDB("dns_one", "root", "pwd", "10.20.63.173", 3306)
	if err != nil {
		t.Error("init db error:", err)
	}

	// 初始化日志
	utils.NewLogFile("/Users/bi/log/dns.log", time.RFC3339)
}

func TestCreateRecord(t *testing.T) {

	initDB(t)

	d := New("fone")

	_, err := d.CreateRecord(RecordInfo{
		SubDomain:  "mafeng",
		DomainID:   5,
		Unit:       sdk.RecordHourUnit,
		TTL:        3,
		Value:      "11.23.34.222",
		RecordType: sdk.RecordCNAMEType,
		Priority:   3,
		LineID:     1,
	})
	if err != nil {
		t.Error("create record error:", err)
		return
	}

}

func TestGetRecordByDomainId(t *testing.T) {

	initDB(t)

	d := New("fone")

	m, err := d.GetRecordByDomainId(6)
	if err != nil {
		t.Error("create record error:", err)
		return
	}
	fmt.Println(utils.GenJson(m))
}

func TestUpdateRecord(t *testing.T) {

	initDB(t)

	d := New("fone")

	err := d.UpdateRecord(RecordInfo{
		ID:         13,
		DomainID:   6,
		SubDomain:  "mn1234g",
		Unit:       sdk.RecordMinUnit,
		TTL:        6,
		Value:      "mf.mm.com",
		RecordType: sdk.RecordMXType,
		Priority:   3,
		LineID:     2,
	})
	if err != nil {
		t.Error("create record error:", err)
		return
	}
}

func TestDisableRecord(t *testing.T) {
	initDB(t)

	d := New("fone")

	err := d.DisableRecord(6, false)
	if err != nil {
		t.Error("create record error:", err)
		return
	}
}

func TestDeleteRecord(t *testing.T) {
	initDB(t)

	d := New("fone")

	err := d.DeleteRecord(11)
	if err != nil {
		t.Error("create record error:", err)
		return
	}
}
