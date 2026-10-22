package dto

import (
	"astor/backend/entity"
	"time"
)

type CreateComputer struct {
	Label      string `json:"label" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Components []uint `json:"components"`
}
type UpdateComputer struct {
	Label      string `json:"label" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Components []uint `json:"components"`
}

type ResponseSpecificationComponent struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
}
type ResponseComponent struct {
	ID             uint                             `json:"id"`
	AssetID        uint                             `json:"asset_id"`
	Specifications []ResponseSpecificationComponent `json:"specifications"`
}
type ResponseComputer struct {
	ID        uint                `json:"id"`
	AssetID   uint                `json:"asset_id"`
	Category  string              `json:"category"`
	Code      string              `json:"code"`
	Label     string              `json:"label"`
	Status    string              `json:"status"`
	Component []ResponseComponent `json:"component"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	CheckedAt time.Time           `json:"checked_at"`
}

func ToResponseComputer(comp *entity.Computers) ResponseComputer {
	components := make([]ResponseComponent, 0, len(comp.ComputerComponents))

	for _, cc := range comp.ComputerComponents {
		specifications := make([]ResponseSpecificationComponent, 0)

		// Mapping spesifikasi komponen
		for _, spec := range cc.Asset.AssetSpecifications {
			specifications = append(specifications, ResponseSpecificationComponent{
				ID:    spec.ID,
				Name:  spec.SpecificatioCategory.Name,
				Value: spec.Value,
				Unit:  spec.SpecificatioCategory.Unit,
			})
		}

		components = append(components, ResponseComponent{
			ID:             cc.ID,
			AssetID:        cc.AssetID,
			Specifications: specifications,
		})
	}
	resp := ResponseComputer{
		ID:        comp.ID,
		AssetID:   comp.AssetID,
		Category:  comp.Asset.Category.Name,
		Code:      comp.Asset.Code,
		Label:     comp.Asset.Label,
		Status:    comp.Asset.Status,
		CreatedAt: comp.Asset.CreatedAt,
		UpdatedAt: comp.Asset.UpdatedAt,
		CheckedAt: comp.Asset.CreatedAt,
		Component: components,
	}
	return resp
}
