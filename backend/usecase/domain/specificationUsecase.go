package du

import (
	"astor/backend/entity"
	"context"
)

type SpecificationUsecase interface {
	Create(ctx context.Context, adminID uint, specification *entity.Specifications) error
	Update(ctx context.Context, adminID uint, specification *entity.Specifications) error
	GetAll(ctx context.Context) ([]entity.Specifications, error)
	GetByID(ctx context.Context, id uint) (*entity.Specifications, error)
	Delete(ctx context.Context, adminID uint, specification *entity.Specifications) error
}
