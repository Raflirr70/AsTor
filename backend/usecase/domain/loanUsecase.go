package du

import (
	"astor/backend/entity"
	"context"
)

type LoanUsecase interface {
	Create(ctx context.Context, userID uint, loan *entity.Loans, quantities []int, codes []string) error
	GetByStatus(ctx context.Context, status string) ([]entity.Loans, error)
	GetAll(ctx context.Context) ([]entity.Loans, error)
	Approved(ctx context.Context, adminID uint, loan *entity.Loans) error
	Returned(ctx context.Context, userID uint, loan *entity.Loans, status string) error
}
