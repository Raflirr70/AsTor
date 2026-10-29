package usecase

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	du "astor/backend/usecase/domain"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type assetUsecaseImpl struct {
	assetRepo        dr.AssetRepository
	categoryRepo     dr.CategoryRepository
	categorySpecRepo dr.CategorySpecificationRepository
	assetSpecRepo    dr.AssetSpecificationRepository
	activityLog      dr.ActivityLogRepository
	txRepo           dr.TransactionRepository
}

func NewAssetUsecaseImpl(
	assetRepo dr.AssetRepository,
	categoryRepo dr.CategoryRepository,
	categorySpecRepo dr.CategorySpecificationRepository,
	assetSpecRepo dr.AssetSpecificationRepository,
	activityLog dr.ActivityLogRepository,
	txRepo dr.TransactionRepository,
) du.AssetUsecase {
	return &assetUsecaseImpl{
		assetRepo:        assetRepo,
		categoryRepo:     categoryRepo,
		categorySpecRepo: categorySpecRepo,
		assetSpecRepo:    assetSpecRepo,
		activityLog:      activityLog,
		txRepo:           txRepo,
	}
}

func (u *assetUsecaseImpl) Create(ctx context.Context, adminID uint, asset *entity.Assets, assetSpecs []entity.AssetSpecifications) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		cat, err := u.categoryRepo.FetchByID(ctx, asset.CategoryID)
		if err != nil {
			return fmt.Errorf("failed to find category: %w", err)
		}

		categorySpecs, err := u.categorySpecRepo.FetchByCategoryID(ctx, cat.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch category specifications: %w", err)
		}

		var requiredIDs []uint
		for _, cs := range categorySpecs {
			requiredIDs = append(requiredIDs, cs.SpecificationID)
		}

		if err := validateAssetSpecifications(requiredIDs, assetSpecs); err != nil {
			return err
		}

		prefix := cat.Code + "-"
		lastCode, err := u.assetRepo.FindLastCode(ctx, prefix)
		nextNum := 1
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to fetch last asset code: %w", err)
		}
		if lastCode != "" {
			var num int
			if _, err := fmt.Sscanf(lastCode, prefix+"%d", &num); err == nil {
				nextNum = num + 1
			}
		}
		asset.Code = fmt.Sprintf("%s%03d", prefix, nextNum)

		if err := u.assetRepo.Create(ctx, asset); err != nil {
			return fmt.Errorf("failed to create asset: %w", err)
		}

		for _, spec := range assetSpecs {
			spec.AssetID = asset.ID
			if err := u.assetSpecRepo.Create(ctx, &spec); err != nil {
				return fmt.Errorf("failed to create asset specification: %w", err)
			}
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "ASSET_CREATED",
			EntityType:  "Assets",
			EntityID:    asset.ID,
			Description: "Asset created with specifications",
		}
		return u.activityLog.Create(ctx, log)
	})
}

func (u *assetUsecaseImpl) Update(ctx context.Context, adminID uint, asset *entity.Assets) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.assetRepo.FetchByID(ctx, asset.ID)
		if err != nil {
			return fmt.Errorf("failed to find asset: %w", err)
		}
		assetChanged := exist.Label != asset.Label || exist.Status != asset.Status
		if !assetChanged {
			return fmt.Errorf("failed to update asset: %w", entity.ErrNoChanges)
		}

		if exist.Status != asset.Status {
			asset.CheckedAt = time.Now()
		}

		if err := u.assetRepo.Update(ctx, asset); err != nil {
			return fmt.Errorf("failed to update asset: %w", err)
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "ASSET_UPDATED",
			EntityType:  "Assets",
			EntityID:    exist.ID,
			Description: "Asset updated",
		}
		if err := u.activityLog.Create(ctx, log); err != nil {
			return fmt.Errorf("failed to create activity log: %w", err)
		}
		return nil
	})
}
func (u *assetUsecaseImpl) Checked(ctx context.Context, adminID uint, asset *entity.Assets) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.assetRepo.FetchByID(ctx, asset.ID)
		if err != nil {
			return fmt.Errorf("failed to find asset: %w", err)
		}
		exist.CheckedAt = asset.CheckedAt
		if err := u.assetRepo.UpdateCheckedAt(ctx, exist); err != nil {
			return fmt.Errorf("failed to update checkedAt asset: %w", err)
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "ASSET_Checked",
			EntityType:  "Assets",
			EntityID:    exist.ID,
			Description: "Asset checked",
		}
		if err := u.activityLog.Create(ctx, log); err != nil {
			return fmt.Errorf("failed to create activity log: %w", err)
		}
		return nil
	})
}

func (u *assetUsecaseImpl) UpdateSpecification(ctx context.Context, adminID uint, spec *entity.AssetSpecifications) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.assetSpecRepo.FetchByID(ctx, spec.ID)
		if err != nil {
			return fmt.Errorf("failed to find asset specification : %w", err)
		}
		if spec.Value == exist.Value {
			return entity.ErrNoChanges
		}
		if err := u.assetRepo.Update(ctx, spec.Asset); err != nil {
			return fmt.Errorf("failed  to update asset specification : %w", err)
		}
		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "ASSET_SPECIFICATION_UPDATED",
			EntityType:  "Asset Specification",
			EntityID:    exist.ID,
			Description: "Asset specification updated",
		}
		if err := u.activityLog.Create(ctx, log); err != nil {
			return fmt.Errorf("failed to create activity log: %w", err)
		}
		return nil
	})
}

func (u *assetUsecaseImpl) Delete(ctx context.Context, adminID uint, asset *entity.Assets) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.assetRepo.FetchByID(ctx, asset.ID)
		if err != nil {
			return fmt.Errorf("failed to find asset: %w", err)
		}
		if err := u.assetRepo.Delete(ctx, exist); err != nil {
			return fmt.Errorf("failed to delete asset: %w", err)
		}
		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "ASSET_DELETED",
			EntityType:  "Assets",
			EntityID:    exist.ID,
			Description: "Asset deleted",
		}
		if err := u.activityLog.Create(ctx, log); err != nil {
			return fmt.Errorf("failed to create activity log: %w", err)
		}
		return nil
	})
}

func (u *assetUsecaseImpl) FetchAll(ctx context.Context) ([]entity.Assets, error) {
	return u.assetRepo.FetchAll(ctx)
}

func (u *assetUsecaseImpl) FetchByID(ctx context.Context, id uint) (*entity.Assets, error) {
	return u.assetRepo.FetchByID(ctx, id)
}

func validateAssetSpecifications(requiredIDs []uint, assetSpecs []entity.AssetSpecifications) error {
	if len(requiredIDs) == 0 {
		if len(assetSpecs) > 0 {
			return fmt.Errorf("%w: this category does not support specifications", entity.ErrInvalidInput)
		}
		return nil
	}

	provided := make(map[uint]string)

	for _, spec := range assetSpecs {
		value := strings.TrimSpace(spec.Value)

		if value == "" {
			return fmt.Errorf("%w: specification %d must not be empty", entity.ErrInvalidInput, spec.SpecificationCategoryID)
		}

		id := spec.SpecificationCategoryID

		if _, exists := provided[id]; exists {
			return fmt.Errorf("%w: duplicate specification ID: %d", entity.ErrInvalidInput, id)
		}

		provided[id] = value
	}

	return nil
}
