package udto

import "astor/backend/entity"

type CreateActivityLogDTO struct {
	UserID      uint
	Action      string
	EntityType  string
	EntityID    uint
	Description string
}

func (r *CreateActivityLogDTO) ToEntity() *entity.ActivityLogs {
	return &entity.ActivityLogs{
		UserID:      r.UserID,
		Action:      r.Action,
		EntityType:  r.EntityType,
		EntityID:    r.EntityID,
		Description: r.Description,
	}
}
