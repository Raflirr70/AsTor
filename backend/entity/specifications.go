package entity

type Specifications struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID uint   `json:"category_id"`
	Name       string `json:"name"`
	Unit       string `json:"unit"`

	// ForeignKey
	Category Categories `json:"category" gorm:"foreignKey:CategoryID;references:ID"`

	//ForeignKey*
	AssetSpecifications []AssetSpecifications `json:"asset_specifications" gorm:"foreignKey:SpecificationID;references:ID"`
}
