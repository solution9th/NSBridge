package fone

import (
	"fmt"

	"github.com/solution9th/NSBridge/sdk"
	"github.com/solution9th/NSBridge/utils"
)

// CreateDNSDomain add a new dns
func (f *FOneDNS) CreateDNSDomain(domain string, authority []sdk.Authority, soaEmail string) (resp sdk.CreateDomainResponse, err error) {

	if !IsOkDomain(domain) {
		err = ErrDomainParams
		return
	}

	if len(authority) == 0 {
		err = ErrMissNSParams
		return
	}

	for k, v := range authority {
		if !IsOkDomain(v.Domain) {
			err = ErrNSParams
			return
		}
		if v.TTL <= 0 {
			authority[k].TTL = 1
		}

		if v.Unit.String() == "" {
			authority[k].Unit = sdk.RecordHourUnit
		}
	}

	uri := fmt.Sprintf("/dns")

	req := struct {
		Zone       string          `json:"zone"`
		NameServer []sdk.Authority `json:"nameserver"`
		SOAEmail   string          `json:"soa_email"`
	}{
		Zone:       domain,
		NameServer: authority,
		SOAEmail:   soaEmail,
	}

	err = f.httpDo("POST", uri, req, &resp)
	if err != nil {
		utils.Error("[fone] new dns error:", err)
		return
	}

	for _, v := range authority {
		resp.NS = append(resp.NS, v.Domain)
	}

	return
}

// UpdateDNSDomain update dns app
// 当 domain 为空时保持原有 domain
func (f *FOneDNS) UpdateDNSDomain(dnsID int, domain string, authority []sdk.Authority, soaEmail string) (err error) {

	uri := fmt.Sprintf("/dns/%d", dnsID)

	if dnsID <= 0 {
		err = ErrDNSIDParams
		return
	}

	if domain != "" && !IsOkDomain(domain) {
		err = ErrDomainParams
		return
	}

	if len(authority) == 0 {
		err = ErrMissNSParams
		return
	}

	for k, v := range authority {
		if !IsOkDomain(v.Domain) {
			err = ErrNSParams
			return
		}
		if v.TTL <= 0 {
			authority[k].TTL = 1
		}

		if v.Unit.String() == "" {
			authority[k].Unit = sdk.RecordHourUnit
		}
	}

	req := struct {
		Zone       string          `json:"zone,omitempty"`
		NameServer []sdk.Authority `json:"nameserver"`
		SOAEmail   string          `json:"soa_email"`
	}{
		Zone:       domain,
		NameServer: authority,
		SOAEmail:   soaEmail,
	}

	err = f.httpDo("PUT", uri, req, nil)
	if err != nil {
		utils.Error("[fone] new dns error:", err)
	}

	return
}

// GetDNSDomain get dns by dns_id
func (f *FOneDNS) GetDNSDomain(dnsID int) (dnsInfo sdk.DNSInfo, err error) {

	if dnsID <= 0 {
		err = ErrDNSIDParams
		return
	}

	uri := fmt.Sprintf("/dns/%d", dnsID)

	err = f.httpDo("GET", uri, nil, &dnsInfo)
	if err != nil {
		utils.Error("[fone] get dns error:", err)
	}

	return
}

// GetAllDNSDomain 获得所有的 Domain
func (f *FOneDNS) GetAllDNSDomain() (infoList sdk.DNSInfos, err error) {

	page := 1
	pagesize := 20

	for {
		tmpInfo, err := f.getAllDomain(page, pagesize)
		if err != nil {
			return infoList, err
		}

		ll := len(tmpInfo.Data)

		infoList.Data = append(infoList.Data, tmpInfo.Data...)
		infoList.Meta.Count = tmpInfo.Meta.Count

		if ll <= pagesize && ll > 0 {
			page++
			continue
		}
		break
	}

	return
}

func (f *FOneDNS) getAllDomain(page, pagesize int) (infoList sdk.DNSInfos, err error) {

	uri := fmt.Sprintf("/dns?page=%d&pagesize=%d", page, pagesize)

	err = f.httpDo("GET", uri, nil, &infoList)
	if err != nil {
		utils.Error("[fone] get dns error:", err)
	}
	return
}

// DeleteDNSDomain delete dns app
func (f *FOneDNS) DeleteDNSDomain(dnsID int) (err error) {

	if dnsID <= 0 {
		err = ErrDNSIDParams
		return
	}

	uri := fmt.Sprintf("/dns/%d", dnsID)

	return f.httpDo("DELETE", uri, nil, nil)
}

//###############
// dns record
//###############

// GetDNSRecord get dns record by id
func (f *FOneDNS) GetDNSRecord(dnsID, recordID int) (record sdk.Record, err error) {

	if dnsID <= 0 || recordID <= 0 {
		err = ErrIDParams
		return
	}

	uri := fmt.Sprintf("/dns/%d/record/%d", dnsID, recordID)

	err = f.httpDo("GET", uri, nil, &record)
	if err != nil {
		utils.Error("[fone] get dns error:", err)
	}

	return
}

