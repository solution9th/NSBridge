package grpc

import (
	"context"
	"strconv"
	"time"

	pb "github.com/solution9th/NSBridge/dns_pb"
	"github.com/solution9th/NSBridge/internal/dns"
	"github.com/solution9th/NSBridge/internal/nserr"
	"github.com/solution9th/NSBridge/internal/service/cache"
	"github.com/solution9th/NSBridge/internal/service/database"
	"github.com/solution9th/NSBridge/internal/utils"
)

const (
	DefaultUnit   = "min"
	DefaultLineID = 0
	DefaultTTL    = 5
)

// RecordList 获取解析记录列表
func (c *RPCServer) RecordList(ctx context.Context, in *pb.RequestRecordList) (*pb.ResponseRecordList, error) {

	key, ok := IsOkRequest(ctx, in.RecordKey)
	if !ok || !IsRecordKey(key) {
		return &pb.ResponseRecordList{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}

	if in.Count == 0 {
		in.Count = 10
	}

	db := database.New()

	domain, err := db.GetDomainByRecordKey(key)
	if err != nil {
		utils.Error("get all record error:", err)
		return &pb.ResponseRecordList{
			ErrCode: nserr.ErrDB,
			ErrMsg:  utils.ErrCode[nserr.ErrDB],
		}, nil
	}

	list, err := db.GetAllRecordByDomainID(domain.ID, uint(in.Start), uint(in.Count))
	if err != nil {
		utils.Error("get all record error:", err)
		return &pb.ResponseRecordList{
			ErrCode: nserr.ErrDB,
			ErrMsg:  utils.ErrCode[nserr.ErrDB],
		}, nil
	}

	totalNum, err := db.CountAllRecordNum(domain.ID)
	if err != nil {
		utils.Error("get total num error:", err)
		return &pb.ResponseRecordList{
			ErrCode: nserr.ErrDB,
			ErrMsg:  utils.ErrCode[nserr.ErrDB],
		}, nil

	}

	info := &pb.ResponseRecordList{}
	info.ErrCode = 0
	info.ErrMsg = "success"
	info.Data = &pb.RecordListData{}
	datas := make([]*pb.Record, 0)

	for _, v := range list {
		data := &pb.Record{}
		err = TypeConvert(v, &data)
		if err != nil {
			continue
		}
		datas = append(datas, data)
	}

	info.Data.List = datas
	info.Data.Total = int32(totalNum)

	return info, nil

}

// RecordDomainOfRK 获取record_key的域名
func (c *RPCServer) RecordDomainOfRK(ctx context.Context, in *pb.RequestDomainOfRK) (*pb.ResponseDomainOfRK, error) {
	key, ok := IsOkRequest(ctx, in.RecordKey)
	if !ok || !IsRecordKey(key) {
		return &pb.ResponseDomainOfRK{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}

	db := database.New()
	m, err := db.GetDomainByRecordKey(key)
	if err != nil {
		utils.Error("get record info error:", err)
		return &pb.ResponseDomainOfRK{
			ErrCode: nserr.ErrDB,
			ErrMsg:  utils.ErrCode[nserr.ErrDB],
		}, nil
	}
	info := &pb.ResponseDomainOfRK{
		ErrCode: 0,
		ErrMsg:  "success",
	}

	data := &pb.DomainInfo{}

	err = TypeConvert(m, data)
	if err != nil {
		utils.Error("type error:", err)
		return &pb.ResponseDomainOfRK{
			ErrCode: nserr.ErrDNSSDK,
			ErrMsg:  utils.ErrCode[nserr.ErrDNSSDK],
		}, nil
	}

	info.Data = data

	utils.Infof("[create domain] key: %v,info: %#v", key, data)

	return info, nil
}

// RecordInfo 获取解析记录详情
func (c *RPCServer) RecordInfo(ctx context.Context, in *pb.RequestRecordInfo) (*pb.ResponseRecordInfo, error) {

	key, ok := IsOkRequest(ctx, in.RecordKey)
	if !ok || !IsRecordKey(key) {
		return &pb.ResponseRecordInfo{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}
	if !CheckRecordPermission(key, int(in.RecordId)) {
		return &pb.ResponseRecordInfo{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}

	db := database.New()
	exist, err := db.ExistRecordByRecordID(int(in.RecordId))
	if err != nil {
		utils.Error("exist record err", err)
	}
	if !exist {
		return &pb.ResponseRecordInfo{
			ErrCode: nserr.ErrRecordNotFound,
			ErrMsg:  utils.ErrCode[nserr.ErrRecordNotFound],
		}, nil
	}
	record, err := db.GetRecordByReID(int(in.RecordId))
	if err != nil {
		utils.Error("get record info error:", err)
		return &pb.ResponseRecordInfo{
			ErrCode: nserr.ErrDB,
			ErrMsg:  utils.ErrCode[nserr.ErrDB],
		}, nil
	}
	info := &pb.ResponseRecordInfo{}
	info.ErrCode = 0
	info.ErrMsg = "success"
	data := pb.Record{}
	err = TypeConvert(record, &data)
	if err != nil {
		utils.Error("get record info error:", err)
		return &pb.ResponseRecordInfo{
			ErrCode: nserr.ErrErr,
			ErrMsg:  utils.ErrCode[nserr.ErrErr],
		}, nil
	}
	info.Data = &data
	return info, nil
}

// RecordCreate 新增解析记录
func (c *RPCServer) RecordCreate(ctx context.Context, in *pb.RequestRecordCreate) (*pb.ResponseRecordCreate, error) {

	key, ok := IsOkRequest(ctx, in.RecordKey)
	if !ok || !IsRecordKey(key) {
		return &pb.ResponseRecordCreate{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}

	if in.SubDomain == "" {
		return &pb.ResponseRecordCreate{
			ErrCode: nserr.ErrParams,
		}, nil
	}

	if in.RecordType == "" {
		return &pb.ResponseRecordCreate{
			ErrCode: nserr.ErrParams,
			ErrMsg:  utils.ErrCode[nserr.ErrParams] + ": miss record_type",
		}, nil
	}

	if in.Value == "" {
		return &pb.ResponseRecordCreate{
			ErrCode: nserr.ErrParams,
			ErrMsg:  utils.ErrCode[nserr.ErrParams] + ": miss value",
		}, nil
	}

	// 如果unit或ttl为空, 给其设置默认值
	if in.Unit == "" {
		in.Unit = DefaultUnit
	}

	if in.Ttl == 0 {
		in.Ttl = DefaultTTL
	}

	db := database.New()
	domain, err := db.GetDomainByRecordKey(key)
	if err != nil {
		utils.Error("get record info error:", err)
		return &pb.ResponseRecordCreate{
			ErrCode: nserr.ErrDB,
			ErrMsg:  utils.ErrCode[nserr.ErrDB],
		}, nil
	}

	d := dns.New("fone")
	var record dns.RecordInfo
	err = TypeConvert(in, &record)
	if err != nil {
		utils.Error("create record type convert error:", err)
		return &pb.ResponseRecordCreate{
			ErrCode: nserr.ErrErr,
			ErrMsg:  utils.ErrCode[nserr.ErrErr],
		}, nil
	}
	record.DomainID = domain.ID
	newRecord, err := d.CreateRecord(record)
	if err != nil {
		utils.Error("create dns error:", err)
		switch err {
		case dns.ErrRecordSubDomain:
			return &pb.ResponseRecordCreate{
				ErrCode: nserr.ErrRecordSubDomain,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordSubDomain],
			}, nil
		case dns.ErrRecordIP:
			return &pb.ResponseRecordCreate{
				ErrCode: nserr.ErrRecordIP,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordIP],
			}, nil
		case dns.ErrRecordDomain:
			return &pb.ResponseRecordCreate{
				ErrCode: nserr.ErrRecordDomain,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordDomain],
			}, nil
		case dns.ErrRecordMissTxt:
			return &pb.ResponseRecordCreate{
				ErrCode: nserr.ErrRecordMissTxt,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordMissTxt],
			}, nil
		case dns.ErrDomainNotExist:
			return &pb.ResponseRecordCreate{
				ErrCode: nserr.ErrDomainExist,
				ErrMsg:  utils.ErrCode[nserr.ErrDomainExist],
			}, nil
		case dns.ErrRecordTypeNotSupport:
			return &pb.ResponseRecordCreate{
				ErrCode: nserr.ErrRecordTypeSupport,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordTypeSupport],
			}, nil
		case dns.ErrRecordTypeMutex: // 互斥
			return &pb.ResponseRecordCreate{
				ErrCode: nserr.ErrRecordTypeMutex,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordTypeMutex],
			}, nil
		case dns.ErrRecordExist:
			return &pb.ResponseRecordCreate{
				ErrCode: nserr.ErrRecordExist,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordExist],
			}, nil
		default:
			return &pb.ResponseRecordCreate{
				ErrCode: nserr.ErrDNSSDK,
				ErrMsg:  utils.ErrCode[nserr.ErrDNSSDK],
			}, nil
		}
	}

	data := pb.Record{}
	err = TypeConvert(newRecord, &data)
	if err != nil {
		utils.Error("get record info error:", err)
		return &pb.ResponseRecordCreate{
			ErrCode: nserr.ErrErr,
			ErrMsg:  utils.ErrCode[nserr.ErrErr],
		}, nil
	}
	info := &pb.ResponseRecordCreate{
		ErrCode: 0,
		ErrMsg:  "success",
		Data:    &data,
	}

	utils.Infof("[create record] key: %v,data: %#v", key, data)

	return info, nil
}

