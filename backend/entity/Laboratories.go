package entity

import (
	"time"

	"gorm.io/gorm"
)

type Laboratories struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID uint           `json:"location_id"`
	Capasity   int            `json:"capasity"`
	CreatedAt  time.Time      `json:"created_at"`
	CheckedAt  time.Time      `json:"checked_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-"`

	// ForeignKey
	Location *Locations `json:"location" gorm:"foreignKey:LocationID;references:ID"`

	// ForeignKey*
	ComputerLaboratories []ComputerLaboratories `json:"computer_laboratories" gorm:"foreignKey:LaboratoryID;references:ID"`
}
