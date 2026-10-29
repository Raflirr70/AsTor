package usecase

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	du "astor/backend/usecase/domain"
	"context"
	"fmt"
)

type categoryUsecaseImpl struct {
	categoryRepo    dr.CategoryRepository
	catSpecRepo     dr.CategorySpecificationRepository
	activityLogRepo dr.ActivityLogRepository
	txRepo          dr.TransactionRepository
}

func NewCategoryUsecaseImpl(categoryRepo dr.CategoryRepository, catSpecRepo dr.CategorySpecificationRepository, activityLogRepo dr.ActivityLogRepository, txRepo dr.TransactionRepository) du.CategoryUsecase {
	return &categoryUsecaseImpl{categoryRepo: categoryRepo, catSpecRepo: catSpecRepo, activityLogRepo: activityLogRepo, txRepo: txRepo}
}

func (u *categoryUsecaseImpl) Create(ctx context.Context, adminID uint, category *entity.Categories, specificationIDs []uint) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		fmt.Print("Name Mulai")
		exists, err := u.categoryRepo.ExistByName(ctx, category.Name)
		if err != nil {
			return err
		}
		fmt.Print("Nama Lewat")
		if exists {
			return fmt.Errorf("%w:category name", entity.ErrAlreadyExists)
		}
		fmt.Print("Code Mulai")

		exists, err = u.categoryRepo.ExistByCode(ctx, category.Code)
		if err != nil {
			return err
		}
		fmt.Print("Code Lewat")
		if exists {
			return fmt.Errorf("%w:category code", entity.ErrAlreadyExists)
		}
		if err := u.categoryRepo.Create(ctx, category); err != nil {
			return fmt.Errorf("failed to create category: %w", err)
		}

		for _, specID := range specificationIDs {
			cs := &entity.CategorySpecifications{
				CategoryID:      category.ID,
				SpecificationID: specID,
			}
			if err := u.catSpecRepo.Create(ctx, cs); err != nil {
				return fmt.Errorf("failed to create category specification: %w", err)
			}
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "CATEGORY_CREATED",
			EntityType:  "Categories",
			EntityID:    category.ID,
			Description: "Category and specifications created",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}

func (u *categoryUsecaseImpl) Update(ctx context.Context, adminID uint, category *entity.Categories, specificationIDs []uint) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.categoryRepo.FetchByID(ctx, category.ID)
		if err != nil {
			return fmt.Errorf("failed to find category: %w", err)
		}

		if err := u.categoryRepo.Update(ctx, category); err != nil {
			return fmt.Errorf("failed to update category: %w", err)
		}

		existingSpecs, err := u.catSpecRepo.FetchByCategoryID(ctx, category.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch existing category specifications: %w", err)
		}

		existingMap := make(map[uint]bool)
		for _, es := range existingSpecs {
			existingMap[es.SpecificationID] = true
		}

		reqMap := make(map[uint]bool)
		for _, id := range specificationIDs {
			reqMap[id] = true
		}

		// 2. Hapus spesifikasi yang tidak ada di request
		for _, es := range existingSpecs {
			if !reqMap[es.SpecificationID] {
				if err := u.catSpecRepo.DeleteByCategoryAndSpecificationID(ctx, category.ID, es.SpecificationID); err != nil {
					return fmt.Errorf("failed to delete category specification: %w", err)
				}
			}
		}

		// 3. Buat spesifikasi baru yang belum ada di DB
		for _, id := range specificationIDs {
			if !existingMap[id] {
				cs := &entity.CategorySpecifications{
					CategoryID:      category.ID,
					SpecificationID: id,
				}
				if err := u.catSpecRepo.Create(ctx, cs); err != nil {
					return fmt.Errorf("failed to create category specification: %w", err)
				}
			}
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "CATEGORY_UPDATED",
			EntityType:  "Categories",
			EntityID:    exist.ID,
			Description: "Category and specifications updated",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}

func (u *categoryUsecaseImpl) Delete(ctx context.Context, adminID uint, category *entity.Categories) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.categoryRepo.FetchByID(ctx, category.ID)
		if err != nil {
			return fmt.Errorf("failed to find category: %w", err)
		}

		if err := u.catSpecRepo.DeleteByCategoryID(ctx, category.ID); err != nil {
			return fmt.Errorf("failed to delete category specifications: %w", err)
		}

		if err := u.categoryRepo.Delete(ctx, exist); err != nil {
			return fmt.Errorf("failed to delete category: %w", err)
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "CATEGORY_DELETED",
			EntityType:  "Categories",
			EntityID:    exist.ID,
			Description: "Category deleted",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}

func (u *categoryUsecaseImpl) GetAll(ctx context.Context) ([]entity.Categories, error) {
	return u.categoryRepo.FetchAll(ctx)
}

func (u *categoryUsecaseImpl) GetByID(ctx context.Context, id uint) (*entity.Categories, error) {
	return u.categoryRepo.FetchByID(ctx, id)
}
