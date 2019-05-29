package utils

import (
	"fmt"
	"testing"

	"github.com/solution9th/NSBridge/sdk"
)

func TestGenToken(t *testing.T) {

	for i := 0; i < 5; i++ {
		fmt.Println(GenToken(10))
	}
}

func TestIsExist(t *testing.T) {

	tests := []struct {
		a    interface{}
		b    interface{}
		want bool
	}{
		{"a", []interface{}{"a", "b"}, true},
		{"c", []interface{}{"a", "b"}, false},
		{sdk.RecordAType, sdk.AllowTypeList, true},
		{"A", sdk.AllowTypeList, true},
		{"MMMM", sdk.AllowTypeList, false},
	}

	for k, v := range tests {
		got := IsExist(v.a, v.b)
		if got != v.want {
			t.Errorf("k: %v got: %v want: %v", k, got, v.want)
		}
	}
}
