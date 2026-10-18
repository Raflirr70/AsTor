package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type authRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepositoryImpl(db *gorm.DB) dr.AuthRepository {
	return &authRepositoryImpl{db: db}
}

func (r *authRepositoryImpl) Create(ctx context.Context, user *entity.Users) error {
	return r.db.WithContext(ctx).Create(user).Error
}
func (r *authRepositoryImpl) FindByNPM(ctx context.Context, npm string) (*entity.Users, error) {
	var user entity.Users
	if err := r.db.WithContext(ctx).Preload("Role").Where("npm = ?", npm).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *authRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.Users, error) {
	var user entity.Users
	if err := r.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *authRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Users, error) {
	var user entity.Users
	if err := r.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *authRepositoryImpl) Update(ctx context.Context, user *entity.Users) error {
	return r.db.WithContext(ctx).Where("id = ?", user.ID).Updates(user).Error
}
func (r *authRepositoryImpl) FindAll(ctx context.Context) ([]entity.Users, error) {
	var users []entity.Users
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (r *authRepositoryImpl) Delete(ctx context.Context, user *entity.Users) error {
	return r.db.WithContext(ctx).Delete(user).Error
}
