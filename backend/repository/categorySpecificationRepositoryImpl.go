package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type categorySpecificationRepositoryImpl struct {
	db *gorm.DB
}

func NewCategorySpecificationRepositoryImpl(db *gorm.DB) dr.CategorySpecificationRepository {
	return &categorySpecificationRepositoryImpl{db: db}
}

func (r *categorySpecificationRepositoryImpl) Create(ctx context.Context, cs *entity.CategorySpecifications) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(cs).Error
}

func (r *categorySpecificationRepositoryImpl) DeleteByCategoryID(ctx context.Context, categoryID uint) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("category_id = ?", categoryID).Delete(&entity.CategorySpecifications{}).Error
}

func (r *categorySpecificationRepositoryImpl) FetchByCategoryID(ctx context.Context, categoryID uint) ([]entity.CategorySpecifications, error) {
	var css []entity.CategorySpecifications
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Where("category_id = ?", categoryID).Find(&css).Error; err != nil {
		return nil, err
	}
	return css, nil
}
