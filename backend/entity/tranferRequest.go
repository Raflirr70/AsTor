package entity

import "time"

type TranferRequests struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetID        uint      `json:"asset_id"`
	FromLocationID uint      `json:"from_location_id"`
	ToLocationID   uint      `json:"to_location_id"`
	TranferedByID  uint      `json:"tranfered_by_id"`
	ApprovedByID   uint      `json:"approved_by_id"`
	Status         string    `json:"status"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
	ApprovedAt     time.Time `json:"approved_at"`

	// ForeignKey
	Asset        *Assets    `json:"asset" gorm:"foreignKey:AssetID;references:ID"`
	FromLocation *Locations `json:"FromLocation" gorm:"foreignKey:FromLocationID;references:ID"`
	ToLocation   *Locations `json:"ToLocation" gorm:"foreignKey:ToLocationID;references:ID"`
	TranferedBy  *Users     `json:"TranferedBy" gorm:"foreignKey:TranferedByID;references:ID"`
	ApprovedBy   *Users     `json:"ApprovedBy" gorm:"foreignKey:ApprovedByID;references:ID"`
}
