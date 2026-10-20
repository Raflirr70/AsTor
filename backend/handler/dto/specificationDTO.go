package dto

import "astor/backend/entity"

type CreateSpecification struct {
	Name string `json:"name" binding:"required"`
	Unit string `json:"unit" binding:"required"`
}
type UpdateSpecification struct {
	Name string `json:"name" binding:"required"`
	Unit string `json:"unit" binding:"required"`
}
type ResponseSpecificationList struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

func ToResponseSpecification(r *entity.Specifications) ResponseSpecificationList {
	return ResponseSpecificationList{
		ID:   r.ID,
		Name: r.Name,
		Unit: r.Unit,
	}
}
func ToResponseSpecificationList(specs []entity.Specifications) []ResponseSpecificationList {
	var result []ResponseSpecificationList
	for _, s := range specs {
		result = append(result, ToResponseSpecification(&s))
	}
	return result
}
func (r *CreateSpecification) ToEntity() *entity.Specifications {
	return &entity.Specifications{
		Name: r.Name,
		Unit: r.Unit,
	}
}
func (r *UpdateSpecification) ToEntity() *entity.Specifications {
	return &entity.Specifications{
		Name: r.Name,
		Unit: r.Unit,
	}
}
