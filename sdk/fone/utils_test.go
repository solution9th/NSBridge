package fone

import "testing"

func TestIsOkDomain(t *testing.T) {

	tests := []struct {
		Domain string
		Result bool
	}{
		{"baidu.com", true},
		{"1.2.a-b.baidu.com", true},
		{"rrrr.baidu.com", true},
		{"-.baidu.com", false},
		{"a.a.a.baidu.com", true},
		{"uuuuu", false},
		{"uuu_uu.com", false},
		{"uuuuu.c_om", false},
		{"uuuuu.c-om", false},
		{".uuuuu.baidu.com", false},
	}

	for _, test := range tests {

		if got := IsOkDomain(test.Domain); got != test.Result {
			t.Errorf("domain: %v, got: %v, want: %v", test.Domain, got, test.Result)
		}
	}
}

func TestIsOkIP(t *testing.T) {

	tests := []struct {
		IP     string
		Result bool
	}{
		{"8.8.8.8", true},
		{"888.8.8.8", false},
		{"22..8.8", false},
		{"22.-.8.8", false},
		{"22.8.8", false},
		{"8.8.8.8.0.0", false},
		{"a.8.8.8", false},
		{"uuuuu", false},
		{"", false},
	}

	for _, test := range tests {

		if got := IsOkIP(test.IP); got != test.Result {
			t.Errorf("ip: %v, got: %v, want: %v", test.IP, got, test.Result)
		}
	}
}