// GetAllDNSRecords get all dns record
func (f *FOneDNS) GetAllDNSRecords(dnsID int) (infos sdk.RecordInfos, err error) {

	page := 1
	pagesize := 50

	for {
		tmpInfo, err := f.getAllDNSRecords(dnsID, page, pagesize)
		if err != nil {
			return infos, err
		}

		ll := len(tmpInfo.Data)

		infos.Data = append(infos.Data, tmpInfo.Data...)
		infos.Meta.Count = tmpInfo.Meta.Count

		if ll <= pagesize && ll > 0 {
			page++
			continue
		}
		break
	}

	return
}

func (f *FOneDNS) getAllDNSRecords(dnsID, page, pagesize int) (infos sdk.RecordInfos, err error) {

	uri := fmt.Sprintf("/dns/%d/record?page=%d&pagesize=%d", dnsID, page, pagesize)

	err = f.httpDo("GET", uri, nil, &infos)
	if err != nil {
		utils.Error("[fone] get dns error:", err)
	}

	return
}

// CreateDNSRecord create a new dns record
// 根据 recordType 的不同，填写对应的 ip，doamain，txt 值
// 当 recordType 类型为 MX 时，需要填写 priority
func (f *FOneDNS) CreateDNSRecord(dnsID int, subDomain string, recordType sdk.RecordType, value string, line, priority, ttl int, recordUnit sdk.RecordUnit) (recordid int, err error) {

	if dnsID <= 0 {
		err = ErrIDParams
		return
	}

	uri := fmt.Sprintf("/dns/%d/record", dnsID)

	p, err := f.checkRecord(subDomain, recordType, value, line, priority, ttl, recordUnit)
	if err != nil {
		return
	}

	var response struct {
		ID int `json:"id"`
	}

	err = f.httpDo("POST", uri, p, &response)
	if err != nil {
		utils.Error("[fone] new record error:", err)
		return
	}

	return response.ID, nil
}

// UpdateDNSRecord update a new dns record
// 根据 recordType 的不同，填写对应的 ip，doamain，txt 值
// 当 recordType 类型为 MX 时，需要填写 priority
func (f *FOneDNS) UpdateDNSRecord(dnsID, recordID int, subDomain string, recordType sdk.RecordType, value string, line, priority, ttl int, recordUnit sdk.RecordUnit) (err error) {

	if dnsID <= 0 || recordID <= 0 {
		err = ErrIDParams
		return
	}

	p, err := f.checkRecord(subDomain, recordType, value, line, priority, ttl, recordUnit)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("/dns/%d/record/%d", dnsID, recordID)

	err = f.httpDo("PUT", uri, p, nil)
	if err != nil {
		utils.Error("[fone] update record error:", err)
	}

	return nil
}

type RecordInfo struct {
	SubDomain string `json:"sub_domain"`
	Line      int    `json:"line"`
	RType     string `json:"type"`
	TTL       int    `json:"ttl"`
	Unit      string `json:"unit"`
	Txt       string `json:"text,omitempty"`
	IP        string `json:"ip,omitempty"`
	Domain    string `json:"domain,omitempty"`
	Priority  int    `json:"priority"`
}

func (f *FOneDNS) checkRecord(subDomain string, recordType sdk.RecordType, value string, line, priority, ttl int, recordUnit sdk.RecordUnit) (info RecordInfo, err error) {

	if subDomain == "" || len(subDomain) > 63 {
		err = fmt.Errorf("sub_domain error")
		return
	}

	if value == "" {
		err = fmt.Errorf("record miss value")
		return
	}

	ip := ""
	domain := ""
	txt := ""

	switch recordType {
	case sdk.RecordAType:
		if !IsOkIP(value) {
			err = fmt.Errorf("ip error")
			return
		}
		ip = value
	case sdk.RecordCNAMEType, sdk.RecordMXType, sdk.RecordNSType:
		if !IsOkDomain(value) {
			err = fmt.Errorf("domain error")
			return
		}
		domain = value
	case sdk.RecordTXTType:
		txt = value
	default:
		err = fmt.Errorf("not support type")
		return
	}

	info = RecordInfo{
		SubDomain: subDomain,
		Line:      line,
		RType:     recordType.String(),
		TTL:       ttl,
		Unit:      recordUnit.String(),
		IP:        ip,
		Domain:    domain,
		Txt:       txt,
		Priority:  1,
	}

	if recordType == sdk.RecordMXType {
		info.Priority = priority
	}

	return
}

// DeleteDNSRecord delete dns record
// 当重复删除 status=1 ，会报错
func (f *FOneDNS) DeleteDNSRecord(dnsID, recordID int) (err error) {

	if dnsID <= 0 || recordID <= 0 {
		err = ErrIDParams
		return
	}

	uri := fmt.Sprintf("/dns/%d/record/%d", dnsID, recordID)

	return f.httpDo("DELETE", uri, nil, nil)
}

// DisableDNSRecord disable or enable,
// disable true => disable
// disable false => enable
func (f *FOneDNS) DisableDNSRecord(dnsID, recordID int, disable bool) (err error) {

	if dnsID <= 0 || recordID <= 0 {
		err = ErrIDParams
		return
	}

	uri := fmt.Sprintf("/dns/%d/record/%d", dnsID, recordID)

	p := struct {
		ID       int  `json:"id"`
		Disabled bool `json:"disabled"`
	}{
		ID:       recordID,
		Disabled: disable,
	}

	return f.httpDo("PUT", uri, p, nil)
}
