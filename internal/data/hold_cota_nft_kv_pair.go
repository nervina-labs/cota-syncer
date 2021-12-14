package data

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"gorm.io/gorm"
)

var _ biz.HoldCotaNftKvPairRepo = (*holdCotaNftKvPairRepo)(nil)

type HoldCotaNftKvPair struct {
	gorm.Model

	BlockNumber    uint64
	CotaId         string
	CotaIdCRC      uint32
	TokenIndex     uint32
	State          uint8
	Configure      uint8
	Characteristic string
	LockHash       string
	LockHashCRC    uint32
}

type holdCotaNftKvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewHoldCotaNftKvPairRepo(data *Data, logger *logger.Logger) biz.HoldCotaNftKvPairRepo {
	return &holdCotaNftKvPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp holdCotaNftKvPairRepo) CreateHoldCotaNftKvPair(ctx context.Context, h *biz.HoldCotaNftKvPair) error {
	if err := rp.data.db.WithContext(ctx).Create(h).Error; err != nil {
		return err
	}
	return nil
}

func (rp holdCotaNftKvPairRepo) DeleteHoldCotaNftKvPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(HoldCotaNftKvPair{}).Error; err != nil {
		return err
	}
	return nil
}
