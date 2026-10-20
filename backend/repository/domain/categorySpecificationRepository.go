package dr

import (
	"astor/backend/entity"
	"context"
)

type CategorySpecificationRepository interface {
	Create(ctx context.Context, cs *entity.CategorySpecifications) error
	DeleteByCategoryID(ctx context.Context, categoryID uint) error
	FetchByCategoryID(ctx context.Context, categoryID uint) ([]entity.CategorySpecifications, error)
}
