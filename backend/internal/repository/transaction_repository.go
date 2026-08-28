package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	db *pgxpool.Pool
}

func NewTxManager(
	db *pgxpool.Pool,
) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (m *TxManager) Begin(
	ctx context.Context,
) (pgx.Tx, error) {

	return m.db.Begin(ctx)
}