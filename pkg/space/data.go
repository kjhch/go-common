package space

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"slices"
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

//------------------------------------------------------------------------------

type RedisDataType int

const (
	RedisString RedisDataType = iota + 1
	RedisList
	RedisSet
	RedisZSet
	RedisHash
	RedisJson
)

func NewRedisTemplate[T any](
	rdb *redis.Client,
	logger *Logger,
) *RedisTemplate[T] {
	return &RedisTemplate[T]{
		rdb:    rdb,
		logger: logger,
	}
}

type RedisTemplate[T any] struct {
	rdb    *redis.Client
	logger *Logger
}

type LoadArgs[T any] struct {
	DataType      RedisDataType
	Key           string
	Field         string
	Expiration    time.Duration
	BackToOrigin  func() (T, error)
	BackToOrigins func() ([]T, error)
}

func (rcv *RedisTemplate[T]) RenewOrLoadOrigin(ctx context.Context, args LoadArgs[T]) error {
	ex, err := rcv.rdb.Expire(ctx, args.Key, args.Expiration).Result()
	if err != nil {
		rcv.logger.Error("[cache]续约失败", "Key", args.Key, "err", err)
		return ErrCacheService
	}
	// cache hit
	if ex {
		return nil
	}

	// cache miss，回源
	_, _, err = rcv.LoadOrigin(ctx, args)
	return err
}

type GetOrLoadArgs[T any] struct {
	LoadArgs[T]
	ConvertCache func(string) (T, error)
}

func (rcv *RedisTemplate[T]) GetOrLoadOrigin(ctx context.Context, args GetOrLoadArgs[T]) (T, error) {
	var res T

	switch args.DataType {
	case RedisString:
		str, err := rcv.rdb.Get(ctx, args.Key).Result()
		// cache hit
		if err == nil {
			return args.ConvertCache(str)
		} else if !errors.Is(err, redis.Nil) {
			rcv.logger.Error("[cache]获取失败", "Key", args.Key, "err", err)
			return res, err
		}
	case RedisHash:
		str, err := rcv.rdb.HGet(ctx, args.Key, args.Field).Result()
		// cache hit
		if err == nil {
			return args.ConvertCache(str)
		} else if !errors.Is(err, redis.Nil) {
			rcv.logger.Error("[cache]获取失败", "Key", args.Key, "err", err)
			return res, err
		}
	case RedisJson:
		str, err := rcv.rdb.JSONGet(ctx, args.Key).Result()
		// cache hit
		if err == nil {
			return args.ConvertCache(str)
		} else if !errors.Is(err, redis.Nil) {
			rcv.logger.Error("[cache]获取失败", "Key", args.Key, "err", err)
			return res, err
		}
	default:
		return res, ErrUnsupportedCacheType
	}

	// cache miss，回源
	res, _, err := rcv.LoadOrigin(ctx, args.LoadArgs)
	return res, err
}

func (rcv *RedisTemplate[T]) LoadOrigin(ctx context.Context, args LoadArgs[T]) (origin T, origins []T, err error) {
	if args.BackToOrigin != nil {
		origin, err = args.BackToOrigin()
	}
	if args.BackToOrigins != nil {
		origins, err = args.BackToOrigins()
	}
	if err != nil {
		return
	}

	toAnySlice := func(items []T) []any {
		res := make([]any, len(items))
		for i, v := range items {
			res[i] = v
		}
		return res
	}
	// 构建缓存
	pipe := rcv.rdb.Pipeline()
	switch args.DataType {
	case RedisString:
		pipe.SetNX(ctx, args.Key, origin, args.Expiration)
	case RedisHash:
		pipe.HSetNX(ctx, args.Key, args.Field, origin)
	case RedisJson:
		pipe.JSONSetMode(ctx, args.Key, "$", origin, "nx")

	case RedisList:
		pipe.RPush(ctx, args.Key, toAnySlice(origins)...)
	case RedisSet:
		pipe.SAdd(ctx, args.Key, toAnySlice(origins)...)
	case RedisZSet:
		members := slices.Collect(func(yield func(redis.Z) bool) {
			for _, v := range origins {
				yield(any(v).(redis.Z))
			}
		})
		pipe.ZAddArgs(ctx, args.Key, redis.ZAddArgs{
			NX:      true,
			Members: members,
		})
	}
	pipe.Expire(ctx, args.Key, args.Expiration)
	_, err = pipe.Exec(ctx)
	if err != nil {
		rcv.logger.Error("[cache]构建失败", "Key", args.Key, "err", err)
		err = ErrCacheService
	}

	return
}
