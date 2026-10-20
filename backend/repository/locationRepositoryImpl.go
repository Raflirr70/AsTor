package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type locationRepositoryImpl struct {
	db *gorm.DB
}

func NewLocationRepositoryImpl(db *gorm.DB) dr.LocationRepository {
	return &locationRepositoryImpl{db: db}
}

func (r *locationRepositoryImpl) Create(ctx context.Context, location *entity.Locations) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(location).Error
}

func (r *locationRepositoryImpl) Update(ctx context.Context, location *entity.Locations) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", location.ID).Updates(location).Error
}

func (r *locationRepositoryImpl) Delete(ctx context.Context, location *entity.Locations) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Delete(location).Error
}

func (r *locationRepositoryImpl) FetchAll(ctx context.Context) ([]entity.Locations, error) {
	var locations []entity.Locations
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *locationRepositoryImpl) FetchByID(ctx context.Context, id uint) (*entity.Locations, error) {
	var location entity.Locations
	if err := dbFromContext(ctx, r.db).WithContext(ctx).First(&location, id).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepositoryImpl) Restore(ctx context.Context, id uint) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Unscoped().Model(&entity.Locations{}).Where("id = ?", id).Update("deleted_at", nil).Error
}
