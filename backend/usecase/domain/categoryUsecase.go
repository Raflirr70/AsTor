package du

import (
	"astor/backend/entity"
	"context"
)

type CategoryUsecase interface {
	Create(ctx context.Context, adminID uint, category *entity.Categories, specificationIDs []uint) error
	Update(ctx context.Context, adminID uint, category *entity.Categories, specificationIDs []uint) error
	Delete(ctx context.Context, adminID uint, category *entity.Categories) error
	GetAll(ctx context.Context) ([]entity.Categories, error)
	GetByID(ctx context.Context, id uint) (*entity.Categories, error)
}
