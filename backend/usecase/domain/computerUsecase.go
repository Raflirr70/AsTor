package du

import (
	"astor/backend/entity"
	"context"
)

const CategoryComputerID = 36
const CategoryComputerCode = "COM"

type UpdateComputerInput struct {
	Label        string
	Status       string
	ComponentIDs []uint
}

type ComputerUsecase interface {
	Create(ctx context.Context, req UpdateComputerInput) error
	Update(ctx context.Context, computerID uint, req UpdateComputerInput) error
	GetAll(ctx context.Context) ([]entity.Computers, error)
	GetByID(ctx context.Context, id uint) (*entity.Computers, error)
	Delete(ctx context.Context, adminID uint, computer *entity.Computers) error
}
