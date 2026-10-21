package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type assetSpecificationRepositoryImpl struct {
	db *gorm.DB
}

func NewAssetSpecificationRepositoryImpl(db *gorm.DB) dr.AssetSpecificationRepository {
	return &assetSpecificationRepositoryImpl{db: db}
}

func (r *assetSpecificationRepositoryImpl) Create(ctx context.Context, as *entity.AssetSpecifications) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(as).Error
}

func (r *assetSpecificationRepositoryImpl) Update(ctx context.Context, as *entity.AssetSpecifications) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", as.ID).Updates(as).Error
}

func (r *assetSpecificationRepositoryImpl) Delete(ctx context.Context, as *entity.AssetSpecifications) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Delete(as).Error
}

func (r *assetSpecificationRepositoryImpl) FetchByID(ctx context.Context, id uint) (*entity.AssetSpecifications, error) {
	var as entity.AssetSpecifications
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("SpecificatioCategory").First(&as, id).Error; err != nil {
		return nil, err
	}
	return &as, nil
}

func (r *assetSpecificationRepositoryImpl) FetchByAssetID(ctx context.Context, assetID uint) ([]entity.AssetSpecifications, error) {
	var list []entity.AssetSpecifications
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Where("asset_id = ?", assetID).Preload("SpecificatioCategory").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
