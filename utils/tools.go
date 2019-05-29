package utils

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/satori/go.uuid"
)

const (
	alphabetnum = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	tokenSize      = 16
	authKeySize    = 16
	authSecretSize = 32
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

// GenToken generate a token
func GenToken(size int) string {
	if size <= 0 {
		size = tokenSize
	}
	bts := genRandomBytes(size, alphabetnum)
	if bts[0] == '0' {
		bts[0] = 'x'
	}
	return string(bts)
}

func genRandomBytes(size int, base string) []byte {
	bts := make([]byte, size)
	n := len(base)
	// ignore error
	rand.Read(bts)
	for i, b := range bts {
		bts[i] = base[b%byte(n)]
	}
	return bts
}

func GenJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		Error("GenJson err:", err)
		return ""
	}
	return string(data)
}

func GenAuthKey() string {
	return uuid.NewV4().String()
}

func GenAuthSecret() string {
	bts := genRandomBytes(authSecretSize, alphabetnum)
	if bts[0] == '0' {
		bts[0] = 'm'
	}
	return string(bts)
}
