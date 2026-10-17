package entity

type Roles struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"not null"`

	// ForeignKey
	Users []Users `json:"users" gorm:"foreignKey:RoleID;references:ID"`
}
