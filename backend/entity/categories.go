package entity

type Categories struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	Code string `json:"code"`

	// ForeignKey*
	CategorySpecifications []CategorySpecifications `json:"category_specifications" gorm:"foreignKey:CategoryID;references:ID"`
	Assets                 []Assets                 `json:"assets" gorm:"foreignKey:CategoryID;references:ID"`
}
