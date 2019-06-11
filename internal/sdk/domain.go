package sdk

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/haozibi/dnsutil"

	"github.com/miekg/dns"
	dnsMsg "github.com/miekg/dns"
	"github.com/robfig/cron"
)

// CheckDomainNS 定时检查检查某个域名的 ns 是否为指定 ns
// func CheckDomainNS2(domain string, ns []string) {

// 	checkDomainNS(domain, ns, 10*time.Second, func(domain string) error {
// 		fmt.Println("SAVE", domain)
// 		return nil
// 	})
// }

// CheckDomainNS 定期检查域名的 ns，如果是指定 ns 则执行 f func(string) error
func CheckDomainNS(domain string, ns []string, timeout time.Duration, f func(string) error) {

	c := cron.New()
	exit := make(chan bool)

	c.AddFunc("*/2 * * * * ?", func() {
		if LookNS(domain, ns) {
			err := f(domain)
			if err != nil {
				fmt.Println("f error:", err)
				return
			}
			exit <- true
		}
	})
	c.Start()

	select {
	case <-time.After(10 * time.Second):
		// fmt.Println("not find an hour")
	case <-exit:
		// fmt.Println("ok")
	}
	c.Stop()

	return
}

// LookNS 检查域名的 ns 是否是指定 ns
func LookNS(domain string, nameServer []string) bool {
	ns, err := GetDomainNS(domain)
	if err != nil {
		log.Println("get domain ns error:", err)
		return false
	}

	for _, a := range ns {
		for _, b := range nameServer {
			if a == b {
				return true
			}
		}
	}

	return false
}

// GetDomainNS 获得域名的 ns 服务器
//
// 原理： dig +trace 获得其最权威的 ns
func GetDomainNS(domain string) (ns []string, err error) {

	d := dnsutil.New()

	resp, err := d.TraceForRecord(domain, dnsMsg.TypeNS)
	if err != nil {
		return nil, err
	}

	for _, r := range resp {
		if r.Msg.Authoritative {
			for _, answer := range r.Msg.Answer {
				ans, ok := answer.(*dns.NS)
				if !ok {
					continue
				}
				s := strings.TrimSuffix(ans.Ns, ".")
				ns = append(ns, s)
			}
		}
	}

	return ns, nil

}
