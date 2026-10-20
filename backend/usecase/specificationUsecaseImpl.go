package usecase

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	du "astor/backend/usecase/domain"
	"context"
	"fmt"
)

type specificationUsecaseImpl struct {
	specRepo        dr.SpecificationRepository
	activityLogRepo dr.ActivityLogRepository
	txRepo          dr.TransactionRepository
}

func NewSpecificationUsecaseImpl(specRepo dr.SpecificationRepository, activityLogRepo dr.ActivityLogRepository, txRepo dr.TransactionRepository) du.SpecificationUsecase {
	return &specificationUsecaseImpl{specRepo: specRepo, activityLogRepo: activityLogRepo, txRepo: txRepo}
}

func (u *specificationUsecaseImpl) Create(ctx context.Context, adminID uint, specification *entity.Specifications) error {
	if specification.Name == "" {
		return entity.ErrInvalidInput
	}

	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := u.specRepo.Create(ctx, specification); err != nil {
			return fmt.Errorf("failed to create specification: %w", err)
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "SPECIFICATION_CREATED",
			EntityType:  "Specifications",
			EntityID:    specification.ID,
			Description: "Specification created",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}

func (u *specificationUsecaseImpl) Update(ctx context.Context, adminID uint, specification *entity.Specifications) error {
	if specification.Name == "" {
		return entity.ErrInvalidInput
	}

	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.specRepo.FetchByID(ctx, specification.ID)
		if err != nil {
			return fmt.Errorf("failed to find specification: %w", err)
		}

		if err := u.specRepo.Update(ctx, specification); err != nil {
			return fmt.Errorf("failed to update specification: %w", err)
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "SPECIFICATION_UPDATED",
			EntityType:  "Specifications",
			EntityID:    exist.ID,
			Description: "Specification updated",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}

func (u *specificationUsecaseImpl) GetAll(ctx context.Context) ([]entity.Specifications, error) {
	return u.specRepo.FetchAll(ctx)
}

func (u *specificationUsecaseImpl) GetByID(ctx context.Context, id uint) (*entity.Specifications, error) {
	return u.specRepo.FetchByID(ctx, id)
}

func (u *specificationUsecaseImpl) Delete(ctx context.Context, adminID uint, specification *entity.Specifications) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.specRepo.FetchByID(ctx, specification.ID)
		if err != nil {
			return fmt.Errorf("failed to find specification: %w", err)
		}

		if err := u.specRepo.Delete(ctx, exist); err != nil {
			return fmt.Errorf("failed to delete specification: %w", err)
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "SPECIFICATION_DELETED",
			EntityType:  "Specifications",
			EntityID:    exist.ID,
			Description: "Specification deleted",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}
