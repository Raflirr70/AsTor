package entity

import "time"

type ComputerLaboratories struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ComputerID   uint      `json:"computer_id"`
	LaboratoryID uint      `json:"laboratory_id"`
	Number       int       `json:"number"`
	StartedAt    time.Time `json:"started_at"`
	EndedAt      time.Time `json:"ended_at" gorm:"default:null"`

	// ForeignKey
	Computer   *Computers    `json:"computer" gorm:"foreignKey:ComputerID;references:ID"`
	Laboratory *Laboratories `json:"laboratory" gorm:"foreignKey:LaboratoryID;references:ID"`
}
