package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepositoryImpl(db *gorm.DB) dr.CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

func (r *categoryRepositoryImpl) Create(ctx context.Context, category *entity.Categories) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(category).Error
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, category *entity.Categories) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", category.ID).Updates(category).Error
}

func (r *categoryRepositoryImpl) Delete(ctx context.Context, category *entity.Categories) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Delete(category).Error
}

func (r *categoryRepositoryImpl) FetchAll(ctx context.Context) ([]entity.Categories, error) {
	var categories []entity.Categories
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("CategorySpecifications.Specification").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepositoryImpl) FetchByID(ctx context.Context, id uint) (*entity.Categories, error) {
	var category entity.Categories
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("CategorySpecifications.Specification").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
