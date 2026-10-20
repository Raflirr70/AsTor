package dto

import "astor/backend/entity"

type CreateLaboratory struct {
	Name     string `json:"name" binding:"required"`
	Capasity int    `json:"capasity" binding:"min=0"`
}

func (r *CreateLaboratory) ToEntity() *entity.Laboratories {
	return &entity.Laboratories{
		Capasity: r.Capasity,
		Location: &entity.Locations{
			Name: r.Name,
		},
	}
}

type UpdateLaboratory struct {
	Name     string `json:"name" binding:"required"`
	Capasity int    `json:"capasity" binding:"min=0"`
}

func (r *UpdateLaboratory) ToEntity() *entity.Laboratories {
	return &entity.Laboratories{
		Capasity: r.Capasity,
		Location: &entity.Locations{
			Name: r.Name,
		},
	}
}
