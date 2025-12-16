package space

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

//------------------------------------------------------------------------------

const txKey = "spaceRepoTx"

type TxAble[T any] interface {
	WithTx(tx pgx.Tx) *T
}

type UnitOfWork[T any] struct {
	logger *Logger
	pool   *pgxpool.Pool
	db     TxAble[T]
}

func NewUnitOfWork[T any](
	logger *Logger,
	pool *pgxpool.Pool,
	db TxAble[T],
) *UnitOfWork[T] {
	return &UnitOfWork[T]{
		logger: logger,
		pool:   pool,
		db:     db,
	}
}

func (uow *UnitOfWork[T]) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// 已存在事务：加入
	if dbqTx := ctx.Value(txKey); dbqTx != nil {
		uow.logger.Info("[db]已存在事务，加入", "qtx", fmt.Sprintf("%p", dbqTx))
		return fn(ctx)
	}

	// 不存在事务：开启
	tx, err := uow.pool.Begin(ctx)
	if err != nil {
		uow.logger.Error("[db]事务开启失败", "err", err)
		return err
	}
	defer tx.Rollback(ctx)
	qtx := uow.db.WithTx(tx)
	uow.logger.Info("[db]开启事务", "qtx", fmt.Sprintf("%p", qtx))

	err = fn(context.WithValue(ctx, txKey, qtx))

	if err != nil {
		uow.logger.Info("[db]事务失败，回滚", "qtx", fmt.Sprintf("%p", qtx))
		return err
	}
	uow.logger.Info("[db]提交事务", "qtx", fmt.Sprintf("%p", qtx))
	return tx.Commit(ctx)
}

func GetDbq[T any](ctx context.Context, dbq *T) *T {
	if dbqTx := ctx.Value(txKey); dbqTx != nil {
		dbq = dbqTx.(*T)
	}
	return dbq
}
