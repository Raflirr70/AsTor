package entity

type Computers struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetID uint   `json:"asset_id"`
	Label   string `json:"label"`

	// ForeignKey
	Asset *Assets `json:"asset" gorm:"foreignKey:AssetID;references:ID"`

	// ForeignKey*
	ComputerLaboratories []ComputerLaboratories `json:"computer_laboratories" gorm:"foreignKey:ComputerID;references:ID"`
	ComputerComponents   []ComputerComponents   `json:"computer_components" gorm:"foreignKey:ComputerID;references:ID"`
}
