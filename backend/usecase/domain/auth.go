package du

import (
	"astor/backend/entity"
	"context"
)

type AuthUsecase interface {
	Register(ctx context.Context, user *entity.Users) error
	Login(ctx context.Context, username, password string) (*entity.Users, string, error)
	Logout(ctx context.Context, token string) error
	GetProfile(ctx context.Context, id uint) (*entity.Users, error)
}
