package grpc

import (
	"context"
	"strconv"
	"time"

	"github.com/solution9th/NSBridge/dns"
	pb "github.com/solution9th/NSBridge/dns_pb"
	"github.com/solution9th/NSBridge/models"
	"github.com/solution9th/NSBridge/service/cache"
	"github.com/solution9th/NSBridge/service/database"
	"github.com/solution9th/NSBridge/utils"
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
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
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
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
		}, nil
	}

	list, err := db.GetAllRecordByDomainID(domain.ID, uint(in.Start), uint(in.Count))
	if err != nil {
		utils.Error("get all record error:", err)
		return &pb.ResponseRecordList{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
		}, nil
	}

	totalNum, err := db.CountAllRecordNum(domain.ID)
	if err != nil {
		utils.Error("get total num error:", err)
		return &pb.ResponseRecordList{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
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
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	db := database.New()
	m, err := db.GetDomainByRecordKey(key)
	if err != nil {
		utils.Error("get record info error:", err)
		return &pb.ResponseDomainOfRK{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
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
			ErrCode: models.ErrDNSSDK,
			ErrMsg:  utils.ErrCode[models.ErrDNSSDK],
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
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}
	if !CheckRecordPermission(key, int(in.RecordId)) {
		return &pb.ResponseRecordInfo{
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	db := database.New()
	exist, err := db.ExistRecordByRecordID(int(in.RecordId))
	if err != nil {
		utils.Error("exist record err", err)
	}
	if !exist {
		return &pb.ResponseRecordInfo{
			ErrCode: models.ErrRecordNotFound,
			ErrMsg:  utils.ErrCode[models.ErrRecordNotFound],
		}, nil
	}
	record, err := db.GetRecordByReID(int(in.RecordId))
	if err != nil {
		utils.Error("get record info error:", err)
		return &pb.ResponseRecordInfo{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
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
			ErrCode: models.ErrErr,
			ErrMsg:  utils.ErrCode[models.ErrErr],
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
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	if in.SubDomain == "" {
		return &pb.ResponseRecordCreate{
			ErrCode: models.ErrParams,
		}, nil
	}

	if in.RecordType == "" {
		return &pb.ResponseRecordCreate{
			ErrCode: models.ErrParams,
			ErrMsg:  utils.ErrCode[models.ErrParams] + ": miss record_type",
		}, nil
	}

	if in.Value == "" {
		return &pb.ResponseRecordCreate{
			ErrCode: models.ErrParams,
			ErrMsg:  utils.ErrCode[models.ErrParams] + ": miss value",
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
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
		}, nil
	}

	d := dns.New("fone")
	var record dns.RecordInfo
	err = TypeConvert(in, &record)
	if err != nil {
		utils.Error("create record type convert error:", err)
		return &pb.ResponseRecordCreate{
			ErrCode: models.ErrErr,
			ErrMsg:  utils.ErrCode[models.ErrErr],
		}, nil
	}
	record.DomainID = domain.ID
	newRecord, err := d.CreateRecord(record)
	if err != nil {
		utils.Error("create dns error:", err)
		switch err {
		case dns.ErrRecordSubDomain:
			return &pb.ResponseRecordCreate{
				ErrCode: models.ErrRecordSubDomain,
				ErrMsg:  utils.ErrCode[models.ErrRecordSubDomain],
			}, nil
		case dns.ErrRecordIP:
			return &pb.ResponseRecordCreate{
				ErrCode: models.ErrRecordIP,
				ErrMsg:  utils.ErrCode[models.ErrRecordIP],
			}, nil
		case dns.ErrRecordDomain:
			return &pb.ResponseRecordCreate{
				ErrCode: models.ErrRecordDomain,
				ErrMsg:  utils.ErrCode[models.ErrRecordDomain],
			}, nil
		case dns.ErrRecordMissTxt:
			return &pb.ResponseRecordCreate{
				ErrCode: models.ErrRecordMissTxt,
				ErrMsg:  utils.ErrCode[models.ErrRecordMissTxt],
			}, nil
		case dns.ErrDomainNotExist:
			return &pb.ResponseRecordCreate{
				ErrCode: models.ErrDomainExist,
				ErrMsg:  utils.ErrCode[models.ErrDomainExist],
			}, nil
		case dns.ErrRecordTypeNotSupport:
			return &pb.ResponseRecordCreate{
				ErrCode: models.ErrRecordTypeSupport,
				ErrMsg:  utils.ErrCode[models.ErrRecordTypeSupport],
			}, nil
		case dns.ErrRecordTypeMutex: // 互斥
			return &pb.ResponseRecordCreate{
				ErrCode: models.ErrRecordTypeMutex,
				ErrMsg:  utils.ErrCode[models.ErrRecordTypeMutex],
			}, nil
		case dns.ErrRecordExist:
			return &pb.ResponseRecordCreate{
				ErrCode: models.ErrRecordExist,
				ErrMsg:  utils.ErrCode[models.ErrRecordExist],
			}, nil
		default:
			return &pb.ResponseRecordCreate{
				ErrCode: models.ErrDNSSDK,
				ErrMsg:  utils.ErrCode[models.ErrDNSSDK],
			}, nil
		}
	}

	data := pb.Record{}
	err = TypeConvert(newRecord, &data)
	if err != nil {
		utils.Error("get record info error:", err)
		return &pb.ResponseRecordCreate{
			ErrCode: models.ErrErr,
			ErrMsg:  utils.ErrCode[models.ErrErr],
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
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	if !CheckRecordPermission(key, int(in.RecordId)) {
		return &pb.ResponseRecordUpdate{
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	var record dns.RecordInfo
	err := TypeConvert(in, &record)
	if err != nil {
		utils.Error("upadte record type convert error:", err)
		return &pb.ResponseRecordUpdate{
			ErrCode: models.ErrErr,
			ErrMsg:  utils.ErrCode[models.ErrErr],
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
				ErrCode: models.ErrRecordSubDomain,
				ErrMsg:  utils.ErrCode[models.ErrRecordSubDomain],
			}, nil
		case dns.ErrRecordIP:
			return &pb.ResponseRecordUpdate{
				ErrCode: models.ErrRecordIP,
				ErrMsg:  utils.ErrCode[models.ErrRecordIP],
			}, nil
		case dns.ErrRecordDomain:
			return &pb.ResponseRecordUpdate{
				ErrCode: models.ErrRecordDomain,
				ErrMsg:  utils.ErrCode[models.ErrRecordDomain],
			}, nil
		case dns.ErrRecordTypeNotSupport:
			return &pb.ResponseRecordUpdate{
				ErrCode: models.ErrRecordTypeSupport,
				ErrMsg:  utils.ErrCode[models.ErrRecordTypeSupport],
			}, nil
		case dns.ErrRecordMissTxt:
			return &pb.ResponseRecordUpdate{
				ErrCode: models.ErrRecordMissTxt,
				ErrMsg:  utils.ErrCode[models.ErrRecordMissTxt],
			}, nil
		case dns.ErrRecordNotExist:
			return &pb.ResponseRecordUpdate{
				ErrCode: models.ErrRecordNotFound,
				ErrMsg:  utils.ErrCode[models.ErrRecordNotFound],
			}, nil
		case dns.ErrRecordTypeMutex:
			return &pb.ResponseRecordUpdate{
				ErrCode: models.ErrRecordTypeMutex,
				ErrMsg:  utils.ErrCode[models.ErrRecordTypeMutex],
			}, nil
		case dns.ErrRecordExist:
			return &pb.ResponseRecordUpdate{
				ErrCode: models.ErrRecordExist,
				ErrMsg:  utils.ErrCode[models.ErrRecordExist],
			}, nil
		default:
			return &pb.ResponseRecordUpdate{
				ErrCode: models.ErrDNSSDK,
				ErrMsg:  utils.ErrCode[models.ErrDNSSDK],
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
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	if !CheckRecordPermission(key, int(in.RecordId)) {
		return &pb.ResponseRecordDelete{
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	d := dns.New("fone")
	err := d.DeleteRecord(int(in.RecordId))
	if err != nil {
		utils.Errorf("delete record by id: %d, error: %v", in.RecordId, err)
		if err == dns.ErrRecordNotExist {
			return &pb.ResponseRecordDelete{
				ErrCode: models.ErrRecordNotFound,
				ErrMsg:  utils.ErrCode[models.ErrRecordNotFound],
			}, nil
		}
		return &pb.ResponseRecordDelete{
			ErrCode: models.ErrDNSSDK,
			ErrMsg:  utils.ErrCode[models.ErrDNSSDK],
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
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	if !CheckRecordPermission(key, int(in.RecordId)) {
		return &pb.ResponseRecordDisable{
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	d := dns.New("fone")

	err := d.DisableRecord(int(in.RecordId), in.Disable)
	if err != nil {
		utils.Errorf("disable record by id: %d, error: %v", in.RecordId, err)
		if err == dns.ErrRecordNotExist {
			return &pb.ResponseRecordDisable{
				ErrCode: models.ErrRecordNotFound,
				ErrMsg:  utils.ErrCode[models.ErrRecordNotFound],
			}, nil
		}
		return &pb.ResponseRecordDisable{
			ErrCode: models.ErrDNSSDK,
			ErrMsg:  utils.ErrCode[models.ErrDNSSDK],
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
