package space

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"slices"
	"sync"
	"time"
)

type Transformer[S any] interface {
	Transform() (S, error)
	Restore(S) error
}

type RedisCache[T Transformer[S], S any] struct {
	rdb *redis.Client

	CacheKey       *CacheKey
	NewTransformer func() T

	Codec RedisCodec[S]
}

func (rc *RedisCache[T, S]) GetOrLoad(ctx context.Context, origin func() (T, error), args ...any) (T, error) {
	key := rc.CacheKey.Format(args...)

	cache, err := rc.Codec.get(ctx, rc.rdb, key)
	// cache hit
	if err == nil {
		var res = rc.NewTransformer()
		err = res.Restore(cache)
		return res, err
	} else if !errors.Is(err, redis.Nil) {
		var res T
		return res, err
	}

	// cache miss，回源构建缓存
	return rc.loadCache(ctx, key, origin)
}

func (rc *RedisCache[T, S]) loadCache(ctx context.Context, key string, origin func() (T, error)) (T, error) {
	res, err := origin()
	if err != nil {
		return res, err
	}
	to, err := res.Transform()
	if err != nil {
		return res, err
	}

	// 构建缓存
	err = rc.Codec.set(ctx, rc.rdb, key, to, rc.CacheKey.Expiration)
	return res, err
}

func (rc *RedisCache[T, S]) BatchGetOrLoad(ctx context.Context, args []struct {
	FormatArgs []any
	Origin     func() (T, error)
}) ([]T, error) {
	var wg sync.WaitGroup
	resChan := make(chan struct {
		idx int
		t   T
		err error
	}, len(args))
	for i, arg := range args {
		wg.Go(func() {
			r, err := rc.GetOrLoad(ctx, arg.Origin, arg.FormatArgs...)
			resChan <- struct {
				idx int
				t   T
				err error
			}{idx: i, t: r, err: err}
		})
	}
	wg.Wait()
	close(resChan)

	results := make([]T, len(args))
	for res := range resChan {
		if res.err != nil {
			return nil, res.err
		}
		results[res.idx] = res.t
	}
	return results, nil
}

func (rc *RedisCache[T, S]) RenewOrLoad(ctx context.Context, origin func() (T, error), args ...any) error {
	key := rc.CacheKey.Format(args...)
	ex, err := rc.rdb.Expire(ctx, key, rc.CacheKey.Expiration).Result()
	if err != nil {
		return ErrCacheService
	}

	// cache hit
	if ex {
		return nil
	}

	// cache miss，回源构建缓存
	_, err = rc.loadCache(ctx, key, origin)
	return err
}

