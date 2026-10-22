package dr

import (
	"astor/backend/entity"
	"context"
)

type ComputerRepository interface {
	Create(ctx context.Context, computer *entity.Computers) error
	Update(ctx context.Context, computer *entity.Computers) error
	Delete(ctx context.Context, computer *entity.Computers) error
	FetchAll(ctx context.Context) ([]entity.Computers, error)
	FetchByID(ctx context.Context, id uint) (*entity.Computers, error)
	FetchByLocation(ctx context.Context, locationID uint) ([]entity.Computers, error)
}

type ComputerComponentRepository interface {
	Create(ctx context.Context, component *entity.ComputerComponents) error
	Update(ctx context.Context, component *entity.ComputerComponents) error
	Delete(ctx context.Context, component *entity.ComputerComponents) error
	FetchByComputerID(ctx context.Context, computerID uint) ([]entity.ComputerComponents, error)
	CheckUsed(ctx context.Context, AssetID uint) (bool, error)
}
