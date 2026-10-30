package dto

import (
	"astor/backend/entity"
	"time"
)

type CreateLoanItem struct {
	AssetCode string `json:"asset_code" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type CreateLoan struct {
	LoansName string           `json:"loans_name" binding:"required"`
	LoansNpm  string           `json:"loans_npm" binding:"required"`
	Deadline  time.Time        `json:"deadline"`
	Items     []CreateLoanItem `json:"items" binding:"required,min=1,dive"`
}
type ResponseLoan struct {
	ID         uint       `json:"id"`
	UserLoanID uint       `json:"user_id"`
	UserReturn uint       `json:"user_return"`
	Status     string     `json:"status"`
	LoansName  string     `json:"loans_name"`
	LoansNpm   string     `json:"loans_npm"`
	CreatedAt  time.Time  `json:"created_at"`
	Deadline   time.Time  `json:"deadline"`
	CheckedAt  *time.Time `json:"checked_at"`
	ReturnedAt *time.Time `json:"returned_at,omitempty"`
}

func ToResponseLoan(l *entity.Loans) ResponseLoan {
	return ResponseLoan{
		ID:         l.ID,
		UserLoanID: l.UserLoanID,
		UserReturn: l.UserReturnID,
		Status:     l.Status,
		LoansName:  l.LoansName,
		LoansNpm:   l.LoansNpm,
		CreatedAt:  l.CreatedAt,
		Deadline:   *l.Deadline,
		CheckedAt:  l.CheckedAt,
		ReturnedAt: l.ReturnedAt,
	}
}
func (req *CreateLoan) ToEntity() *entity.Loans {
	deadline := req.Deadline

	if deadline.IsZero() {
		deadline = time.Now().AddDate(0, 0, 7)
	}
	return &entity.Loans{
		LoansName: req.LoansName,
		LoansNpm:  req.LoansNpm,
		Deadline:  &deadline,
		Status:    "PENDING",
	}
}
