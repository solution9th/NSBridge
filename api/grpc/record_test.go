package grpc

import (
	"fmt"
	"testing"
	"time"

	"github.com/solution9th/NSBridge/config"
	"github.com/solution9th/NSBridge/utils"

	pb "github.com/solution9th/NSBridge/dns_pb"
)

func MFInit(t *testing.T) {
	InitDB(t)
	config.InitConfig("config")
	utils.NewLogFile("/Users/bi/log/dns.log", time.RFC3339)
}

func TestCreateRecord(t *testing.T) {

	MFInit(t)
	c := New()

	r, err := c.RecordCreate(ctx, &pb.RequestRecordCreate{
		RecordKey:  "cmVjb3Jk2cab9d2edfd19539",
		SubDomain:  "uus2u.com",
		RecordType: "MX",
		Unit:       "day",
		Value:      "123443",
		LineId:     5,
	})
	if err != nil {
		t.Error("create domain error:", err)
	}

	fmt.Println(utils.GenJson(r))
}

func TestUpdateRecord(t *testing.T) {

	MFInit(t)
	c := New()

	r, err := c.RecordUpdate(ctx, &pb.RequestRecordUpdate{
		RecordKey:  "cmVjb3Jk2cab9d2edfd19539",
		SubDomain:  "uus2u.com",
		RecordType: "NS",
		Unit:       "day",
		Value:      "mafeng.com",
		LineId:     2,
		RecordId:   37,
	})
	if err != nil {
		t.Error("update record error:", err)
	}

	fmt.Println(utils.GenJson(r))
}

func TestListRecords(t *testing.T) {

	MFInit(t)
	c := New()

	r, err := c.RecordList(ctx, &pb.RequestRecordList{
		RecordKey: "cmVjb3Jk2cab9d2edfd19539",
	})
	if err != nil {
		t.Error("domain list error:", err)
	}

	fmt.Println(err, utils.GenJson(r))

}

func TestRecordInfo(t *testing.T) {

	MFInit(t)

	c := New()

	r, err := c.RecordInfo(ctx, &pb.RequestRecordInfo{
		RecordKey: "cmVjb3Jk2cab9d2edfd19539",
		RecordId:  13,
	})

	if err != nil {
		t.Error("domain status error:", err)
	}

	fmt.Println(err, utils.GenJson(r))
}

func TestRecordDisable(t *testing.T) {

	MFInit(t)
	c := New()

	r, err := c.RecordDisable(ctx, &pb.RequestRecordDisable{
		RecordKey: "cmVjb3Jk2cab9d2edfd19539",
		RecordId:  18,
		Disable:   false,
	})

	if err != nil {
		t.Error("domain status error:", err)
	}

	fmt.Println(err, utils.GenJson(r))
}

func TestRecordDelete(t *testing.T) {

	MFInit(t)
	c := New()

	r, err := c.RecordDelete(ctx, &pb.RequestRecordDelete{
		RecordKey: "cmVjb3Jk2cab9d2edfd19539",
		RecordId:  13,
	})

	if err != nil {
		t.Error("domain status error:", err)
	}

	fmt.Println(err, utils.GenJson(r))
}
