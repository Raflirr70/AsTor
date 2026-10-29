package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type assetRepositoryImpl struct {
	db *gorm.DB
}

func NewAssetRepositoryImpl(db *gorm.DB) dr.AssetRepository {
	return &assetRepositoryImpl{db: db}
}

func (r *assetRepositoryImpl) Create(ctx context.Context, asset *entity.Assets) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(asset).Error
}

func (r *assetRepositoryImpl) Update(ctx context.Context, asset *entity.Assets) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", asset.ID).Updates(asset).Error
}

func (r *assetRepositoryImpl) Delete(ctx context.Context, asset *entity.Assets) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Delete(asset).Error
}

func (r *assetRepositoryImpl) FetchAll(ctx context.Context) ([]entity.Assets, error) {
	var assets []entity.Assets
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("Category").Preload("AssetSpecifications.SpecificatioCategory").Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *assetRepositoryImpl) FetchByID(ctx context.Context, id uint) (*entity.Assets, error) {
	var asset entity.Assets
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("Category").Preload("AssetSpecifications.SpecificatioCategory").First(&asset, id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}
func (r *assetRepositoryImpl) FetchByCategory(ctx context.Context, categoryID uint) ([]entity.Assets, error) {
	var assets []entity.Assets
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Where("category_id = ?", categoryID).Preload("Category").Preload("AssetSpecifications.SpecificatioCategory").Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *assetRepositoryImpl) FindLastCode(ctx context.Context, prefix string) (string, error) {
	var asset entity.Assets
	err := dbFromContext(ctx, r.db).WithContext(ctx).Where("code LIKE ?", prefix+"%").Order("code DESC").First(&asset).Error
	if err != nil {
		return "", err
	}
	return asset.Code, nil
}

func (r *assetRepositoryImpl) UpdateCheckedAt(ctx context.Context, asset *entity.Assets) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Model(&entity.Assets{}).Where("id = ?", asset.ID).UpdateColumn("checked_at", asset.CheckedAt).Error
}
