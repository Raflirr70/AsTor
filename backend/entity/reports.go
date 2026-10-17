package entity

import "time"

type Reports struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ReportedByID uint      `json:"reported_by_id"`
	ApprovedByID uint      `json:"approved_by_id"`
	AssetID      uint      `json:"asset_id"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
	ApprovedAt   time.Time `json:"approved_at"`

	// ForeignKey
	ReportedBy *Users  `json:"reported_by" gorm:"foreignKey:ReportedByID;references:ID"`
	ApprovedBy *Users  `json:"approved_by" gorm:"foreignKey:ApprovedByID;references:ID"`
	Asset      *Assets `json:"asset" gorm:"foreignKey:AssetID;references:ID"`
}
