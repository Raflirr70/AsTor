package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type loanRepositoryImpl struct {
	db *gorm.DB
}

func NewLoanRepositoryImpl(db *gorm.DB) dr.LoanRepository {
	return &loanRepositoryImpl{db: db}
}

func (r *loanRepositoryImpl) Create(ctx context.Context, loan *entity.Loans) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(loan).Error
}

func (r *loanRepositoryImpl) CreateDetail(ctx context.Context, detail *entity.LoanDetails) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(detail).Error
}

func (r *loanRepositoryImpl) FetchByStatus(ctx context.Context, status string) ([]entity.Loans, error) {
	var loans []entity.Loans
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("User").Preload("LoanDetails.Asset").Where("status = ?", status).Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *loanRepositoryImpl) FetchAll(ctx context.Context) ([]entity.Loans, error) {
	var loans []entity.Loans
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("UserLoan").Preload("UserReturn").Preload("LoanDetails.Asset").Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}
func (r *loanRepositoryImpl) FetchByID(ctx context.Context, id uint) (*entity.Loans, error) {
	var loan entity.Loans
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("UserLoan").Preload("UserReturn").Preload("LoanDetails.Asset").First(&loan, id).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}
func (r *loanRepositoryImpl) Update(ctx context.Context, loan *entity.Loans) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", loan.ID).Updates(loan).Error
}
