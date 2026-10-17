package entity

type AssetSpecifications struct {
	ID              uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	SpecificationID uint    `json:"specification_id"`
	AssetID         uint    `json:"asset_id"`
	Value           float32 `json:"value"`

	// ForeignKey
	Specification *Specifications `json:"specification" gorm:"foreignKey:SpecificationID;references:ID"`
	Asset         *Assets         `json:"asset" gorm:"foreignKey:AssetID;references:ID"`
}
