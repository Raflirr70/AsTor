package entity

type Specifications struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name gorm:"unique;not null""`
	Unit string `json:"unit"`

	//ForeignKey*
	CategorySpecifications []CategorySpecifications `json:"category_specifications" gorm:"foreignKey:SpecificationID;references:ID"`
}
