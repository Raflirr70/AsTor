package usecase

import (
	"astor/backend/config"
	"astor/backend/entity"
	"astor/backend/pkg/jwt"
	dr "astor/backend/repository/domain"
	du "astor/backend/usecase/domain"
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authUsecaseImpl struct {
	authRepo dr.AuthRepository
	JWT      config.JWTConfig
	Redis    *redis.Client
}

func NewAuthUsecaseImpl(authRepo dr.AuthRepository) du.AuthUsecase {
	return &authUsecaseImpl{authRepo: authRepo}
}
func (u *authUsecaseImpl) Register(ctx context.Context, user *entity.Users) error {
	exist, err := u.authRepo.FindByNPM(ctx, user.Npm)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check NPM: %w", err)
	}
	if exist != nil {
		return entity.ErrAlreadyExists
	}
	//validasi password
	if len(user.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	//validasi npm
	matched, _ := regexp.MatchString(`^[0-9]{10}$`, user.Npm)
	if !matched {
		return fmt.Errorf("NPM must contain exactly 10 digits")
	}
	return u.authRepo.Create(ctx, user)
}
func (u *authUsecaseImpl) Login(ctx context.Context, username, password string) (*entity.Users, string, error) {
	exist, err := u.authRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, "", fmt.Errorf("invalid Credential")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("invalid Credential")
	}

	token, err := jwt.GenerateToken(exist.ID, exist.Username, exist.Role.Name, u.JWT.Secret, u.JWT.Expiration)
	if err != nil {
		return nil, "", err
	}

	redisKey := fmt.Sprintf("sessions:%d", exist.ID)
	u.Redis.Set(ctx, redisKey, token, u.JWT.Expiration)

	return exist, token, nil
}
func (u *authUsecaseImpl) Logout(ctx context.Context, token string) error {
	claim, err := jwt.ParseToken(token, u.JWT.Secret)
	if err != nil {
		return errors.New("Invalid Token")
	}
	redisKey := fmt.Sprintf("sessions:%d", claim.UserID)
	u.Redis.Del(ctx, redisKey)

	ttl := time.Until(claim.ExpiresAt.Time)
	if ttl > 0 {
		blacklistKey := fmt.Sprintf("blacklist:%s", token)
		u.Redis.Set(ctx, blacklistKey, "1", ttl)
	}
	return nil
}
func (u *authUsecaseImpl) GetProfile(ctx context.Context, id uint) (*entity.Users, error) {
	return u.authRepo.FindByID(ctx, id)
}

