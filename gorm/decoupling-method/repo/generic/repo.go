package repo

import (
	"context"
	"go-tricks/gorm/decoupling-method/model"
)

type GenericRepo[T any, PT model.PointerModel[T]] interface {
	Create(ctx context.Context, ptrModel PT) error
	CreateInBatches(ctx context.Context, ptrModels []PT, batchSize int) error
	GetByID(ctx context.Context, id uint64) (PT, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]PT, error)
	GetByPage(ctx context.Context, page, pageSize int) ([]PT, error)
	GetByCursor(ctx context.Context, cursor uint64, pageSize int) ([]PT, uint64, bool, error)
	Update(ctx context.Context, ptrModel PT) error
	DeleteByID(ctx context.Context, id uint64) error
	DeleteByIDs(ctx context.Context, ids []uint64) error
}
