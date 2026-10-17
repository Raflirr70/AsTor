package entity

import "time"

type Loans struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint      `json:"user_id"`
	Status     string    `json:"status"`
	LoansName  string    `json:"loans_name"`
	LoansNpm   string    `json:"loans_npm"`
	CreatedAt  time.Time `json:"created_at"`
	Deadline   time.Time `json:"deadline"`
	CheckedAt  time.Time `json:"checked_at"`
	ReturnedAt time.Time `json:"returned_at"`

	// ForeignKey
	User *Users `json:"approved_by" gorm:"foreignKey:UserID;references:ID"`
}
