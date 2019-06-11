package fone

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/solution9th/NSBridge/sdk"
)

var (
	isDelete = false
)

func InitDNS(t *testing.T) sdk.DNSSDK {
	f, _ := New("https://admin.f1.com", "user", "pwd", 3)

	err := f.Login()
	if err != nil {
		panic(err)
	}

	return f
}

// 一个创建，查看，删除 流程
// 测试环境保证不要删除 11 号域名
func TestDomainFlow(t *testing.T) {

	f := InitDNS(t)

	testDomain := "test.dev.newio.cc"

	domain := fmt.Sprintf("%s.%s", time.Now().Format("20060102150405"), testDomain)
	domainID := 0

	auths := make([]sdk.Authority, 0)
	auth := sdk.Authority{
		Domain: "ns1.newio.cc",
	}
	auths = append(auths, auth)

	resp, err := f.CreateDNSDomain(domain, auths, "test@uu.com")
	if err != nil {
		t.Error("[test] create domain error:", err)
	}

	domainID = resp.ID

	info, err := f.GetDNSDomain(domainID)
	if err != nil {
		t.Error("[test] get domain error:", err)
	}

	if info.ID != domainID {
		t.Error("[test] domain not match")
	}

	list, err := f.GetAllDNSDomain()
	if err != nil {
		t.Error("[test] get all domain error:", err)
	}

	fmt.Println("total:", len(list.Data))

	isFound := false

	for _, v := range list.Data {
		if v.ID == domainID {
			isFound = true
		}
	}

	if !isFound {
		t.Error("[test] all domain error: not found")
	}

	domain = "update-" + domain

	err = f.UpdateDNSDomain(domainID, domain, auths, "")
	if err != nil {
		t.Error("[test] update error:", err)
	}

	// 测试记录的操作
	testReocrdFlow(domainID, t)

	if isDelete {
		err = f.DeleteDNSDomain(domainID)
		if err != nil {
			t.Error("[test] delete domain error:", err)
		}
	}
}

func TestErrorCreateDomain(t *testing.T) {
	f := InitDNS(t)

	tests := []struct {
		domainid int
		domain   string
		auth     []sdk.Authority
		want     string
	}{
		{-1, "", nil, "domain error"},
		{-1, "nnn.com.-", nil, "domain error"},
		{-1, "nnn.com", []sdk.Authority{
			{Domain: "123", TTL: 5, Unit: sdk.RecordDayUnit},
		}, "name server"},
		{-1, "dev.newio.cc", []sdk.Authority{
			{Domain: "123.com", TTL: 5, Unit: sdk.RecordDayUnit},
		}, "error"},
		{99, "1111", nil, "domain error"},
		{99, "1111.com", nil, "miss authority"},
		{99, "nnn.com", []sdk.Authority{
			{Domain: "123", TTL: 5, Unit: sdk.RecordDayUnit},
		}, "name server"},
		{-99, "1111", nil, "dnsid error"},
	}

	for k, v := range tests {
		if v.domainid == -1 {
			_, err := f.CreateDNSDomain(v.domain, v.auth, "")
			if err == nil {
				t.Error("[test] error not nil")
			}

			if !strings.Contains(err.Error(), v.want) {
				t.Errorf("k %v want: %v got: %v", k, v.want, err)
			}
			continue
		}
		err := f.UpdateDNSDomain(v.domainid, v.domain, v.auth, "")
		if err == nil {
			t.Error("[test] error not nil")
		}

		if !strings.Contains(err.Error(), v.want) {
			t.Errorf("k %v want: %v got: %v", k, v.want, err)
		}
	}
}

func TestErrorDomainGet(t *testing.T) {

	id := 0

	f := InitDNS(t)

	_, err := f.GetDNSDomain(id)
	if err == nil {
		t.Error("[test] get domain nil")
	}

	if err != ErrDNSIDParams {
		t.Errorf("want: %v, got: %v", ErrDNSIDParams, err)
	}

	err = f.DeleteDNSDomain(id)
	if err == nil {
		t.Error("[test] delete domain nil")
	}

	if err != ErrDNSIDParams {
		t.Errorf("want: %v, got: %v", ErrDNSIDParams, err)
	}
}

func TestGetAllDNSDomain(t *testing.T) {

	f := InitDNS(t)

	infos, err := f.GetAllDNSDomain()

	if err != nil {
		panic(err)
	}

	for k, v := range infos.Data {
		fmt.Println(k, v.Zone)
	}

	fmt.Println("count:", infos.Meta.Count)
}

func TestGetAllDNSRecords(t *testing.T) {

	f := InitDNS(t)

	infos, err := f.GetAllDNSRecords(16)

	if err != nil {
		panic(err)
	}

	for k, v := range infos.Data {
		fmt.Println(k, v.SubDomain)
	}

	fmt.Println("count:", infos.Meta.Count)
}

func TestFOneNew(t *testing.T) {

	_, err := New("", "", "", 5)
	if err == nil {
		t.Error("[test] new fone not nil")
	}

	if err != ErrMissInfoParams {
		t.Error("[test] error:", err)
	}

}

func TestRecordFlow(t *testing.T) {
	dnsid := 58
	testReocrdFlow(dnsid, t)
}

func testReocrdFlow(dnsid int, t *testing.T) {

	f := InitDNS(t)

	subDomain := fmt.Sprintf("%s", time.Now().Format("20060102150405"))

	recordid, err := f.CreateDNSRecord(dnsid, subDomain, sdk.RecordNSType, "ns1.newio.cc", 0, 1, 5, sdk.RecordSecUnit)

	if err != nil {
		t.Error("[test] create record error:", err)
	}

	subDomain = "update-" + subDomain

	err = f.UpdateDNSRecord(dnsid, recordid, subDomain, sdk.RecordAType, "1.1.1.1", 0, 1, 5, sdk.RecordSecUnit)
	if err != nil {
		t.Error("[test] update record error:", err)
	}

	err = f.DisableDNSRecord(dnsid, recordid, true)
	if err != nil {
		t.Error("[test] disable record error:", err)
	}

	info, err := f.GetDNSRecord(dnsid, recordid)
	if err != nil {
		t.Error("[test] get record error:", err)
	}

	fmt.Println(info.Disabled)

	list, err := f.GetAllDNSRecords(dnsid)
	if err != nil {
		t.Error("[test] get all domain error:", err)
	}

	fmt.Println("total:", len(list.Data))

	isFound := false

	for _, v := range list.Data {
		if v.ID == recordid {
			isFound = true
		}
	}

	if !isFound {
		t.Error("[test] all record error: not found")
	}

	if isDelete {
		err = f.DeleteDNSRecord(dnsid, recordid)
		if err != nil {
			t.Error("[test] delete record error:", err)
		}
	}

	fmt.Println(recordid)
}
