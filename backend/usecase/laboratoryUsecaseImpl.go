package usecase

import (
	"astor/backend/entity"

	dr "astor/backend/repository/domain"
	du "astor/backend/usecase/domain"
	"context"
	"fmt"
)

type laboratoryUsecaseImpl struct {
	laboratoryRepo  dr.LaboratoryRepository
	locationRepo    dr.LocationRepository
	activityLogRepo dr.ActivityLogRepository
	txRepo          dr.TransactionRepository
}

func NewLaboratoryUsecaseImpl(laboratoryRepo dr.LaboratoryRepository, locationRepo dr.LocationRepository, activityLogRepo dr.ActivityLogRepository, txRepo dr.TransactionRepository) du.LaboratoryUsecase {
	return &laboratoryUsecaseImpl{laboratoryRepo: laboratoryRepo, locationRepo: locationRepo, activityLogRepo: activityLogRepo, txRepo: txRepo}
}

func (u *laboratoryUsecaseImpl) Create(ctx context.Context, adminID uint, laboratory *entity.Laboratories) error {
	if laboratory.Location == nil || laboratory.Location.Name == "" {
		return entity.ErrInvalidInput
	}

	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := u.locationRepo.Create(ctx, laboratory.Location); err != nil {
			return fmt.Errorf("failed to create location: %w", err)
		}

		laboratory.LocationID = laboratory.Location.ID

		if err := u.laboratoryRepo.Create(ctx, laboratory); err != nil {
			return fmt.Errorf("failed to create laboratory: %w", err)
		}
		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "LABORATORY_CREATED",
			EntityType:  "Laboratory",
			EntityID:    laboratory.ID,
			Description: "Laboratory location created",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}

func (u *laboratoryUsecaseImpl) Update(ctx context.Context, adminID uint, laboratory *entity.Laboratories) error {
	if laboratory.Location == nil || laboratory.Location.Name == "" {
		return entity.ErrInvalidInput
	}

	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.laboratoryRepo.FetchByID(ctx, laboratory.ID)
		if err != nil {
			return fmt.Errorf("failed to find laboratory: %w", err)
		}
		if err := u.laboratoryRepo.Update(ctx, &entity.Laboratories{ID: exist.ID, Capasity: laboratory.Capasity}); err != nil {
			return fmt.Errorf("failed to update laboratory: %w", err)
		}
		if err := u.locationRepo.Update(ctx, &entity.Locations{ID: exist.LocationID, Name: laboratory.Location.Name}); err != nil {
			return fmt.Errorf("failed to update location name: %w", err)
		}
		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "LABORATORY_UPDATE",
			EntityType:  "Laboratory",
			EntityID:    laboratory.ID,
			Description: "Laboratory location updated",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}

func (u *laboratoryUsecaseImpl) Delete(ctx context.Context, adminID uint, laboratory *entity.Laboratories) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.laboratoryRepo.FetchByID(ctx, laboratory.ID)
		if err != nil {
			return fmt.Errorf("failed to find laboratory: %w", err)
		}

		if err := u.laboratoryRepo.Delete(ctx, exist); err != nil {
			return fmt.Errorf("failed to delete laboratory: %w", err)
		}
		if err := u.locationRepo.Delete(ctx, exist.Location); err != nil {
			return fmt.Errorf("failed to delete locations: %w", err)
		}
		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "LABORATORY_DELETED",
			EntityType:  "Laboratory",
			EntityID:    laboratory.ID,
			Description: "Laboratory location deleted",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}

func (u *laboratoryUsecaseImpl) Restore(ctx context.Context, adminID uint, laboratory *entity.Laboratories) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.laboratoryRepo.FetchDeletedByID(ctx, laboratory.ID)
		if err != nil {
			return fmt.Errorf("failed to find laboratory: %w", err)
		}
		if !exist.DeletedAt.Valid {
			return entity.ErrAlreadyExists
		}

		if err := u.laboratoryRepo.Restore(ctx, exist.ID); err != nil {
			return fmt.Errorf("failed to restore laboratory: %w", err)
		}
		if err := u.locationRepo.Restore(ctx, exist.LocationID); err != nil {
			return fmt.Errorf("failed to restore location: %w", err)
		}
		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "LABORATORY_RESTORED",
			EntityType:  "Laboratory",
			EntityID:    exist.ID,
			Description: "Laboratory and location restored",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}
