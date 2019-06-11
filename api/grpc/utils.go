package grpc

import (
	"context"
	"encoding/json"
	"net/textproto"
	"strings"

	"github.com/solution9th/NSBridge/internal/dns"
	"github.com/solution9th/NSBridge/internal/utils"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/metadata"
)

const (
	// MetadataPrefix 请求头前缀
	MetadataPrefix = "gg-"

	// MetadataHeaderPrefix grpc 中的前缀
	MetadataHeaderPrefix = "Gg-"
)

// GatewayHaderMatcher 自定义 gRPC 网关 header key 转换方式
// 如果是 http 固定 header，取值的 key(例如：Cookie) => gg-cookie
// 如果是自定义的 header 则必须有 Gg- 前缀（例如: Gg-Token），取值 key => token
func GatewayHaderMatcher(key string) (string, bool) {

	key = textproto.CanonicalMIMEHeaderKey(key)

	if isPermanentHTTPHeader(key) {
		return MetadataPrefix + key, true
	} else if strings.HasPrefix(key, MetadataHeaderPrefix) {
		return key[len(MetadataHeaderPrefix):], true
	}
	return "", false
}

func isPermanentHTTPHeader(hdr string) bool {
	switch hdr {
	case
		"Accept",
		"Accept-Charset",
		"Accept-Language",
		"Accept-Ranges",
		"Authorization",
		"Cache-Control",
		"Content-Type",
		"Cookie",
		"Date",
		"Expect",
		"From",
		"Host",
		"If-Match",
		"If-Modified-Since",
		"If-None-Match",
		"If-Schedule-Tag-Match",
		"If-Unmodified-Since",
		"Max-Forwards",
		"Origin",
		"Pragma",
		"Referer",
		"User-Agent",
		"Via",
		"Warning",
		"X-Api-Key",
		"X-Api-Timestamp",
		"X-Api-Hmac":
		return true
	}
	return false
}

func respReturn(code int, msg string, data interface{}) *httpbody.HttpBody {

	resp := utils.ParseResult(code, msg, data)

	body, _ := json.Marshal(resp)

	return &httpbody.HttpBody{
		ContentType: "text/json",
		Data:        body,
	}
}

// IsOkRequest 判断请求是否通过鉴权的
// gRPCKey 是在 grpc 请求中的 key
func IsOkRequest(ctx context.Context, gRPCKey string) (string, bool) {

	key := ""

	if IsHTTPRequest(ctx) {
		key = GetHeaderKey(ctx)
		return key, true
	}

	key = gRPCKey
	secret := GetSecret(key)
	if secret == "" {
		return key, false
	}
	return key, true
}

// IsHTTPRequest 判断流量是否来自 http
func IsHTTPRequest(ctx context.Context) bool {
	value := getGRPCValue(ctx, GRPCFlag)
	if value != "" {
		return true
	}
	return false
}

// IsRecordKey 判断 key 是否是操作 record 的 key
func IsRecordKey(key string) bool {
	if key == "" {
		return false
	}
	return strings.HasPrefix(key, dns.GetRecordKeyPrefix())
}

// GetHeaderKey 获取头文件中的 key
func GetHeaderKey(ctx context.Context) string {
	return getGRPCValue(ctx, "gg-x-api-key")
}

func getGRPCValue(ctx context.Context, name string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	name = strings.ToLower(name)

	key := ""

	keys := md.Get(name)
	if len(keys) > 0 {
		key = keys[0]
	}

	return key
}
