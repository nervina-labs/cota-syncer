package data

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"time"
)

var _ biz.CheckInfoRepo = (*checkInfoRepo)(nil)

type CheckInfo struct {
	ID          uint `gorm:"primaryKey"`
	BlockNumber uint64
	BlockHash   string
	CheckType   biz.CheckType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type checkInfoRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewCheckInfoRepo(data *Data, logger *logger.Logger) biz.CheckInfoRepo {
	return &checkInfoRepo{
		data:   data,
		logger: logger,
	}
}

func (rp checkInfoRepo) FindLastCheckInfo(ctx context.Context, info *biz.CheckInfo) error {
	c := &CheckInfo{}
	if err := rp.data.db.WithContext(ctx).Last(&c).Error; err != nil {
		return err
	}
	info.Id = uint64(c.ID)
	info.BlockNumber = c.BlockNumber
	info.BlockHash = c.BlockHash
	return nil
}

func (rp checkInfoRepo) CreateCheckInfo(ctx context.Context, info *biz.CheckInfo) error {
	if err := rp.data.db.WithContext(ctx).Create(CheckInfo{
		BlockNumber: info.BlockNumber,
		BlockHash:   info.BlockHash,
		CheckType:   info.CheckType,
	}).Error; err != nil {
		return err
	}
	return nil
}
