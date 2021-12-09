package data

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
)

var _ biz.HoldCotaNftKVPairRepo = (*holdCotaNftKVPairRepo)(nil)

type HoldCotaNftKVPair struct {
	BlockNumber    uint64
	CotaId         string
	CotaIdCRC      uint32
	Total          uint32
	TokenIndex     uint32
	State          uint8
	Configure      uint8
	Characteristic string
	LockHash       string
	LockHashCRC    uint32
}

func NewHoldCotaNftKVPairRepo(data *Data, logger *logger.Logger) *holdCotaNftKVPairRepo {
	return &holdCotaNftKVPairRepo{
		data:   data,
		logger: logger,
	}
}

type holdCotaNftKVPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func (rp holdCotaNftKVPairRepo) CreateHoldCotaNftKVPair(ctx context.Context, h *biz.HoldCotaNftKVPair) error {
	if err := rp.data.db.WithContext(ctx).Create(h).Error; err != nil {
		return err
	}
	return nil
}

func (rp holdCotaNftKVPairRepo) DeleteHoldCotaNftKVPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(HoldCotaNftKVPair{}).Error; err != nil {
		return err
	}
	return nil
}
