package utils

import "fmt"

// GetSecretKey key 在 redis 中的缓存 key
func GetSecretKey(key string) string {
	return fmt.Sprintf("mid-key-%v", key)
}
