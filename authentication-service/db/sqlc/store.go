package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
}

type PgStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &PgStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *PgStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}

	db := New(tx)

	err = fn(db)

	if err != nil {

		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
