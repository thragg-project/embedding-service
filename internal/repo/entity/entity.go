package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/fr33dman/go-template/internal/core"
	apperrors "github.com/fr33dman/go-template/internal/errors"
	"github.com/fr33dman/go-template/internal/repo/entity/sqlc"
	"github.com/fr33dman/go-template/pkg/database"
)

type EntityRepository struct {
	dbtx sqlc.DBTX
}

func NewEntityRepository(dbtx sqlc.DBTX) *EntityRepository {
	return &EntityRepository{dbtx: dbtx}
}

func (r *EntityRepository) FetchOne(ctx context.Context, id int64) (core.Entity, error) {
	row, err := r.queries(ctx).GetEntity(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.Entity{}, apperrors.ErrNotFound
		}
		return core.Entity{}, err
	}
	return mapEntity(row), nil
}

func (r *EntityRepository) Insert(ctx context.Context, entity core.Entity) (core.Entity, error) {
	row, err := r.queries(ctx).CreateEntity(ctx, sqlc.CreateEntityParams{
		Field1: entity.Field1,
		Field2: int32(entity.Field2),
	})
	if err != nil {
		return core.Entity{}, err
	}

	return mapEntity(row), nil
}

func (r *EntityRepository) queries(ctx context.Context) *sqlc.Queries {
	return sqlc.New(database.DBTX(ctx, r.dbtx))
}

func mapEntity(row sqlc.Entity) core.Entity {
	return core.Entity{
		Id:     &row.ID,
		Field1: row.Field1,
		Field2: int(row.Field2),
	}
}
