package txs

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transactor interface {
	WithTransaction(context.Context, func(ctx context.Context) error) error
}

type QueryPerformer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func GetQueryPerformer(ctx context.Context, defaultPerformer QueryPerformer) QueryPerformer {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}

	return defaultPerformer
}

type TxBeginner struct {
	db *pgxpool.Pool
}

type txKey struct{}

func NewTxBeginner(db *pgxpool.Pool) Transactor {
	return &TxBeginner{db: db}
}

func injectCtx(ctx context.Context, tx pgx.Tx) context.Context {
	ctx = context.WithValue(ctx, txKey{}, tx)

	return ctx
}

func (tb *TxBeginner) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tb.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("TxBeginner: failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	err = fn(injectCtx(ctx, tx))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
