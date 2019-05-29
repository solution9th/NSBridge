package utils

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// IsOkDomain check domain
func IsOkDomain(domain string) bool {

	if domain == "" || len(strings.Replace(domain, ".", "", -1)) > 255 {
		return false
	}

	re := `^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})*(\.[a-zA-Z0-9]{0,62}){1}\.?$`

	return Match(re, domain)
}

// Match 正则匹配
func Match(re, s string) bool {
	ok, err := regexp.MatchString(re, s)
	if err != nil {
		Error("regexp error:", err)
		return false
	}
	if ok {
		return true
	}
	return false
}

// IsExist 判断是否存在
//
// 注意 a 必须是 list 中的类型，别名也不行
func IsExist(a interface{}, list interface{}) bool {

	r := interfaceSlice(list)

	for _, v := range r {
		if reflect.DeepEqual(a, v) {
			return true
		}
	}

	return false
}

func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil
		// panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func GetCacheTakeOverKey(key string) string {
	return fmt.Sprintf("takeover:%s", key)
}

func FindFile(file ...string) (string, bool) {
	for _, v := range file {
		if Exists(v) {
			return v, true
		}
	}
	return "", false
}

// Exists checks if a file or directory exists.
func Exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		return false
	}
	return false
}
