package data

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
	"gorm.io/gorm"
)

var _ biz.WithdrawCotaNftKvPairRepo = (*withdrawCotaNftKvPairRepo)(nil)

type WithdrawCotaNftKvPair struct {
	gorm.Model

	BlockNumber         uint64
	CotaId              string
	CotaIdCRC           uint32
	Total               uint32
	TokenIndex          uint32
	OutPoint            string
	OutPointCrc         uint32
	State               uint8
	Configure           uint8
	Characteristic      string
	ReceiverLockHash    string
	ReceiverLockHashCrc uint32
	LockHash            string
	LockHashCrc         uint32
}

func NewWithdrawCotaNftKvPairRepo(data *Data, logger *logger.Logger) *withdrawCotaNftKvPairRepo {
	return &withdrawCotaNftKvPairRepo{
		data:   data,
		logger: logger,
	}
}

type withdrawCotaNftKvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func (rp withdrawCotaNftKvPairRepo) CreateWithdrawCotaNftKvPair(ctx context.Context, w *biz.WithdrawCotaNftKvPair) error {
	if err := rp.data.db.WithContext(ctx).Create(w).Error; err != nil {
		return err
	}
	return nil
}

func (rp withdrawCotaNftKvPairRepo) DeleteWithdrawCotaNftKvPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(WithdrawCotaNftKvPair{}).Error; err != nil {
		return err
	}
	return nil
}
