package dto

import (
	"astor/backend/entity"
	"time"
)

type AssetSpecificationInput struct {
	SpecificationCategoryID uint   `json:"specification_category_id" binding:"required"`
	Value                   string `json:"value" binding:"required"`
}
type UpdateAssetSpecification struct {
	Value string `json:"value" binding:"required"`
}

type CreateAssetRequest struct {
	CategoryID          uint                      `json:"category_id" binding:"required"`
	Label               string                    `json:"label" binding:"required"`
	Status              string                    `json:"status" binding:"required"`
	AssetSpecifications []AssetSpecificationInput `json:"asset_specifications"`
}

type UpdateAssetRequest struct {
	Label  string `json:"label" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type ResponseAsset struct {
	ID                  uint                               `json:"id"`
	CategoryID          uint                               `json:"category_id"`
	Code                string                             `json:"code"`
	Label               string                             `json:"label"`
	Status              string                             `json:"status"`
	CheckedAt           time.Time                          `json:"checked_at"`
	UpdatedAt           time.Time                          `json:"updated_at"`
	AssetSpecifications []ResponseAssetSpecificationDetail `json:"asset_specifications,omitempty"`
}

type ResponseAssetSpecificationDetail struct {
	ID                      uint   `json:"id"`
	SpecificationCategoryID uint   `json:"specification_category_id"`
	Value                   string `json:"value"`
	Name                    string `json:"name,omitempty"`
	Unit                    string `json:"unit,omitempty"`
}

func ToResponseAsset(a *entity.Assets) ResponseAsset {
	var specs []ResponseAssetSpecificationDetail
	for _, as := range a.AssetSpecifications {
		name, unit := "", ""
		if as.SpecificatioCategory != nil {
			name = as.SpecificatioCategory.Name
			unit = as.SpecificatioCategory.Unit
		}
		specs = append(specs, ResponseAssetSpecificationDetail{
			ID:                      as.ID,
			SpecificationCategoryID: as.SpecificationCategoryID,
			Value:                   as.Value,
			Name:                    name,
			Unit:                    unit,
		})
	}
	return ResponseAsset{
		ID:                  a.ID,
		CategoryID:          a.CategoryID,
		Code:                a.Code,
		Label:               a.Label,
		Status:              a.Status,
		CheckedAt:           a.CheckedAt,
		UpdatedAt:           a.UpdatedAt,
		AssetSpecifications: specs,
	}
}

func ToResponseAssetList(assets []entity.Assets) []ResponseAsset {
	var list []ResponseAsset
	for _, a := range assets {
		list = append(list, ToResponseAsset(&a))
	}
	return list
}

func (r *CreateAssetRequest) ToEntity() (*entity.Assets, []entity.AssetSpecifications) {
	asset := &entity.Assets{
		CategoryID: r.CategoryID,
		Label:      r.Label,
		Status:     r.Status,
	}
	var specs []entity.AssetSpecifications
	for _, s := range r.AssetSpecifications {
		specs = append(specs, entity.AssetSpecifications{
			SpecificationCategoryID: s.SpecificationCategoryID,
			Value:                   s.Value,
		})
	}
	return asset, specs
}

func (r *UpdateAssetRequest) ToEntity(id uint) *entity.Assets {
	asset := &entity.Assets{
		ID:     id,
		Label:  r.Label,
		Status: r.Status,
	}
	return asset
}
func (r *UpdateAssetSpecification) ToEntity() *entity.AssetSpecifications {
	return &entity.AssetSpecifications{
		Value: r.Value,
	}
}
