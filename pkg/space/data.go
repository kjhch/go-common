package space

import (
	"fmt"
	"time"
)

type CacheKey interface {
	Format(args ...any) string
	Expiration() time.Duration
}

func NewCacheKey(template string, exp time.Duration) CacheKey {
	return &cacheKey{
		template:   template,
		expiration: exp,
	}
}

type cacheKey struct {
	template   string
	expiration time.Duration
}

func (k *cacheKey) Format(args ...any) string {
	return fmt.Sprintf(k.template, args...)
}
func (k *cacheKey) Expiration() time.Duration {
	return k.expiration
}
