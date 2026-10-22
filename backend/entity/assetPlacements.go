package entity

import "time"

type AssetPlacements struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetID    uint      `json:"asset_id"`
	LocationID uint      `json:"location_id"`
	StartedAt  time.Time `json:"started_at"`
	EndedAt    time.Time `json:"ended_at" gorm:"default:null"`

	Asset    *Assets    `json:"asset" gorm:"foreignKey:AssetID;references:ID"`
	Location *Locations `json:"location" gorm:"foreignKey:LocationID;references:ID"`
}
