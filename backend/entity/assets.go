package entity

import (
	"time"

	"gorm.io/gorm"
)

type Assets struct {
	ID         uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID uint            `json:"category_id"`
	Code       string          `json:"code"`
	Label      string          `json:"label"`
	Status     string          `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
	CheckedAt  time.Time       `json:"checked_at" gorm:"default:null"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  *gorm.DeletedAt `json:"-"`

	// ForeignKey
	Category Categories `json:"category" gorm:"foreignKey:CategoryID;references:ID"`

	// ForeignKey*
	LoanDetails         []LoanDetails         `json:"loan_details" gorm:"foreignKey:AssetID;references:ID"`
	Reports             []Reports             `json:"reports" gorm:"foreignKey:AssetID;references:ID"`
	TranferRequests     []TranferRequests     `json:"tranfer_requests" gorm:"foreignKey:AssetID;references:ID"`
	AssetPlacements     []AssetPlacements     `json:"asset_placements" gorm:"foreignKey:AssetID;references:ID"`
	Computers           []Computers           `json:"computers" gorm:"foreignKey:AssetID;references:ID"`
	ComputerComponents  []ComputerComponents  `json:"computer_components" gorm:"foreignKey:AssetID;references:ID"`
	AssetSpecifications []AssetSpecifications `json:"asset_specifications" gorm:"foreignKey:AssetID;references:ID"`
}
