package entity

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID    uint           `json:"role_id" gorm:"not null"`
	Firstname string         `json:"firstname" gorm:"not null"`
	Lastname  string         `json:"lastname"`
	Npm       string         `json:"npm" gorm:"unique"`
	Address   string         `json:"address"`
	Birthday  time.Time      `json:"birthday"`
	Username  string         `json:"username" gorm:"not null"`
	Password  string         `json:"password" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`

	// ForeignKey
	Role *Roles `json:"role" gorm:"foreignKey:RoleID;references:ID"`

	// ForeignKey*
	ActivityLogs            []ActivityLogs    `json:"activity_logs" gorm:"foreignKey:UserID;references:ID"`
	TranferRequests         []TranferRequests `json:"tranfer_requests" gorm:"foreignKey:TranferedByID;references:ID"`
	ApprovedTranferRequests []TranferRequests `json:"approved_tranfer_requests" gorm:"foreignKey:ApprovedByID;references:ID"`
	Reporteds               []Reports         `json:"reporteds" gorm:"foreignKey:ReportedByID;references:ID"`
	ApprovedReporteds       []Reports         `json:"approved_reporteds" gorm:"foreignKey:ApprovedByID;references:ID"`
	Notifications           []Notifications   `json:"notifications" gorm:"foreignKey:UserID;references:ID"`
	Loans                   []Loans           `json:"loans" gorm:"foreignKey:UserLoanID;references:ID"`
	ReturnedLoans           []Loans           `json:"returned_loans" gorm:"foreignKey:UserReturnID;references:ID"`
}
