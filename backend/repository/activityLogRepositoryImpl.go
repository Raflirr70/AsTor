package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type activityLogRepositoryImpl struct {
	db *gorm.DB
}

func NewActivityLogRepositoryImpl(db *gorm.DB) dr.ActivityLogRepository {
	return &activityLogRepositoryImpl{db: db}
}
func (r *activityLogRepositoryImpl) Create(ctx context.Context, al *entity.ActivityLogs) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(al).Error
}
func (r *activityLogRepositoryImpl) FetchAll(ctx context.Context) ([]entity.ActivityLogs, error) {
	var als []entity.ActivityLogs
	err := dbFromContext(ctx, r.db).WithContext(ctx).Find(als).Error
	if err != nil {
		return nil, err
	}
	return als, nil
}
