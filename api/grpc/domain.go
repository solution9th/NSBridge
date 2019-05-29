package grpc

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/solution9th/NSBridge/dns"
	"github.com/solution9th/NSBridge/models"
	"github.com/solution9th/NSBridge/sdk"
	"github.com/solution9th/NSBridge/service/cache"
	"github.com/solution9th/NSBridge/service/database"
	"github.com/solution9th/NSBridge/utils"

	"github.com/haozibi/gendry/scanner"
	pb "github.com/solution9th/NSBridge/dns_pb"
)

var (
	// ErrMetaData 从元数据中获得 metadata 错误
	ErrMetaData = errors.New("Get MetaData From Ctx error")
)

// DomainsList 展示所有域名
func (r *RPCServer) DomainsList(ctx context.Context, in *pb.RequestDomainsList) (*pb.ResponseDomainList, error) {

	if in.Count == 0 {
		in.Count = 10
	}

	if in.Count > 100 {
		in.Count = 50
	}

	db := database.New()

	list, err := db.GetAllDomains(uint(in.Start), uint(in.Count))
	if err != nil {
		utils.Error("get all domains error:", err)
		return &pb.ResponseDomainList{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
		}, nil
	}

	totalNum, err := db.CountAllDomainNum()
	if err != nil {
		utils.Error("get total num error:", err)
		return &pb.ResponseDomainList{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
		}, nil
	}

	info := &pb.ResponseDomainList{}
	info.ErrCode = 0
	info.ErrMsg = "success"
	info.Data = &pb.DomainListData{}
	datas := make([]*pb.DomainInfo, 0)

	for _, v := range list {
		data := &pb.DomainInfo{}
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

// OwnDomainsList 获取key下所有域名
func (r *RPCServer) OwnDomainsList(ctx context.Context, in *pb.RequestOwnDomainsList) (*pb.ResponseDomainList, error) {

	key, ok := IsOkRequest(ctx, in.Key)
	if !ok || IsRecordKey(key) {
		utils.Errorf("[DomainCreate-PD] key: %v,ok: %v,rpcKey: %v, isRecordKey: %v", key, ok, in.Key, IsRecordKey(key))
		return &pb.ResponseDomainList{
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	if in.Count == 0 {
		in.Count = 10
	}

	if in.Count > 100 {
		in.Count = 50
	}

	db := database.New()

	list, err := db.GetDomainByDomainKey(key, uint(in.Start), uint(in.Count))
	if err != nil {
		utils.Error("get own domains error:", err)
		return &pb.ResponseDomainList{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
		}, nil
	}

	totalNum, err := db.CountOwnDomainNum(key)
	if err != nil {
		utils.Error("get total num error:", err)
		return &pb.ResponseDomainList{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
		}, nil
	}

	info := &pb.ResponseDomainList{}
	info.ErrCode = 0
	info.ErrMsg = "success"
	info.Data = &pb.DomainListData{}
	datas := make([]*pb.DomainInfo, 0)

	for _, v := range list {
		data := &pb.DomainInfo{}
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

// DomainCreate 创建域名
func (r *RPCServer) DomainCreate(ctx context.Context, in *pb.RequestDomainCreate) (*pb.ResponseDomainCreate, error) {

	key, ok := IsOkRequest(ctx, in.ApiKey)
	if !ok || IsRecordKey(key) {
		utils.Errorf("[DomainCreate-PD] key: %v,ok: %v,rpcKey: %v, isRecordKey: %v", key, ok, in.ApiKey, IsRecordKey(key))
		return &pb.ResponseDomainCreate{
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	in.Domain = strings.TrimSuffix(in.Domain, ".")

	if !utils.IsOkDomain(in.Domain) {
		return &pb.ResponseDomainCreate{
			ErrCode: models.ErrParams,
			ErrMsg:  utils.ErrCode[models.ErrParams] + ": domain error",
		}, nil
	}

	if in.Domain == "" {
		return &pb.ResponseDomainCreate{
			ErrCode: models.ErrParams,
			ErrMsg:  utils.ErrCode[models.ErrParams] + ": miss domain",
		}, nil
	}

	db := database.New()

	mAuth, err := db.GetAuthByKey(key)
	if err != nil {
		if err == scanner.ErrEmptyResult {
			utils.Errorf("[DomainCreate-PD-notfind] key: %v,ok: %v", in.ApiKey, ok)
			return &pb.ResponseDomainCreate{
				ErrCode: models.ErrPermissionDenied,
				ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
			}, nil
		}
		utils.Error("get auth error:", err)
		return &pb.ResponseDomainCreate{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
		}, nil
	}

	d := dns.New("fone")
	m, err := d.Create(mAuth.DomainKey, in.Domain, in.Remark)
	if err != nil {
		if err == dns.ErrDomainHasExist {
			return &pb.ResponseDomainCreate{
				ErrCode: models.ErrDomainExist,
				ErrMsg:  utils.ErrCode[models.ErrDomainExist],
			}, nil
		}
		utils.Error("create dns error:", err)
		return &pb.ResponseDomainCreate{
			ErrCode: models.ErrDNSSDK,
			ErrMsg:  utils.ErrCode[models.ErrDNSSDK],
		}, nil
	}

	info := &pb.ResponseDomainCreate{
		ErrCode: 0,
		ErrMsg:  "success",
	}

	data := &pb.DomainInfo{}

	err = TypeConvert(m, data)
	if err != nil {
		utils.Error("type error:", err)
		return &pb.ResponseDomainCreate{
			ErrCode: models.ErrDNSSDK,
			ErrMsg:  utils.ErrCode[models.ErrDNSSDK],
		}, nil
	}

	info.Data = data

	utils.Infof("[create domain] key: %v,info: %#v", key, data)

	return info, nil
}

// DomainDelete 删除域名
func (r *RPCServer) DomainDelete(ctx context.Context, in *pb.RequestDomainDelete) (*pb.ResponseDomainDelete, error) {

	key, ok := IsOkRequest(ctx, in.ApiKey)
	if !ok || IsRecordKey(key) {
		utils.Errorf("[DomainDelete-PD] key: %v,ok: %v,rpcKey: %v, isRecordKey: %v", key, ok, in.ApiKey, IsRecordKey(key))
		// utils.Errorf("[DomainDelete-PD] key: %v,ok: %v", in.ApiKey, ok)
		return &pb.ResponseDomainDelete{
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	id := int(in.Id)

	db := database.New()
	m, err := db.GetDomainByID(id)
	if err != nil {
		if err == scanner.ErrEmptyResult {
			return &pb.ResponseDomainDelete{
				ErrCode: models.ErrDomainNotFound,
				ErrMsg:  utils.ErrCode[models.ErrDomainNotFound],
			}, nil
		}
		utils.Errorf("get domain by id: %d error: %v", id, err)
		return &pb.ResponseDomainDelete{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
		}, nil
	}

	// 只能用创建的 key 删除创建的域名
	if m.DomainKey != key {
		utils.Errorf("[DomainDelete-PD] key: %v,DomainKey: %v", key, m.DomainKey)
		return &pb.ResponseDomainDelete{
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	d := dns.New("fone")
	err = d.DeleteDomain(id, m.FoneDomainID)
	if err != nil {
		utils.Errorf("delete domain by id: %d,domainid: %d, error: %v", id, m.FoneDomainID, err)
		return &pb.ResponseDomainDelete{
			ErrCode: models.ErrDNSSDK,
			ErrMsg:  utils.ErrCode[models.ErrDNSSDK],
		}, nil
	}

	utils.Infof("[delete domain] key: %v,id: %v", key, id)

	return &pb.ResponseDomainDelete{
		ErrCode: 0,
		ErrMsg:  "success",
	}, nil
}

// DomainStatus 查看域名托管状态
// key 是操控记录的 key
func (r *RPCServer) DomainStatus(ctx context.Context, in *pb.RequestDomainStatus) (*pb.ResponseDomainStatus, error) {

	key, ok := IsOkRequest(ctx, in.ApiKey)
	if !ok || !IsRecordKey(key) {
		utils.Errorf("[DomainStatus-PD] key: %v,ok: %v,rpcKey: %v, isRecordKey: %v", key, ok, in.ApiKey, IsRecordKey(key))
		// utils.Errorf("[DomainStatus-PD] key: %v,ok: %v", in.ApiKey, ok)
		return &pb.ResponseDomainStatus{
			ErrCode: models.ErrPermissionDenied,
			ErrMsg:  utils.ErrCode[models.ErrPermissionDenied],
		}, nil
	}

	result := &pb.ResponseDomainStatus{}

	cacheKey := utils.GetCacheTakeOverKey(key)

	err := cache.DefaultCache.Get(cacheKey, &result)
	if err == nil {
		utils.Infof("get is_take_over from cache: %v", cacheKey)
		return result, nil
	}

	db := database.New()
	m, err := db.GetDomainByRecordKey(key)
	if err != nil {
		if err == scanner.ErrEmptyResult {
			return &pb.ResponseDomainStatus{
				ErrCode: models.ErrDomainNotFound,
				ErrMsg:  utils.ErrCode[models.ErrDomainNotFound],
			}, nil
		}
		utils.Errorf("get domain by id: %s error: %v", key, err)
		return &pb.ResponseDomainStatus{
			ErrCode: models.ErrDB,
			ErrMsg:  utils.ErrCode[models.ErrDB],
		}, nil
	}

	ns, err := sdk.GetDomainNS(m.Domain)
	if err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			utils.Debugf("get ns timeout:%v", err)
			return &pb.ResponseDomainStatus{
				// ErrCode: models.ErrDNSSDK,
				// ErrMsg:  utils.ErrCode[models.ErrDNSSDK],
				Data: &pb.TakeOver{
					IsTakeOver: int32(0),
				},
			}, nil
		}
		utils.Error("get ns error:", err)
		return &pb.ResponseDomainStatus{
			ErrCode: models.ErrDNSSDK,
			ErrMsg:  utils.ErrCode[models.ErrDNSSDK],
		}, nil
	}

	isTakeOver := 0

	for _, a := range ns {
		for _, b := range *m.NameServer {
			if a == b {
				isTakeOver = 1
				break
			}
		}
	}

	where := map[string]interface{}{
		"id": m.ID,
	}

	update := make(map[string]interface{})

	if isTakeOver == 1 && m.IsTakeOver != 1 {
		update["is_take_over"] = 1
	} else if isTakeOver == 0 && m.IsTakeOver != 0 {
		update["is_take_over"] = 0
	}

	if len(update) > 0 {
		err = db.UpdateDomain(where, update)
		if err != nil {
			utils.Error("update domain error:", m.ID, update, err)
			return &pb.ResponseDomainStatus{
				ErrCode: models.ErrDB,
				ErrMsg:  utils.ErrCode[models.ErrDB],
			}, nil
		}
	}

	tmpData := &pb.TakeOver{
		IsTakeOver: int32(isTakeOver),
	}

	result.Data = tmpData

	go func() {
		cache.DefaultCache.Set(cacheKey, *result, 5*time.Minute)
	}()

	return result, nil
}
