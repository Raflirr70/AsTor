package dto

import (
	"astor/backend/entity"
	"time"
)

type CreateLaboratory struct {
	Name     string `json:"name" binding:"required"`
	Capasity int    `json:"capasity" binding:"min=0"`
}
type UpdateLaboratory struct {
	Name     string `json:"name" binding:"required"`
	Capasity int    `json:"capasity" binding:"min=0"`
}
type ResponseLaboratoriesForAdmin struct {
	ID         uint      `json:"id"`          //diambil dari laboratory.ID
	LocationID uint      `json:"location_id"` //diambil dari laboratory.locations.id
	Name       string    `json:"name"`        //diambil dari laboratories.locations.Name
	Capasity   int       `json:"capasity"`    //diambil dari laboratories.Cacpasity
	CheckedAt  time.Time `json:"checked_at"`  //diambil dari laboratory.CheckedAt
	UpdatedAt  time.Time `json:"updated_at"`  //diambil dari laboratory.UpdatedAt
}

func ToResponseLaboratoryForAdmin(r *entity.Laboratories) ResponseLaboratoriesForAdmin {
	locID := r.LocationID
	name := ""
	if r.Location != nil {
		locID = r.Location.ID
		name = r.Location.Name
	}
	return ResponseLaboratoriesForAdmin{
		ID:         r.ID,
		LocationID: locID,
		Name:       name,
		Capasity:   r.Capasity,
		CheckedAt:  r.CheckedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func ToResponseLaboratoriesForAdminList(labs []entity.Laboratories) []ResponseLaboratoriesForAdmin {
	var result []ResponseLaboratoriesForAdmin
	for _, l := range labs {
		result = append(result, ToResponseLaboratoryForAdmin(&l))
	}
	return result
}

func (r *CreateLaboratory) ToEntity() *entity.Laboratories {
	return &entity.Laboratories{
		Capasity: r.Capasity,
		Location: &entity.Locations{
			Name: r.Name,
		},
	}
}
func (r *UpdateLaboratory) ToEntity() *entity.Laboratories {
	return &entity.Laboratories{
		Capasity: r.Capasity,
		Location: &entity.Locations{
			Name: r.Name,
		},
	}
}
