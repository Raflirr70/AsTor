package du

import (
	"astor/backend/entity"
	"context"
)

type LaboratoryUsecase interface {
	Create(ctx context.Context, admin uint, laboratory *entity.Laboratories) error
	Update(ctx context.Context, admin uint, laboratory *entity.Laboratories) error
	Delete(ctx context.Context, admin uint, laboratory *entity.Laboratories) error
	Restore(ctx context.Context, admin uint, laboratory *entity.Laboratories) error
	GetAll(ctx context.Context) ([]entity.Laboratories, error)
}
