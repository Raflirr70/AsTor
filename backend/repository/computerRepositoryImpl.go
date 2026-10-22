package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type computerRepositoryImpl struct {
	db *gorm.DB
}

func NewComputerRepositoryImpl(db *gorm.DB) dr.ComputerRepository {
	return &computerRepositoryImpl{db: db}
}

func (r *computerRepositoryImpl) Create(ctx context.Context, computer *entity.Computers) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(computer).Error
}

func (r *computerRepositoryImpl) Update(ctx context.Context, computer *entity.Computers) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", computer.ID).Updates(computer).Error
}

func (r *computerRepositoryImpl) Delete(ctx context.Context, computer *entity.Computers) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Delete(computer).Error
}

func (r *computerRepositoryImpl) FetchAll(ctx context.Context) ([]entity.Computers, error) {
	var computers []entity.Computers
	if err := dbFromContext(ctx, r.db).WithContext(ctx).
		Preload("Asset.Category").
		Preload("Asset.AssetSpecifications.SpecificatioCategory").
		Preload("ComputerLaboratories.Laboratory").
		Preload("ComputerComponents", "ended_at IS NULL").
		Preload("ComputerComponents.Asset.Category").
		Preload("ComputerComponents.Asset.AssetSpecifications.SpecificatioCategory").
		Find(&computers).Error; err != nil {
		return nil, err
	}
	return computers, nil
}

func (r *computerRepositoryImpl) FetchByID(ctx context.Context, id uint) (*entity.Computers, error) {
	var computer entity.Computers
	if err := dbFromContext(ctx, r.db).WithContext(ctx).
		Preload("Asset.Category").
		Preload("Asset.AssetSpecifications.SpecificatioCategory").
		Preload("ComputerLaboratories.Laboratory").
		Preload("ComputerComponents.Asset.Category").
		Preload("ComputerComponents.Asset.AssetSpecifications.SpecificatioCategory").
		First(&computer, id).Error; err != nil {
		return nil, err
	}
	return &computer, nil
}

func (r *computerRepositoryImpl) FetchByLocation(ctx context.Context, locationID uint) ([]entity.Computers, error) {
	var computers []entity.Computers
	if err := dbFromContext(ctx, r.db).WithContext(ctx).
		Joins("JOIN computer_laboratories ON computer_laboratories.computer_id = computers.id").
		Where("computer_laboratories.location_id = ?", locationID).
		Preload("Asset.Category").
		Preload("Asset.AssetSpecifications.SpecificatioCategory").
		Preload("ComputerLaboratories.Laboratory").
		Preload("ComputerComponents.Asset.Category").
		Preload("ComputerComponents.Asset.AssetSpecifications.SpecificatioCategory").
		Find(&computers).Error; err != nil {
		return nil, err
	}
	return computers, nil
}
func (r *computerRepositoryImpl) FetchAssetID(ctx context.Context, ID uint) (uint, error) {
	var computer entity.Computers
	err := dbFromContext(ctx, r.db).WithContext(ctx).Select("asset_id").Where("id = ?", ID).First(&computer).Error
	if err != nil {
		return 0, err
	}

	return computer.AssetID, nil
}

type computerComponentRepositoryImpl struct {
	db *gorm.DB
}

func NewComputerComponentRepositoryImpl(db *gorm.DB) dr.ComputerComponentRepository {
	return &computerComponentRepositoryImpl{db: db}
}

func (r *computerComponentRepositoryImpl) Create(ctx context.Context, component *entity.ComputerComponents) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(component).Error
}

func (r *computerComponentRepositoryImpl) Update(ctx context.Context, component *entity.ComputerComponents) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", component.ID).Updates(component).Error
}

func (r *computerComponentRepositoryImpl) Delete(ctx context.Context, component *entity.ComputerComponents) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Delete(component).Error
}

func (r *computerComponentRepositoryImpl) FetchByComputerID(ctx context.Context, computerID uint) ([]entity.ComputerComponents, error) {
	var components []entity.ComputerComponents
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Where("computer_id = ?", computerID).Preload("Asset").Find(&components).Error; err != nil {
		return nil, err
	}
	return components, nil
}
func (r *computerComponentRepositoryImpl) CheckUsed(ctx context.Context, assetID uint) (bool, error) {
	var count int64

	err := dbFromContext(ctx, r.db).WithContext(ctx).Model(&entity.ComputerComponents{}).Where("asset_id = ?", assetID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
