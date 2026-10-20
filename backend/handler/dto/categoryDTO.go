package dto

import "astor/backend/entity"

type CreateCategory struct {
	Name             string `json:"name" binding:"required"`
	Code             string `json:"code" binding:"required"`
	SpecificationIDs []uint `json:"specification_ids" binding:"required,min=1"`
}

type UpdateCategory struct {
	Name             string `json:"name" binding:"required"`
	Code             string `json:"code" binding:"required"`
	SpecificationIDs []uint `json:"specification_ids" binding:"required,min=1"`
}

type ResponseSpecificationForCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type ResponseCategory struct {
	ID             uint                               `json:"id"`
	Name           string                             `json:"name"`
	Code           string                             `json:"code"`
	Specifications []ResponseSpecificationForCategory `json:"specifications"`
}

func ToResponseCategory(r *entity.Categories) ResponseCategory {
	var specs []ResponseSpecificationForCategory
	for _, cs := range r.CategorySpecifications {
		if cs.Specification != nil {
			specs = append(specs, ResponseSpecificationForCategory{
				ID:   cs.Specification.ID,
				Name: cs.Specification.Name,
				Unit: cs.Specification.Unit,
			})
		}
	}

	return ResponseCategory{
		ID:             r.ID,
		Name:           r.Name,
		Code:           r.Code,
		Specifications: specs,
	}
}

func ToResponseCategoryList(categories []entity.Categories) []ResponseCategory {
	var result []ResponseCategory
	for _, cat := range categories {
		result = append(result, ToResponseCategory(&cat))
	}
	return result
}

func (r *CreateCategory) ToEntity() *entity.Categories {
	return &entity.Categories{
		Name: r.Name,
		Code: r.Code,
	}
}

func (r *UpdateCategory) ToEntity() *entity.Categories {
	return &entity.Categories{
		Name: r.Name,
		Code: r.Code,
	}
}
