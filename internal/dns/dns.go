package dns

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/solution9th/NSBridge/internal/config"
	"github.com/solution9th/NSBridge/internal/sdk"
	"github.com/solution9th/NSBridge/internal/utils"
)

// GetDefaultNS 获得 fone 默认 ns
func GetDefaultNS() []sdk.Authority {

	datas := make([]sdk.Authority, 0)

	data := sdk.Authority{
		TTL:  1,
		Unit: sdk.RecordDayUnit,
	}

	for _, v := range config.NameServerConfig {
		data.Domain = v
		datas = append(datas, data)
	}

	return datas
}

// GetRecordKey 获得操作域名的 appkey
func GetRecordKey(domain string) string {

	key := md5.Sum([]byte(fmt.Sprintf("%s%v", domain, time.Now().Unix())))

	keyStr := fmt.Sprintf("%x", key)
	keyStr = keyStr[8 : len(keyStr)-8]

	// record base64 => cmVjb3Jk
	return fmt.Sprintf("%s%s", GetRecordKeyPrefix(), keyStr)
}

// GetRecordSecret 获得操作域名的秘钥
func GetRecordSecret() string {
	size := 20
	return utils.GenToken(size)
}

// GetRecordKeyPrefix record app-key 生成前缀
func GetRecordKeyPrefix() string {
	return "cmVjb3Jk"
}
