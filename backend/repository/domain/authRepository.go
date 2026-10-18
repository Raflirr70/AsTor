package dr

import (
	"astor/backend/entity"
	"context"
)

type AuthRepository interface {
	Create(ctx context.Context, user *entity.Users) error
	Update(ctx context.Context, user *entity.Users) error
	Delete(ctx context.Context, user *entity.Users) error
	FindByID(ctx context.Context, id uint) (*entity.Users, error)
	FindByUsername(ctx context.Context, username string) (*entity.Users, error)
	FindByNPM(ctx context.Context, npm string) (*entity.Users, error)
	FindAll(ctx context.Context) ([]entity.Users, error)
}
