package entity

type CategorySpecifications struct {
	ID              uint `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID      uint `json:"category_id"`
	SpecificationID uint `json:"specification_id"`

	Category      *Categories     `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
	Specification *Specifications `json:"category" gorm:"foreignKey:SpecificationID;references:ID"`
}
