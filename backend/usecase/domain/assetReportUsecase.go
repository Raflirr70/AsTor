package du

import (
	"astor/backend/entity"
	"context"
)

type AssetReportUsecase interface {
	CreateReport(ctx context.Context, userID uint, report *entity.Reports, code string) error
	UpdateReport(ctx context.Context, adminID uint, report *entity.Reports, status string) error
	GetAll(ctx context.Context) ([]entity.Reports, error)
	GetByID(ctx context.Context, id uint) (*entity.Reports, error)
}
