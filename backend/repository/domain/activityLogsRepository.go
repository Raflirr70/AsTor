package dr

import (
	"astor/backend/entity"
	"context"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, al *entity.ActivityLogs) error
	FetchAll(ctx context.Context) ([]entity.ActivityLogs, error)
}
