package entity

import "time"

type ComputerComponents struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetID    uint      `json:"asset_id"`
	ComputerID uint      `json:"computer_id"`
	StartedAt  time.Time `json:"started_at"`
	EndedAt    time.Time `json:"ended_at" gorm:"default:null"`

	// ForeignKey
	Asset    *Assets    `json:"asset" gorm:"foreignKey:AssetID;references:ID"`
	Computer *Computers `json:"computer" gorm:"foreignKey:ComputerID;references:ID"`
}
