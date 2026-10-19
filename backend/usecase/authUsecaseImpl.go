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
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authUsecaseImpl struct {
	authRepo        dr.AuthRepository
	activityLogRepo dr.ActivityLogRepository
	txRepo          dr.TransactionRepository
	JWT             config.JWTConfig
	Redis           *redis.Client
}

func NewAuthUsecaseImpl(authRepo dr.AuthRepository, activityLogRepo dr.ActivityLogRepository, txRepo dr.TransactionRepository, JWT config.JWTConfig, Redis *redis.Client) du.AuthUsecase {
	return &authUsecaseImpl{authRepo: authRepo, activityLogRepo: activityLogRepo, txRepo: txRepo, JWT: JWT, Redis: Redis}
}
func (u *authUsecaseImpl) Register(ctx context.Context, adminID uint, user *entity.Users) error {
	existNpm, err := u.authRepo.FindByNPM(ctx, user.Npm)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check NPM: %w", err)
	}
	if existNpm != nil {
		return entity.ErrAlreadyExists
	}
	existUsername, err := u.authRepo.FindByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check Username: %w", err)
	}
	if existUsername != nil {
		return entity.ErrAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hash)

	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := u.authRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "USER_CREATED",
			EntityType:  "Users",
			EntityID:    user.ID,
			Description: "User account created",
		}

		return u.activityLogRepo.Create(ctx, log)

	})
}
func (u *authUsecaseImpl) Login(ctx context.Context, username, password string) (*entity.Users, string, error) {
	exist, err := u.authRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, "", entity.ErrInvalidCredentials
		}

		return nil, "", fmt.Errorf("failed to find user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(password)); err != nil {
		log := &entity.ActivityLogs{
			UserID:      exist.ID,
			Action:      "LOGIN_FAILED",
			EntityType:  "Users",
			EntityID:    exist.ID,
			Description: "User Failed login attempt",
		}
		if err := u.activityLogRepo.Create(ctx, log); err != nil {
			return nil, "", errors.New("failed Create Log")
		}

		return nil, "", entity.ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(exist.ID, exist.Username, exist.Role.Name, u.JWT.Secret, u.JWT.Expiration)
	if err != nil {
		return nil, "", errors.New("failed Generate Token JWT")
	}

	redisKey := fmt.Sprintf("sessions:%d", exist.ID)
	u.Redis.Set(ctx, redisKey, token, u.JWT.Expiration)

	log := &entity.ActivityLogs{
		UserID:      exist.ID,
		Action:      "LOGIN_SUCCESS",
		EntityType:  "Users",
		EntityID:    exist.ID,
		Description: "User successfully login",
	}
	if err := u.activityLogRepo.Create(ctx, log); err != nil {
		return nil, "", errors.New("failed Create Log")
	}

	return exist, token, nil
}
func (u *authUsecaseImpl) Logout(ctx context.Context, token string) error {
	claim, err := jwt.ParseToken(token, u.JWT.Secret)
	if err != nil {
		return errors.New("invalid token")
	}

	redisKey := fmt.Sprintf("sessions:%d", claim.UserID)

	if err := u.Redis.Del(ctx, redisKey).Err(); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	ttl := time.Until(claim.ExpiresAt.Time)
	if ttl > 0 {
		blacklistKey := fmt.Sprintf("blacklist:%s", token)

		if err := u.Redis.Set(ctx, blacklistKey, "1", ttl).Err(); err != nil {
			return fmt.Errorf("failed to blacklist token: %w", err)
		}
	}

	log := &entity.ActivityLogs{
		UserID:      claim.UserID,
		Action:      "LOGOUT_SUCCESS",
		EntityType:  "Users",
		EntityID:    claim.UserID,
		Description: "Successful logout",
	}

	if err := u.activityLogRepo.Create(ctx, log); err != nil {
		return fmt.Errorf("failed to create activity log: %w", err)
	}

	return nil
}
func (u *authUsecaseImpl) GetProfile(ctx context.Context, id uint) (*entity.Users, error) {
	return u.authRepo.FindByID(ctx, id)
}
