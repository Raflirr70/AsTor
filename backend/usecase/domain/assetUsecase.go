package du

import (
	"astor/backend/entity"
	"context"
)

type AssetUsecase interface {
	Create(ctx context.Context, adminID uint, asset *entity.Assets, assetSpecs []entity.AssetSpecifications) error
	Update(ctx context.Context, adminID uint, asset *entity.Assets) error
	Checked(ctx context.Context, adminID uint, asset *entity.Assets) error
	Delete(ctx context.Context, adminID uint, asset *entity.Assets) error
	FetchAll(ctx context.Context) ([]entity.Assets, error)
	FetchByID(ctx context.Context, id uint) (*entity.Assets, error)
	UpdateSpecification(ctx context.Context, adminID uint, spec *entity.AssetSpecifications) error
}
