package dto

import (
	"astor/backend/entity"
	"time"
)

type CreateReport struct {
	AssetCode string `json:"asset_code" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

type UpdateReport struct {
	Status  string `json:"status" binding:"required,oneof=Baik Buruk"`
	Message string `json:"message"`
}

type ResponseReport struct {
	ID           uint           `json:"id"`
	ReportedByID uint           `json:"reported_by_id"`
	ApprovedByID uint           `json:"approved_by_id"`
	AssetID      uint           `json:"asset_id"`
	Status       string         `json:"status"`
	Message      string         `json:"message"`
	CreatedAt    time.Time      `json:"created_at"`
	ApprovedAt   time.Time      `json:"approved_at"`
	Asset        *ResponseAsset `json:"asset,omitempty"`
}

func (r *CreateReport) ToEntity() *entity.Reports {
	return &entity.Reports{
		Message: r.Message,
		Status:  "PENDING",
	}
}

func (r *UpdateReport) ToEntity(id uint) *entity.Reports {
	return &entity.Reports{
		ID:      id,
		Message: r.Message,
	}
}

func ToResponseReport(r *entity.Reports) ResponseReport {
	var asset *ResponseAsset
	if r.Asset != nil {
		a := ToResponseAsset(r.Asset)
		asset = &a
	}
	return ResponseReport{
		ID:           r.ID,
		ReportedByID: r.ReportedByID,
		ApprovedByID: *r.ApprovedByID,
		AssetID:      r.AssetID,
		Status:       r.Status,
		Message:      r.Message,
		CreatedAt:    r.CreatedAt,
		ApprovedAt:   *r.ApprovedAt,
		Asset:        asset,
	}
}

func ToResponseReportList(reports []entity.Reports) []ResponseReport {
	var list []ResponseReport
	for i := range reports {
		list = append(list, ToResponseReport(&reports[i]))
	}
	return list
}
