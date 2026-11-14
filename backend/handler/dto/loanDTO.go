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
type ResponseLoanItem struct {
	ID        uint        `json:"id"`
	AssetID   uint        `json:"asset_id"`
	AssetCode string      `json:"asset_code"`
	AssetName string      `json:"asset_name"`
	Quantity  int         `json:"quantity"`
	Asset     interface{} `json:"asset,omitempty"`
}

type ResponseLoan struct {
	ID             uint               `json:"id"`
	UserLoanID     uint               `json:"user_id"`
	UserLoanName   string             `json:"user_loan_name"`
	UserReturn     *uint              `json:"user_return"`
	UserReturnName *string            `json:"user_return_name"`
	Status         string             `json:"status"`
	LoansName      string             `json:"loans_name"`
	LoansNpm       string             `json:"loans_npm"`
	CreatedAt      time.Time          `json:"created_at"`
	Deadline       time.Time          `json:"deadline"`
	CheckedAt      *time.Time         `json:"checked_at"`
	ReturnedAt     *time.Time         `json:"returned_at,omitempty"`
	LoanDetails    []ResponseLoanItem `json:"loan_details"`
}

func ToResponseLoan(l *entity.Loans) ResponseLoan {
	var loanName string
	if l.UserLoan != nil {
		names := []string{}
		if l.UserLoan.Firstname != "" {
			names = append(names, l.UserLoan.Firstname)
		}
		if len(names) > 0 {
			loanName = names[0]
			if len(names) > 1 {
				loanName += " " + names[1]
			}
		} else {
			loanName = l.UserLoan.Username
		}
	}

	var returnName *string
	if l.UserReturn != nil {
		names := []string{}
		if l.UserReturn.Firstname != "" {
			names = append(names, l.UserReturn.Firstname)
		}
		var nameStr string
		if len(names) > 0 {
			nameStr = names[0]
			if len(names) > 1 {
				nameStr += " " + names[1]
			}
		} else {
			nameStr = l.UserReturn.Username
		}
		returnName = &nameStr
	}

	var details []ResponseLoanItem
	for _, d := range l.LoanDetails {
		code := ""
		name := ""
		if d.Asset != nil {
			code = d.Asset.Code
			name = d.Asset.Label
		}
		details = append(details, ResponseLoanItem{
			ID:        d.ID,
			AssetID:   d.AssetID,
			AssetCode: code,
			AssetName: name,
			Quantity:  d.Quantity,
			Asset:     d.Asset,
		})
	}

	return ResponseLoan{
		ID:             l.ID,
		UserLoanID:     l.UserLoanID,
		UserLoanName:   loanName,
		UserReturn:     l.UserReturnID,
		UserReturnName: returnName,
		Status:         l.Status,
		LoansName:      l.LoansName,
		LoansNpm:       l.LoansNpm,
		CreatedAt:      l.CreatedAt,
		Deadline:       *l.Deadline,
		CheckedAt:      l.CheckedAt,
		ReturnedAt:     l.ReturnedAt,
		LoanDetails:    details,
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
