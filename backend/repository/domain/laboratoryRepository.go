package dr

import (
	"astor/backend/entity"
	"context"
)

type LaboratoryRepository interface {
	Create(ctx context.Context, laboratory *entity.Laboratories) error
	Update(ctx context.Context, laboratory *entity.Laboratories) error
	Delete(ctx context.Context, laboratory *entity.Laboratories) error
	FetchAll(ctx context.Context) ([]entity.Laboratories, error)
	FetchByID(ctx context.Context, id uint) (*entity.Laboratories, error)
	FetchByLocation(ctx context.Context, locationID uint) (*entity.Laboratories, error)
	Restore(ctx context.Context, id uint) error
	FetchDeletedByID(ctx context.Context, id uint) (*entity.Laboratories, error)
}
