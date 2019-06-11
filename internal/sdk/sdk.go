package sdk

type DNSSDK interface {
	Login() (string, error)

	// 域名操作

	// 添加新的 DNS 域名
	CreateDNSDomain(domain string, authority []Authority, soaEmail string) (resp CreateDomainResponse, err error)

	// 更新 DNS 域名
	UpdateDNSDomain(dnsID int, domain string, authority []Authority, soaEmail string) (err error)

	// 查看 DNS 域名详情
	GetDNSDomain(dnsID int) (dnsInfo DNSInfo, err error)

	// 获得所有 DNS 域名
	GetAllDNSDomain() (infoList DNSInfos, err error)

	// 删除 DNS 域名
	DeleteDNSDomain(dnsID int) error

	// 记录操作

	// 获得 DNS 记录详情
	GetDNSRecord(dnsID, recordID int) (record Record, err error)

	// 获得所有 DNS 记录
	GetAllDNSRecords(dnsID int) (infos RecordInfos, err error)

	// 添加新的 DNS 记录
	CreateDNSRecord(dnsID int, subDomain string, recordType RecordType, value string, line, priority, ttl int, recordUnit RecordUnit) (recordid int, err error)

	// 更新 DNS 记录
	UpdateDNSRecord(dnsID, recordID int, subDomain string, recordType RecordType, value string, line, priority, ttl int, recordUnit RecordUnit) (err error)

	// 删除 DNS 记录
	DeleteDNSRecord(dnsID, recordID int) (err error)

	// 暂停启动 DNS 记录
	DisableDNSRecord(dnsID, recordID int, disable bool) (err error)
}
