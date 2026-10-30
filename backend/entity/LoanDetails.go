package entity

type LoanDetails struct {
	ID       uint `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetID  uint `json:"asset_id"`
	LoanID   uint `json:"loan_id"`
	Quantity int  `json:"quantity"`

	// ForeignKey
	Asset *Assets `json:"asset" gorm:"foreignKey:AssetID;references:ID"`
	Loan  *Loans  `json:"loan" gorm:"foreignKey:LoanID;references:ID"`
}
