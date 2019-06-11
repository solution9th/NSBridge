package cache

import (
	"fmt"

	"github.com/solution9th/NSBridge/internal/utils"

	"github.com/go-redis/redis"
)

// InitRedis 初始化 redis
func InitRedis(host, port, passwd string, dataBase int) error {
	if host == "" || port == "" {
		return fmt.Errorf("host or port is nil")
	}
	DefaultCache = &Cache{
		Cache: redis.NewClient(&redis.Options{
			Addr:       host + ":" + port,
			Password:   passwd,   // no password set
			DB:         dataBase, // use default DB
			MaxRetries: 2,
		}),
	}

	pong, err := DefaultCache.Cache.Ping().Result()
	if err != nil {
		return err
	}

	utils.Infof("PING redis: %v", pong)

	return nil
}
