package space

import (
	"fmt"
	"time"
)

type CacheKey struct {
	Template   string
	Expiration time.Duration
}

func NewCacheKey(template string, exp time.Duration) *CacheKey {
	return &CacheKey{
		Template:   template,
		Expiration: exp,
	}
}

func (k *CacheKey) Format(args ...any) string {
	return fmt.Sprintf(k.Template, args...)
}

//------------------------------------------------------------------------------

type OrderBy string

const (
	OrderAsc  OrderBy = "ASC"
	OrderDesc OrderBy = "DESC"
)

type CursorPagination struct {
	LastIdx int64
	OrderBy OrderBy
	Limit   int
}
