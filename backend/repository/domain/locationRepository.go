package dr

import (
	"astor/backend/entity"
	"context"
)

type LocationRepository interface {
	Create(ctx context.Context, location *entity.Locations) error
	Update(ctx context.Context, location *entity.Locations) error
	Delete(ctx context.Context, location *entity.Locations) error
	FetchAll(ctx context.Context) ([]entity.Locations, error)
	FetchByID(ctx context.Context, id uint) (*entity.Locations, error)
	Restore(ctx context.Context, id uint) error
}
