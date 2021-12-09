package data

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
	"gorm.io/gorm"
)

var _ biz.DefineCotaNftKVPairRepo = (*defineCotaNftKVPairRepo)(nil)

type DefineCotaNftKVPair struct {
	gorm.Model

	BlockNumber uint64
	CotaId      string
	CotaIdCRC   uint32
	Total       uint32
	Issued      uint32
	Configure   uint8
	LockHash    string
	LockHashCRC uint32
}

func NewDefineCotaNftKVPairRepo(data *Data, logger *logger.Logger) *defineCotaNftKVPairRepo {
	return &defineCotaNftKVPairRepo{
		data:   data,
		logger: logger,
	}
}

type defineCotaNftKVPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func (rp defineCotaNftKVPairRepo) CreateDefineCotaNftKVPair(ctx context.Context, d *biz.DefineCotaNftKVPair) error {
	if err := rp.data.db.WithContext(ctx).Create(d).Error; err != nil {
		return err
	}
	return nil
}

func (rp defineCotaNftKVPairRepo) DeleteDefineCotaNftKVPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Error; err != nil {
		return err
	}
	return nil
}
