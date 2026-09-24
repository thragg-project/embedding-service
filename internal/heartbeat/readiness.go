package heartbeat

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ReadinessHeartbeat struct {
		pool *pgxpool.Pool
	}
)

func NewReadinessHeartbeat(pool *pgxpool.Pool) *ReadinessHeartbeat {
	return &ReadinessHeartbeat{
		pool: pool,
	}
}

func (h *ReadinessHeartbeat) Beat(ctx context.Context) error {
	if err := h.pool.Ping(ctx); err != nil {
		return err
	}
	return nil
}
