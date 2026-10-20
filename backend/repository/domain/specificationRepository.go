package dr

import (
	"astor/backend/entity"
	"context"
)

type SpecificationRepository interface {
	Create(ctx context.Context, specification *entity.Specifications) error
	Update(ctx context.Context, specification *entity.Specifications) error
	FetchAll(ctx context.Context) ([]entity.Specifications, error)
	FetchByID(ctx context.Context, id uint) (*entity.Specifications, error)
	Delete(ctx context.Context, specification *entity.Specifications) error
}
