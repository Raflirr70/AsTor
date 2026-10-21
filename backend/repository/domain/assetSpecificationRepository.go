package dr

import (
	"astor/backend/entity"
	"context"
)

type AssetSpecificationRepository interface {
	Create(ctx context.Context, as *entity.AssetSpecifications) error
	Update(ctx context.Context, as *entity.AssetSpecifications) error
	Delete(ctx context.Context, as *entity.AssetSpecifications) error
	FetchByID(ctx context.Context, id uint) (*entity.AssetSpecifications, error)
	FetchByAssetID(ctx context.Context, assetID uint) ([]entity.AssetSpecifications, error)
}
