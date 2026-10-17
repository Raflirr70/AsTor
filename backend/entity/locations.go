package entity

type Locations struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`

	Laboratory      *Laboratories     `json:"laboratory" gorm:"foreignKey:LocationID;references:ID"`
	AssetPlacements []AssetPlacements `json:"asset_placements" gorm:"foreignKey:LocationID;references:ID"`
}
