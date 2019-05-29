package sdk

import (
	"testing"
)

func TestCheckNS(t *testing.T) {

	// CheckDomainNS("baidu.com")
	// CheckDomainNS("dev.newio.cc")
}

func TestLookUp(t *testing.T) {

	// a := lookns("2fa.cc", "lv3ns1.ffdns.net.")
	// fmt.Println(a)
}

func TestGetDomainNS(t *testing.T) {

	n, err := GetDomainNS("dada.2fa.cc")
	t.Log(n, err)
}
