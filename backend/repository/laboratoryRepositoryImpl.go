package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type laboratoryRepositoryImpl struct {
	db *gorm.DB
}

func NewLaboratoryRepositoryImpl(db *gorm.DB) dr.LaboratoryRepository {
	return &laboratoryRepositoryImpl{db: db}
}

func (r *laboratoryRepositoryImpl) Create(ctx context.Context, laboratory *entity.Laboratories) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(laboratory).Error
}

func (r *laboratoryRepositoryImpl) Update(ctx context.Context, laboratory *entity.Laboratories) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", laboratory.ID).Updates(laboratory).Error
}

func (r *laboratoryRepositoryImpl) Delete(ctx context.Context, laboratory *entity.Laboratories) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Delete(laboratory).Error
}

func (r *laboratoryRepositoryImpl) FetchAll(ctx context.Context) ([]entity.Laboratories, error) {
	var laboratories []entity.Laboratories
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Joins("Location").Order(`"Location"."name" ASC`).Find(&laboratories).Error; err != nil {
		return nil, err
	}
	return laboratories, nil
}

func (r *laboratoryRepositoryImpl) FetchByID(ctx context.Context, id uint) (*entity.Laboratories, error) {
	var laboratory entity.Laboratories
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("Location").First(&laboratory, id).Error; err != nil {
		return nil, err
	}
	return &laboratory, nil
}

func (r *laboratoryRepositoryImpl) FetchByLocation(ctx context.Context, locationID uint) (*entity.Laboratories, error) {
	var laboratory entity.Laboratories
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Where("location_id = ?", locationID).First(&laboratory).Error; err != nil {
		return nil, err
	}
	return &laboratory, nil
}

func (r *laboratoryRepositoryImpl) Restore(ctx context.Context, id uint) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Unscoped().Model(&entity.Laboratories{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *laboratoryRepositoryImpl) FetchDeletedByID(ctx context.Context, id uint) (*entity.Laboratories, error) {
	var laboratory entity.Laboratories
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Unscoped().First(&laboratory, id).Error; err != nil {
		return nil, err
	}
	return &laboratory, nil
}
