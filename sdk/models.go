package sdk

// RecordType dns record type [A, CNAME, MX, TXT]
type RecordType string

func (r RecordType) String() string {
	return string(r)
}

const (
	// RecordAType record type A
	RecordAType RecordType = "A"

	// RecordCNAMEType record type CNAME
	RecordCNAMEType RecordType = "CNAME"

	// RecordMXType record type MX
	RecordMXType RecordType = "MX"

	// RecordTXTType record type TXT
	RecordTXTType RecordType = "TXT"

	// RecordNSType record type NS
	RecordNSType RecordType = "NS"
)

var (
	// AllowTypeList 支持的 type 列表
	AllowTypeList = []RecordType{RecordAType,
		RecordCNAMEType, RecordMXType,
		RecordTXTType, RecordNSType}
)

// RecordUnit unit [sec, min, hour, day]
type RecordUnit string

func (r RecordUnit) String() string {
	return string(r)
}

const (
	// RecordSecUnit sec
	RecordSecUnit RecordUnit = "sec"

	// RecordMinUnit min
	RecordMinUnit RecordUnit = "min"

	// RecordHourUnit hour
	RecordHourUnit RecordUnit = "hour"

	// RecordDayUnit day
	RecordDayUnit RecordUnit = "day"
)

var (
	// AllowUnitList 允许的 Unit 列表
	AllowUnitList = []RecordUnit{
		RecordSecUnit, RecordMinUnit,
		RecordHourUnit, RecordDayUnit}
)

// Authority authority dns server
type Authority struct {
	ID     int        `json:"id,omitempty"`
	Domain string     `json:"domain"`
	TTL    int        `json:"ttl"`
	Unit   RecordUnit `json:"unit"`
}

// DNSInfo dns info
type DNSInfo struct {
	ModifiedUnix float64     `json:"_modified_unix"`
	ID           int         `json:"id,omitempty"`
	SoaEmail     string      `json:"soa_email"`
	Nameserver   []Authority `json:"nameserver"`
	Zone         string      `json:"zone"`
	CreatedUnix  float64     `json:"_created_unix"`
}

// DNSInfos dns info list
type DNSInfos struct {
	Data []DNSInfo `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
}

// Record dns record
type Record struct {
	ID           int        `json:"id,omitempty"`
	ModifiedUnix float64    `json:"_modified_unix"`
	CreatedUnix  float64    `json:"_created_unix"`
	Text         string     `json:"text"`
	Unit         RecordUnit `json:"unit"` // [sec, min, hour, day]
	Priority     int        `json:"priority"`
	SubDomain    string     `json:"sub_domain"`
	Disabled     bool       `json:"disabled"`
	TTL          int        `json:"ttl"`
	Line         int        `json:"line"`
	Type         RecordType `json:"type"` // [A, CNAME, MX, TXT]
}

// RecordInfos dns record info list
type RecordInfos struct {
	Data []Record `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
}

type CreateDomainResponse struct {
	ID int      `json:"id,omitempty"`
	NS []string `json:"ns,omitempty"`
}
