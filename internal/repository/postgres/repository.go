package postgres

import (
	"context"
	"fmt"
	"time"

	"example.com/taskservice/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ctxTimeout = 3 * time.Second
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, domain.NewInternal(domain.MessageInternal, fmt.Errorf("repository.BeginTx: %w", err))
	}

	return tx, nil
}
