package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type specificationRepositoryImpl struct {
	db *gorm.DB
}

func NewSpecificationRepositoryImpl(db *gorm.DB) dr.SpecificationRepository {
	return &specificationRepositoryImpl{db: db}
}

func (r *specificationRepositoryImpl) Create(ctx context.Context, specification *entity.Specifications) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(specification).Error
}

func (r *specificationRepositoryImpl) Update(ctx context.Context, specification *entity.Specifications) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", specification.ID).Updates(specification).Error
}

func (r *specificationRepositoryImpl) Delete(ctx context.Context, specification *entity.Specifications) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Delete(specification).Error
}

func (r *specificationRepositoryImpl) FetchAll(ctx context.Context) ([]entity.Specifications, error) {
	var specs []entity.Specifications
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Order("name ASC").Find(&specs).Error; err != nil {
		return nil, err
	}
	return specs, nil
}

func (r *specificationRepositoryImpl) FetchByID(ctx context.Context, id uint) (*entity.Specifications, error) {
	var spec entity.Specifications
	if err := dbFromContext(ctx, r.db).WithContext(ctx).First(&spec, id).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}