func (rc *RedisCache[T, S]) BatchRenewOrLoad(ctx context.Context, args []struct {
	FormatArgs []any
	Origin     func() (T, error)
}) error {
	pipe := rc.rdb.Pipeline()
	for _, arg := range args {
		key := rc.CacheKey.Format(arg.FormatArgs...)
		pipe.Expire(ctx, key, rc.CacheKey.Expiration)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(args))
	for i, cmd := range cmds {
		ex := cmd.(*redis.BoolCmd).Val()
		// cache miss，回源构建缓存
		if !ex {
			wg.Go(func() {
				_, err = rc.loadCache(ctx, rc.CacheKey.Format(args[i].FormatArgs...), args[i].Origin)
				errChan <- err
			})
		}
	}
	wg.Wait()
	close(errChan)

	for err = range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (rc *RedisCache[T, S]) Delete(ctx context.Context, args ...any) (int, error) {
	key := rc.CacheKey.Format(args...)
	result, err := rc.rdb.Del(ctx, key).Result()
	return int(result), err
}

func (rc *RedisCache[T, S]) BatchDelete(ctx context.Context, args [][]any) (int, error) {
	keys := slices.Collect(func(yield func(string) bool) {
		for _, arg := range args {
			yield(rc.CacheKey.Format(arg...))
		}
	})
	result, err := rc.rdb.Del(ctx, keys...).Result()
	return int(result), err
}

//------------------------------------------------------------------------------

func NewRedisStringCache[T Transformer[string]](
	rdb *redis.Client,
	CacheKey *CacheKey,
	NewTransformer func() T,
) *RedisCache[T, string] {
	return &RedisCache[T, string]{
		rdb:            rdb,
		CacheKey:       CacheKey,
		NewTransformer: NewTransformer,

		Codec: RedisStringCodec,
	}
}

func NewRedisJsonCache[T Transformer[string]](
	rdb *redis.Client,
	CacheKey *CacheKey,
	NewTransformer func() T,
) *RedisCache[T, string] {
	return &RedisCache[T, string]{
		rdb:            rdb,
		CacheKey:       CacheKey,
		NewTransformer: NewTransformer,

		Codec: RedisJsonCodec,
	}
}

func NewRedisHashCache[T Transformer[map[string]string]](
	rdb *redis.Client,
	CacheKey *CacheKey,
	NewTransformer func() T,
) *RedisCache[T, map[string]string] {
	return &RedisCache[T, map[string]string]{
		rdb:            rdb,
		CacheKey:       CacheKey,
		NewTransformer: NewTransformer,

		Codec: RedisHashCodec,
	}
}

func NewRedisZsetCache[T Transformer[[]redis.Z]](
	rdb *redis.Client,
	CacheKey *CacheKey,
	NewTransformer func() T,
) *RedisCache[T, []redis.Z] {
	return &RedisCache[T, []redis.Z]{
		rdb:            rdb,
		CacheKey:       CacheKey,
		NewTransformer: NewTransformer,

		Codec: RedisZsetCodec,
	}
}

//------------------------------------------------------------------------------

type RedisCodec[S any] struct {
	get func(ctx context.Context, rdb *redis.Client, key string) (S, error)
	set func(ctx context.Context, rdb *redis.Client, key string, value S, expiration time.Duration) error
}

var RedisStringCodec = RedisCodec[string]{
	get: func(ctx context.Context, rdb *redis.Client, key string) (string, error) {
		return rdb.Get(ctx, key).Result()
	},
	set: func(ctx context.Context, rdb *redis.Client, key string, value string, expiration time.Duration) error {
		return rdb.Set(ctx, key, value, expiration).Err()
	},
}

var RedisJsonCodec = RedisCodec[string]{
	get: func(ctx context.Context, rdb *redis.Client, key string) (string, error) {
		return rdb.JSONGet(ctx, key).Result()
	},
	set: func(ctx context.Context, rdb *redis.Client, key string, value string, expiration time.Duration) error {
		pipe := rdb.Pipeline()
		pipe.JSONSet(ctx, key, "$", value)
		pipe.Expire(ctx, key, expiration)
		_, err := pipe.Exec(ctx)
		return err
	},
}

var RedisHashCodec = RedisCodec[map[string]string]{
	get: func(ctx context.Context, rdb *redis.Client, key string) (map[string]string, error) {
		result, err := rdb.HGetAll(ctx, key).Result()
		if err != nil {
			return result, err
		}
		if len(result) == 0 {
			return result, redis.Nil
		}
		return result, nil
	},
	set: func(ctx context.Context, rdb *redis.Client, key string, value map[string]string, expiration time.Duration) error {
		pipe := rdb.Pipeline()
		pipe.HSet(ctx, key, value)
		pipe.Expire(ctx, key, expiration)
		_, err := pipe.Exec(ctx)
		return err
	},
}

var RedisListCodec = RedisCodec[[]string]{
	get: func(ctx context.Context, rdb *redis.Client, key string) ([]string, error) {
		result, err := rdb.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			return result, err
		}
		if len(result) == 0 {
			return result, redis.Nil
		}
		return result, err
	},
	set: func(ctx context.Context, rdb *redis.Client, key string, value []string, expiration time.Duration) error {
		pipe := rdb.Pipeline()
		pipe.RPush(ctx, key, slices.Collect(func(yield func(any) bool) {
			for _, v := range value {
				yield(v)
			}
		})...)
		pipe.Expire(ctx, key, expiration)
		_, err := pipe.Exec(ctx)
		return err
	},
}

var RedisSetCodec = RedisCodec[[]string]{
	get: func(ctx context.Context, rdb *redis.Client, key string) ([]string, error) {
		result, err := rdb.SMembers(ctx, key).Result()
		if err != nil {
			return result, err
		}
		if len(result) == 0 {
			return result, redis.Nil
		}
		return result, err
	},
	set: func(ctx context.Context, rdb *redis.Client, key string, value []string, expiration time.Duration) error {
		pipe := rdb.Pipeline()
		pipe.SAdd(ctx, key, slices.Collect(func(yield func(any) bool) {
			for _, v := range value {
				yield(v)
			}
		})...)
		pipe.Expire(ctx, key, expiration)
		_, err := pipe.Exec(ctx)
		return err
	},
}

var RedisZsetCodec = RedisCodec[[]redis.Z]{
	get: func(ctx context.Context, rdb *redis.Client, key string) ([]redis.Z, error) {
		result, err := rdb.ZRangeWithScores(ctx, key, 0, -1).Result()
		if err != nil {
			return result, err
		}
		if len(result) == 0 {
			return result, redis.Nil
		}
		return result, err
	},
	set: func(ctx context.Context, rdb *redis.Client, key string, value []redis.Z, expiration time.Duration) error {
		pipe := rdb.Pipeline()
		pipe.ZAdd(ctx, key, value...)
		pipe.Expire(ctx, key, expiration)
		_, err := pipe.Exec(ctx)
		return err
	},
}
