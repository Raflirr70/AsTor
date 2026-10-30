package usecase

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	du "astor/backend/usecase/domain"
	"context"
	"fmt"
	"time"
)

type loanUsecaseImpl struct {
	loanRepo        dr.LoanRepository
	assetRepo       dr.AssetRepository
	txRepo          dr.TransactionRepository
	activityLogRepo dr.ActivityLogRepository
}

func NewLoanUsecaseImpl(loanRepo dr.LoanRepository, assetRepo dr.AssetRepository, txRepo dr.TransactionRepository, activityLogRepo dr.ActivityLogRepository) du.LoanUsecase {
	return &loanUsecaseImpl{loanRepo: loanRepo, assetRepo: assetRepo, txRepo: txRepo, activityLogRepo: activityLogRepo}
}

func (u *loanUsecaseImpl) Create(ctx context.Context, userID uint, loan *entity.Loans, quantities []int, codes []string) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		loan.UserLoanID = userID
		if err := u.loanRepo.Create(ctx, loan); err != nil {
			return fmt.Errorf("failed to create loan: %w", err)
		}
		for i := range codes {
			asset, err := u.assetRepo.FetchByCode(ctx, codes[i])
			if err != nil {
				return fmt.Errorf("failed to fetch asset: %w", err)
			}
			if asset.Status != "Dipinjam" {
				asset.Status = "Dipinjam"
				if err := u.assetRepo.Update(ctx, asset); err != nil {
					return fmt.Errorf("failed to updated asset: %w", err)
				}
				loanDetail := &entity.LoanDetails{
					AssetID:  asset.ID,
					LoanID:   loan.ID,
					Quantity: quantities[i],
				}
				if err := u.loanRepo.CreateDetail(ctx, loanDetail); err != nil {
					return fmt.Errorf("failed to create loan detail: %w", err)
				}
			} else {
				return fmt.Errorf("failed to create loans")
			}
		}
		log := &entity.ActivityLogs{
			UserID:      userID,
			Action:      "LOAN_CREATED",
			EntityType:  "Loan",
			EntityID:    loan.ID,
			Description: "Loan created",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}
func (u *loanUsecaseImpl) Approved(ctx context.Context, adminID uint, loan *entity.Loans) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		loan, err := u.loanRepo.FetchByID(ctx, loan.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch loans")
		}
		loan.Status = "Approve"
		check := time.Now()
		loan.CheckedAt = &check
		if err := u.loanRepo.Update(ctx, loan); err != nil {
			return fmt.Errorf("failed to approved loans")
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "LOAN_APPROVED",
			EntityType:  "Loan",
			EntityID:    loan.ID,
			Description: "Loan approved",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}
func (u *loanUsecaseImpl) Returned(ctx context.Context, userID uint, loan *entity.Loans, status string) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		loan, err := u.loanRepo.FetchByID(ctx, loan.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch loan: %w", err)
		}
		loan.Status = "Done"
		loan.UserReturnID = userID
		returned := time.Now()
		loan.ReturnedAt = &returned
		if err := u.loanRepo.Update(ctx, loan); err != nil {
			return fmt.Errorf("failed to update loan: %w", err)
		}
		asset, err := u.assetRepo.FetchByID(ctx, loan.LoanDetails.AssetID)
		if err != nil || asset == nil {
			return fmt.Errorf("failed to fetch loan: %w", err)
		}
		asset.Status = status
		if err := u.assetRepo.Update(ctx, asset); err != nil {
			return fmt.Errorf("failed to update status assets: %w", err)
		}

		log := &entity.ActivityLogs{
			UserID:      userID,
			Action:      "LOAN_RETURNED",
			EntityType:  "Loan",
			EntityID:    loan.ID,
			Description: "Loan returned",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}
func (u *loanUsecaseImpl) GetByStatus(ctx context.Context, status string) ([]entity.Loans, error) {
	return u.loanRepo.FetchByStatus(ctx, status)
}
func (u *loanUsecaseImpl) GetAll(ctx context.Context) ([]entity.Loans, error) {
	return u.loanRepo.FetchAll(ctx)
}
