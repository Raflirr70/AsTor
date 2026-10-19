package entity

type AssetSpecifications struct {
	ID                      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	SpecificationCategoryID uint   `json:"specification_category_id"`
	AssetID                 uint   `json:"asset_id"`
	Value                   string `json:"value"`

	// ForeignKey
	SpecificatioCategory *Specifications `json:"specification" gorm:"foreignKey:SpecificationCategoryID;references:ID"`
	Asset                *Assets         `json:"asset" gorm:"foreignKey:AssetID;references:ID"`
}
