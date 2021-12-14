package data

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
	"gorm.io/gorm"
)

var _ biz.CheckInfoRepo = (*checkInfoRepo)(nil)

type CheckInfo struct {
	gorm.Model

	BlockNumber uint64
	BlockHash   string
	CheckType   biz.CheckType
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

func (rp checkInfoRepo) FindOrCreateCheckInfo(ctx context.Context, info *biz.CheckInfo) error {
	if err := rp.data.db.WithContext(ctx).FirstOrCreate(info, CheckInfo{BlockNumber: info.BlockNumber, CheckType: info.CheckType}).Error; err != nil {
		return err
	}
	return nil
}

func (rp checkInfoRepo) UpdateCheckInfo(ctx context.Context, info biz.CheckInfo) error {
	if err := rp.data.db.WithContext(ctx).Save(info).Error; err != nil {
		return err
	}
	return nil
}
