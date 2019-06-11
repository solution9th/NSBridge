package sdk

import (
	"testing"
)

func TestGetDomainNS(t *testing.T) {

	n, err := GetDomainNS("dada.2fa.cc")
	t.Log(n, err)
}
