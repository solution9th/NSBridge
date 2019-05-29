package v1

import (
	"testing"
)

func TestCheck(t *testing.T) {

	tests := []struct {
		s    string
		want bool
	}{
		{"", false},
		{"sss", true},
		{"@", true},
		{"*", true},
		{"-", false},
		{"a-a", true},
		{"a-", false},
		{"-a", false},
		{"a.a", false},
		{"a**a", false},
		{"000a**00000a", false},
	}

	for k, v := range tests {
		got := check(v.s)
		if got != v.want {
			t.Errorf("k: %v, got: %v, want: %v", k, got, v.want)
		}
	}
}
