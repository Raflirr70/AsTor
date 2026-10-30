package entity

import "time"

type Loans struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserLoanID   uint       `json:"user_loan_id"`
	UserReturnID uint       `json:"user_return_id"`
	Status       string     `json:"status"`
	LoansName    string     `json:"loans_name"`
	LoansNpm     string     `json:"loans_npm"`
	CreatedAt    time.Time  `json:"created_at"`
	Deadline     *time.Time `json:"deadline" gorm:"default:null"`
	CheckedAt    *time.Time `json:"checked_at"`
	ReturnedAt   *time.Time `json:"returned_at"`

	// ForeignKey
	UserLoan    *Users       `json:"user_loan" gorm:"foreignKey:UserLoanID;references:ID"`
	UserReturn  *Users       `json:"user_return" gorm:"foreignKey:UserReturnID;references:ID"`
	LoanDetails *LoanDetails `json:"loan_details" gorm:"foreignKey:LoanID;references:ID"`
}
