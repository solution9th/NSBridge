package fone

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	c *cache.Cache
)

func init() {
	c = cache.New(5*time.Minute, 10*time.Minute)
}
