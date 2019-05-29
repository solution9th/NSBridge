package models

import (
	"github.com/solution9th/NSBridge/utils"
)

const (
	// WebErrDomainIDParam domainID 参数错误
	WebErrDomainIDParam = 10100

	// WebErrParams 参数错误
	WebErrParams = 10101

	// WebErrOffsetParam offset 参数错误
	WebErrOffsetParam = 10102

	// WebErrCountParam count 参数错误
	WebErrCountParam = 10103

	// WebErrParamsFormat 参数格式错误
	WebErrParamsFormat = 10104

	// WebErrDisableParam disable 参数错误
	WebErrDisableParam = 10105

	// WebErrAuthInsert 授权新增错误
	WebErrAuthInsert = 10205

	// WebErrAuthUpdate 授权更新错误
	WebErrAuthUpdate = 10206

	// WebErrAuthDisable 启禁用错误
	WebErrAuthDisable = 10207

	// WebErrAuthDelete 授权删除错误
	WebErrAuthDelete = 10208

	// WebErrAuthIDNotFound authID 找不到
	WebErrAuthIDNotFound = 10209

	// WebErrSrever 服务器错误
	WebErrSrever = 10301

	// WebErrEmptyData 数据为空
	WebErrEmptyData = 10302

	// WebErrUserNotLogin 用户未登录
	WebErrUserNotLogin = 10303

	// WebErrDelAuthNotdisable  删除未禁用的授权错误
	WebErrDelAuthNotDisable = 10304
)

const (
	// ErrErr just test
	ErrErr = 201233

	// ErrPermissionDenied 没有权限
	ErrPermissionDenied = 200199

	// ErrParams 参数错误
	ErrParams = 20100

	// ErrDB 数据库错误
	ErrDB = 200101

	// ErrDomainExist 域名已经存在
	ErrDomainExist = 200102

	// ErrDomainNotFound 域名不存在
	ErrDomainNotFound = 200103

	// ErrRecordNotFound 记录不存在
	ErrRecordNotFound = 200104

	// ErrDNSSDK sdk 错误
	ErrDNSSDK = 200105

	// ErrDomainExist 域名已经存在
	ErrRecordTypeSupport = 200106

	// ErrRecordTypeMutex 解析记录互斥
	ErrRecordTypeMutex = 200107

	// ErrRecordExist 解析记录已存在
	ErrRecordExist = 200108

	// ErrRecordMissTxt txt记录不可为空
	ErrRecordMissTxt = 200109

	// ErrRecordIP ip记录不正确
	ErrRecordIP = 200110

	// ErrRecordDomain 域名记录不正确
	ErrRecordDomain = 200111

	// ErrRecordSubDomain sub_domain 不正确
	ErrRecordSubDomain = 200112
)

var (
	webErrCode = map[int]string{
		WebErrDomainIDParam:     "domain_id参数错误",
		WebErrParams:            "参数错误",
		WebErrOffsetParam:       "offset参数错误",
		WebErrCountParam:        "count参数错误",
		WebErrParamsFormat:      "参数格式错误",
		WebErrDisableParam:      "disable参数错误",
		WebErrAuthUpdate:        "服务器异常,请联系管理员",
		WebErrAuthDisable:       "服务器异常,请联系管理员",
		WebErrAuthDelete:        "服务器异常,请联系管理员",
		WebErrAuthInsert:        "服务器异常,请联系管理员",
		WebErrSrever:            "服务器异常,请联系管理员",
		WebErrAuthIDNotFound:    "服务器异常,请联系管理员",
		WebErrEmptyData:         "查询结果为空",
		WebErrUserNotLogin:      "没有登录",
		WebErrDelAuthNotDisable: "启用状态的授权信息不可删除",
	}

	apiErrCode = map[int]string{
		ErrErr:               "error",
		ErrParams:            "params error",
		ErrDB:                "database error",
		ErrDomainExist:       "domain exist",
		ErrDomainNotFound:    "domain not found",
		ErrRecordNotFound:    "record not found",
		ErrDNSSDK:            "dns sdk error",
		ErrPermissionDenied:  "permission denied",
		ErrRecordTypeSupport: "record type not support",

		ErrRecordTypeMutex: "record type is mutex",
		ErrRecordExist:     "record is duplicate, can't create it again",
		ErrRecordMissTxt:   "type txt miss txt content",
		ErrRecordIP:        "error ip",
		ErrRecordDomain:    "error domain value",
		ErrRecordSubDomain: "error sub_domain",
	}
)

func init() {
	utils.AddErrCodes(webErrCode)
	utils.AddErrCodes(apiErrCode)
}
