package usecase

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	du "astor/backend/usecase/domain"
	"context"
	"fmt"
	"time"
)

type assetReportUsecaseImpl struct {
	reportRepo      dr.ReportRepository
	assetRepo       dr.AssetRepository
	txRepo          dr.TransactionRepository
	activityLogRepo dr.ActivityLogRepository
}

func NewAssetReportUsecaseImpl(reportRepo dr.ReportRepository, assetRepo dr.AssetRepository, txRepo dr.TransactionRepository, activityLogRepo dr.ActivityLogRepository) du.AssetReportUsecase {
	return &assetReportUsecaseImpl{reportRepo: reportRepo, assetRepo: assetRepo, txRepo: txRepo, activityLogRepo: activityLogRepo}
}

func (u *assetReportUsecaseImpl) CreateReport(ctx context.Context, userID uint, report *entity.Reports, code string) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		asset, err := u.assetRepo.FetchByCode(ctx, code)
		if err != nil {
			return fmt.Errorf("failed to fetch asset: %w", err)
		}
		report.AssetID = asset.ID
		report.ReportedByID = userID
		if err := u.reportRepo.Create(ctx, report); err != nil {
			return fmt.Errorf("failed to create report: %w", err)
		}
		log := &entity.ActivityLogs{
			UserID:      userID,
			Action:      "REPORT_CREATED",
			EntityType:  "Reports",
			EntityID:    report.ID,
			Description: "Report created",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}

func (u *assetReportUsecaseImpl) UpdateReport(ctx context.Context, adminID uint, report *entity.Reports, status string) error {
	return u.txRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		exist, err := u.reportRepo.FetchByID(ctx, report.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch report: %w", err)
		}
		asset, err := u.assetRepo.FetchByID(ctx, exist.AssetID)
		if err != nil {
			return fmt.Errorf("failed to fetch asset: %w", err)
		}

		checkedAt := time.Now()
		switch status {
		case "Baik":
			asset.CheckedAt = checkedAt
			if err := u.assetRepo.UpdateCheckedAt(ctx, asset); err != nil {
				return fmt.Errorf("failed to update checked_at asset: %w", err)
			}
		case "Buruk":
			asset.Status = status
			asset.CheckedAt = checkedAt
			if err := u.assetRepo.Update(ctx, asset); err != nil {
				return fmt.Errorf("failed to update asset: %w", err)
			}
		default:
			return entity.ErrInvalidInput
		}

		exist.Status = status
		if report.Message != "" {
			exist.Message = report.Message
		}
		exist.ApprovedByID = &adminID
		now := time.Now()
		exist.ApprovedAt = &now
		if err := u.reportRepo.Update(ctx, exist); err != nil {
			return fmt.Errorf("failed to update report: %w", err)
		}

		log := &entity.ActivityLogs{
			UserID:      adminID,
			Action:      "REPORT_UPDATED",
			EntityType:  "Reports",
			EntityID:    exist.ID,
			Description: "Report updated",
		}
		return u.activityLogRepo.Create(ctx, log)
	})
}

func (u *assetReportUsecaseImpl) GetAll(ctx context.Context) ([]entity.Reports, error) {
	return u.reportRepo.FetchAll(ctx)
}

func (u *assetReportUsecaseImpl) GetByID(ctx context.Context, id uint) (*entity.Reports, error) {
	return u.reportRepo.FetchByID(ctx, id)
}
