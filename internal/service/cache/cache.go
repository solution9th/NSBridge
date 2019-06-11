package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/solution9th/NSBridge/internal/utils"

	"github.com/go-redis/redis"
)

// RedisNotFound not found value
var RedisNotFound = redis.Nil

// DefaultCache default cache
var DefaultCache *Cache

type Cache struct {
	Cache *redis.Client
}

func (c *Cache) Set(key string, value interface{}, timeout time.Duration) error {
	g, err := GobEnValue(value)
	if err != nil {
		utils.Error("cache set error:", err)
		return err
	}
	return c.Cache.Set(key, g, timeout).Err()
}

func (c *Cache) Get(key string, p interface{}) error {
	r, err := c.Cache.Get(key).Result()
	if err == redis.Nil {
		return RedisNotFound
	} else if err != nil {
		return err
	}
	return GobDeValue([]byte(r), p)
}

func (c *Cache) Delete(key string) error {
	return c.Cache.Del(key).Err()
}

// Expire 重新设定一个 key 的过期时间
func (c *Cache) Expire(key string, timeout time.Duration) error {
	return c.Cache.Expire(key, timeout).Err()
}

func (c *Cache) Exist(key string) (bool, error) {
	r := c.Cache.Exists(key)
	if r.Err() != nil {
		return false, r.Err()
	}

	if r.Val() == 1 {
		return true, r.Err()
	}
	return false, r.Err()
}

func (c *Cache) Incr(key string) int64 {
	id, err := c.Cache.Incr(key).Result()
	if err != nil {
		return -1
	}
	return id
}

// 设置成功则标记true，设置失败或者已经存在则返回false
func (c *Cache) SetNX(key string, value interface{}, timeout time.Duration) bool {
	return c.Cache.SetNX(key, value, timeout).Val()
}

func (c *Cache) SetNoGob(key, value string, timeout time.Duration) error {
	return c.Cache.Set(key, value, timeout).Err()
}

func (c *Cache) GetNoGob(key string) (string, error) {
	r, err := c.Cache.Get(key).Result()
	if err == redis.Nil {
		return "", RedisNotFound
	} else if err != nil {
		return "", err
	}
	return r, nil
}

func GobEnValue(value interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(value)
	if err != nil {
		return []byte(""), err
	}
	return b.Bytes(), nil
}

func GobDeValue(data []byte, p interface{}) error {
	var b = bytes.NewBuffer(data)
	d := gob.NewDecoder(b)
	return d.Decode(p)
}
