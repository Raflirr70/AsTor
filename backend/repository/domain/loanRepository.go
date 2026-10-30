package dr

import (
	"astor/backend/entity"
	"context"
)

type LoanRepository interface {
	Create(ctx context.Context, loan *entity.Loans) error
	CreateDetail(ctx context.Context, detail *entity.LoanDetails) error
	FetchByStatus(ctx context.Context, status string) ([]entity.Loans, error)
	FetchAll(ctx context.Context) ([]entity.Loans, error)
	FetchByID(ctx context.Context, id uint) (*entity.Loans, error)
	Update(ctx context.Context, loan *entity.Loans) error
}
