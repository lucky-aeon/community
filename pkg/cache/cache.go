package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var c *cache.Cache

func Init() {
	c = cache.New(5*time.Minute, 10*time.Minute)
}

func GetInstance() *cache.Cache {
	return c
}

func CountLimit(key string, limit int, ttl time.Duration) bool {
	v, e := c.Get(key)
	if !e {
		v = 0
	}
	i := v.(int)
	i++
	if i == limit {
		return false
	}
	c.Set(key, i, ttl)
	return true
}
