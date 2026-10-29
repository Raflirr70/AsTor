package dr

import (
	"astor/backend/entity"
	"context"
)

type CategorySpecificationRepository interface {
	Create(ctx context.Context, cs *entity.CategorySpecifications) error
	DeleteByCategoryID(ctx context.Context, categoryID uint) error
	DeleteByCategoryAndSpecificationID(ctx context.Context, categoryID uint, specificationID uint) error
	FetchByCategoryID(ctx context.Context, categoryID uint) ([]entity.CategorySpecifications, error)
	ExistsByCategoryAndSpecificationID(ctx context.Context, categoryID uint, specificationID uint) (bool, error)
}
