package repository

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	"context"

	"gorm.io/gorm"
)

type reportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepositoryImpl(db *gorm.DB) dr.ReportRepository {
	return &reportRepositoryImpl{db: db}
}

func (r *reportRepositoryImpl) Create(ctx context.Context, report *entity.Reports) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Create(report).Error
}

func (r *reportRepositoryImpl) Update(ctx context.Context, report *entity.Reports) error {
	return dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", report.ID).Updates(report).Error
}

func (r *reportRepositoryImpl) FetchAll(ctx context.Context) ([]entity.Reports, error) {
	var reports []entity.Reports
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("Asset").Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *reportRepositoryImpl) FetchByID(ctx context.Context, id uint) (*entity.Reports, error) {
	var report entity.Reports
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Preload("Asset").First(&report, id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}
