package dr

import (
	"astor/backend/entity"
	"context"
)

type AssetRepository interface {
	Create(ctx context.Context, asset *entity.Assets) error
	Update(ctx context.Context, asset *entity.Assets) error
	Delete(ctx context.Context, asset *entity.Assets) error
	FetchAll(ctx context.Context) ([]entity.Assets, error)
	FetchByID(ctx context.Context, id uint) (*entity.Assets, error)
	FetchByCategory(ctx context.Context, categoryID uint) ([]entity.Assets, error)
	FindLastCode(ctx context.Context, prefix string) (string, error)
}
