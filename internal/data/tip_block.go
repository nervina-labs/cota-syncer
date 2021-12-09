package data

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
	"gorm.io/gorm"
)

type TipBlock struct {
	gorm.Model

	BlockNumber uint64
	BlockHash   string
}

type tipBlockRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewTipBlockRepo(data *Data, logger *logger.Logger) biz.TipBlockRepo {
	return &tipBlockRepo{
		data:   data,
		logger: logger,
	}
}

func (t tipBlockRepo) FindOrCreateTipBlock(ctx context.Context, tipBlock *biz.TipBlock) error {
	if err := t.data.db.WithContext(ctx).FirstOrCreate(tipBlock, TipBlock{BlockNumber: tipBlock.BlockNumber}).Error; err != nil {
		return err
	}
	return nil
}

func (t tipBlockRepo) UpdateTipBlock(ctx context.Context, tipBlock *biz.TipBlock) error {
	if err := t.data.db.WithContext(ctx).Save(tipBlock).Error; err != nil {
		return err
	}
	return nil
}
