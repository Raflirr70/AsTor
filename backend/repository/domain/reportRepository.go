package dr

import (
	"astor/backend/entity"
	"context"
)

type ReportRepository interface {
	Create(ctx context.Context, report *entity.Reports) error
	Update(ctx context.Context, report *entity.Reports) error
	FetchAll(ctx context.Context) ([]entity.Reports, error)
	FetchByID(ctx context.Context, id uint) (*entity.Reports, error)
}
