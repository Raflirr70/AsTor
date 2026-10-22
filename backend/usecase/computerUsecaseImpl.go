package usecase

import (
	"astor/backend/entity"
	dr "astor/backend/repository/domain"
	du "astor/backend/usecase/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type computerUsecaseImpl struct {
	assetRepo        dr.AssetRepository
	computerRepo     dr.ComputerRepository
	computerCompRepo dr.ComputerComponentRepository
	trxRepo          dr.TransactionRepository
}

func NewComputerUsecaseImpl(assetRepo dr.AssetRepository, computerRepo dr.ComputerRepository, computerCompRepo dr.ComputerComponentRepository, trxRepo dr.TransactionRepository) du.ComputerUsecase {
	return &computerUsecaseImpl{assetRepo: assetRepo, computerRepo: computerRepo, computerCompRepo: computerCompRepo, trxRepo: trxRepo}
}
func (u *computerUsecaseImpl) Create(ctx context.Context, req du.UpdateComputerInput) error {
	return u.trxRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		prefix := du.CategoryComputerCode + "-"
		lastCode, err := u.assetRepo.FindLastCode(ctx, prefix)
		nextNum := 1
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to fetch last asset code: %w", err)
		}
		if lastCode != "" {
			var num int
			if _, err := fmt.Sscanf(lastCode, prefix+"%d", &num); err == nil {
				nextNum = num + 1
			}
		}
		code := fmt.Sprintf("%s%03d", prefix, nextNum)
		asset := &entity.Assets{
			CategoryID: du.CategoryComputerID,
			Code:       code,
			Label:      req.Label,
			Status:     req.Status,
		}

		if err := u.assetRepo.Create(ctx, asset); err != nil {
			return fmt.Errorf("failed to create asset: %w", err)
		}

		computer := &entity.Computers{AssetID: asset.ID}
		if err := u.computerRepo.Create(ctx, computer); err != nil {
			return fmt.Errorf("failed to create computer: %w", err)
		}

		for _, compID := range req.ComponentIDs {
			used, err := u.computerCompRepo.CheckUsed(ctx, compID)
			if err != nil {
				return fmt.Errorf("failed to check component usage: %w", err)
			}
			if used {
				return fmt.Errorf("component asset ID %d is already installed", compID)
			} else {
				component := &entity.ComputerComponents{
					AssetID:    compID,
					ComputerID: computer.ID,
					StartedAt:  time.Now(),
				}
				if err := u.computerCompRepo.Create(ctx, component); err != nil {
					return fmt.Errorf("failed to create computer component: %w", err)
				}
			}
		}
		return nil
	})
}
func (u *computerUsecaseImpl) Update(ctx context.Context, computerID uint, req du.UpdateComputerInput) error {
	return u.trxRepo.WithinTransaction(ctx, func(ctx context.Context) error {
		computer, err := u.computerRepo.FetchByID(ctx, computerID)
		if err != nil {
			return fmt.Errorf("failed to find computer : %w", err)
		}
		assetChanged := false
		if req.Label != "" && computer.Asset.Label != req.Label {
			assetChanged = true
		}
		if req.Status != "" && computer.Asset.Status != req.Status {
			assetChanged = true
		}

		// 3. Cek apakah component berubah — komponen AKTIF (EndedAt.IsZero)
		active := map[uint]bool{} // AssetID aktif existing
		for _, c := range computer.ComputerComponents {
			if c.EndedAt.IsZero() {
				active[c.AssetID] = true
			}
		}
		incoming := map[uint]bool{} // dari req.ComponentIDs (dedupe)
		for _, id := range req.ComponentIDs {
			incoming[id] = true
		}

		compsChanged := len(active) != len(incoming)
		if !compsChanged {
			for id := range incoming {
				if !active[id] {
					compsChanged = true
					break
				}
			}
		}

		// Tidak ada perubahan sama sekali → selesai
		if !assetChanged && !compsChanged {
			return entity.ErrNoChanges
		}
		// 4. Update asset jika berubah (GORM struct Updates: field kosong auto-skip)
		if assetChanged {
			asset := &entity.Assets{ID: computer.Asset.ID, Label: req.Label, Status: req.Status}
			fmt.Printf("%d-%d\n", computer.Asset.Status, asset.Status)
			fmt.Printf("%d-%d\n", computer.Asset.CheckedAt, asset.CheckedAt)
			if req.Status != computer.Asset.Status {
				fmt.Printf("Masuk ke sini\n")
				asset.CheckedAt = time.Now()
			}
			if err := u.assetRepo.Update(ctx, asset); err != nil {
				return fmt.Errorf("failed update asset : %w", err)
			}
		}

		// 5. Diff components
		for _, c := range computer.ComputerComponents {
			// aktif tapi tidak di input → isi endedAt
			if c.EndedAt.IsZero() && !incoming[c.AssetID] {
				c.EndedAt = time.Now()
				if err := u.computerCompRepo.Update(ctx, &c); err != nil {
					return fmt.Errorf("failed update computer component : %w", err)
				}
			}
		}
		for _, id := range req.ComponentIDs {
			if active[id] {
				continue
			} // yang sama → tetap, skip
			// bertambah → cek belum terpakai, lalu buat baru
			used, err := u.computerCompRepo.CheckUsed(ctx, id)
			if err != nil {
				return fmt.Errorf("failed checkUsed computer component : %w", err)

			}
			if used {
				return fmt.Errorf("component asset ID %d is already installed", id)
			}
			if err := u.computerCompRepo.Create(ctx, &entity.ComputerComponents{
				AssetID: id, ComputerID: computer.ID, StartedAt: time.Now(),
			}); err != nil {
				return fmt.Errorf("failed created computer component : %w", err)
			}
		}
		return nil
	})
}
func (u *computerUsecaseImpl) GetAll(ctx context.Context) ([]entity.Computers, error) {
	return u.computerRepo.FetchAll(ctx)
}