// RecordUpdate 更新解析记录
func (c *RPCServer) RecordUpdate(ctx context.Context, in *pb.RequestRecordUpdate) (*pb.ResponseRecordUpdate, error) {

	key, ok := IsOkRequest(ctx, in.RecordKey)
	if !ok || !IsRecordKey(key) {
		return &pb.ResponseRecordUpdate{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}

	if !CheckRecordPermission(key, int(in.RecordId)) {
		return &pb.ResponseRecordUpdate{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}

	var record dns.RecordInfo
	err := TypeConvert(in, &record)
	if err != nil {
		utils.Error("upadte record type convert error:", err)
		return &pb.ResponseRecordUpdate{
			ErrCode: nserr.ErrErr,
			ErrMsg:  utils.ErrCode[nserr.ErrErr],
		}, nil
	}
	d := dns.New("fone")
	record.ID = int(in.RecordId)
	err = d.UpdateRecord(record)
	if err != nil {
		utils.Errorf("update record by id: %d, error: %v", in.RecordId, err)
		switch err {
		case dns.ErrRecordSubDomain:
			return &pb.ResponseRecordUpdate{
				ErrCode: nserr.ErrRecordSubDomain,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordSubDomain],
			}, nil
		case dns.ErrRecordIP:
			return &pb.ResponseRecordUpdate{
				ErrCode: nserr.ErrRecordIP,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordIP],
			}, nil
		case dns.ErrRecordDomain:
			return &pb.ResponseRecordUpdate{
				ErrCode: nserr.ErrRecordDomain,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordDomain],
			}, nil
		case dns.ErrRecordTypeNotSupport:
			return &pb.ResponseRecordUpdate{
				ErrCode: nserr.ErrRecordTypeSupport,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordTypeSupport],
			}, nil
		case dns.ErrRecordMissTxt:
			return &pb.ResponseRecordUpdate{
				ErrCode: nserr.ErrRecordMissTxt,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordMissTxt],
			}, nil
		case dns.ErrRecordNotExist:
			return &pb.ResponseRecordUpdate{
				ErrCode: nserr.ErrRecordNotFound,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordNotFound],
			}, nil
		case dns.ErrRecordTypeMutex:
			return &pb.ResponseRecordUpdate{
				ErrCode: nserr.ErrRecordTypeMutex,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordTypeMutex],
			}, nil
		case dns.ErrRecordExist:
			return &pb.ResponseRecordUpdate{
				ErrCode: nserr.ErrRecordExist,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordExist],
			}, nil
		default:
			return &pb.ResponseRecordUpdate{
				ErrCode: nserr.ErrDNSSDK,
				ErrMsg:  utils.ErrCode[nserr.ErrDNSSDK],
			}, nil
		}
	}

	utils.Infof("[update record] key: %v,data: %#v", key, record)

	return &pb.ResponseRecordUpdate{
		ErrCode: 0,
		ErrMsg:  "success",
	}, nil
}

