package dr

import (
	"astor/backend/entity"
	"context"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Categories) error
	Update(ctx context.Context, category *entity.Categories) error
	Delete(ctx context.Context, category *entity.Categories) error
	FetchAll(ctx context.Context) ([]entity.Categories, error)
	FetchByID(ctx context.Context, id uint) (*entity.Categories, error)
}
