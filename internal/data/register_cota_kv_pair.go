package data

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
	"gorm.io/gorm"
)

var _ biz.RegisterCotaKVPairRepo = (*registerCotaKVPairRepo)(nil)

type RegisterCotaKVPair struct {
	gorm.Model

	BlockNumber uint64
	LockHash    string
	LockHashCRC uint32
}

func NewRegisterCotaKVPairRepo(data *Data, logger *logger.Logger) *registerCotaKVPairRepo {
	return &registerCotaKVPairRepo{
		data:   data,
		logger: logger,
	}
}

type registerCotaKVPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func (rp registerCotaKVPairRepo) CreateRegisterCotaKVPair(ctx context.Context, register *biz.RegisterCotaKVPair) error {
	if err := rp.data.db.WithContext(ctx).Create(register).Error; err != nil {
		return err
	}
	return nil
}

func (rp registerCotaKVPairRepo) DeleteRegisterCotaKVPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(RegisterCotaKVPair{}).Error; err != nil {
		return err
	}
	return nil
}