// RecordDelete 删除解析记录
func (c *RPCServer) RecordDelete(ctx context.Context, in *pb.RequestRecordDelete) (*pb.ResponseRecordDelete, error) {

	key, ok := IsOkRequest(ctx, in.RecordKey)
	if !ok || !IsRecordKey(key) {
		return &pb.ResponseRecordDelete{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}

	if !CheckRecordPermission(key, int(in.RecordId)) {
		return &pb.ResponseRecordDelete{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}

	d := dns.New("fone")
	err := d.DeleteRecord(int(in.RecordId))
	if err != nil {
		utils.Errorf("delete record by id: %d, error: %v", in.RecordId, err)
		if err == dns.ErrRecordNotExist {
			return &pb.ResponseRecordDelete{
				ErrCode: nserr.ErrRecordNotFound,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordNotFound],
			}, nil
		}
		return &pb.ResponseRecordDelete{
			ErrCode: nserr.ErrDNSSDK,
			ErrMsg:  utils.ErrCode[nserr.ErrDNSSDK],
		}, nil
	}

	utils.Infof("[delete record] key: %v,id: %v", key, in.RecordId)

	return &pb.ResponseRecordDelete{
		ErrCode: 0,
		ErrMsg:  "success",
	}, nil
}

// RecordDisable 暂停、启动解析记录
func (c *RPCServer) RecordDisable(ctx context.Context, in *pb.RequestRecordDisable) (*pb.ResponseRecordDisable, error) {

	return &pb.ResponseRecordDisable{
		ErrCode: 0,
		ErrMsg:  "success",
	}, nil
	key, ok := IsOkRequest(ctx, in.RecordKey)
	if !ok || !IsRecordKey(key) {
		return &pb.ResponseRecordDisable{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}

	if !CheckRecordPermission(key, int(in.RecordId)) {
		return &pb.ResponseRecordDisable{
			ErrCode: nserr.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[nserr.ErrPermissionDenied],
		}, nil
	}

	d := dns.New("fone")

	err := d.DisableRecord(int(in.RecordId), in.Disable)
	if err != nil {
		utils.Errorf("disable record by id: %d, error: %v", in.RecordId, err)
		if err == dns.ErrRecordNotExist {
			return &pb.ResponseRecordDisable{
				ErrCode: nserr.ErrRecordNotFound,
				ErrMsg:  utils.ErrCode[nserr.ErrRecordNotFound],
			}, nil
		}
		return &pb.ResponseRecordDisable{
			ErrCode: nserr.ErrDNSSDK,
			ErrMsg:  utils.ErrCode[nserr.ErrDNSSDK],
		}, nil
	}

	utils.Infof("[disable record] key: %v,id: %v,disable: %v", key, in.RecordId, in.Disable)

	return &pb.ResponseRecordDisable{
		ErrCode: 0,
		ErrMsg:  "success",
	}, nil
}

// CheckRecordPermission 校验Recordid 和 key 是否匹配
func CheckRecordPermission(recordKey string, recordId int) bool {

	key := "exist:" + recordKey + "+" + strconv.Itoa(recordId)
	exist, err := cache.DefaultCache.Exist(key)
	if err == nil && exist {
		return true
	}

	db := database.New()
	domain, err := db.GetDomainByRecordKey(recordKey)
	if err != nil {
		utils.Error("CheckRecordPermission error:", err)
		return false
	}

	record, err := db.GetRecordByReID(recordId)
	if err != nil {
		utils.Error("CheckRecordPermission error:", err)
		return false
	}
	if domain.ID == record.DomainID {
		go func() {
			cache.DefaultCache.Set(key, true, 5*time.Minute)
		}()
		return true
	}
	return false
}
