package entity

import "gorm.io/gorm"

type Locations struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name"`
	DeletedAt gorm.DeletedAt `json:"-"`

	Laboratory      *Laboratories     `json:"laboratory" gorm:"foreignKey:LocationID;references:ID"`
	AssetPlacements []AssetPlacements `json:"asset_placements" gorm:"foreignKey:LocationID;references:ID"`
}
