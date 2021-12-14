package data

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"gorm.io/gorm"
)

var _ biz.ClaimedCotaNftKvPairRepo = (*claimedCotaNftKvPairRepo)(nil)

type ClaimedCotaNftKvPair struct {
	gorm.Model

	BlockNumber uint64
	CotaId      string
	CotaIdCRC   uint32
	Total       uint32
	TokenIndex  uint32
	OutPoint    string
	OutPointCrc uint32
	LockHash    string
	LockHashCrc uint32
}

type claimedCotaNftKvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewClaimedCotaNftKvPairRepo(data *Data, logger *logger.Logger) biz.ClaimedCotaNftKvPairRepo {
	return &claimedCotaNftKvPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp claimedCotaNftKvPairRepo) CreateClaimedCotaNftKvPair(ctx context.Context, c *biz.ClaimedCotaNftKvPair) error {
	if err := rp.data.db.WithContext(ctx).Create(c).Error; err != nil {
		return err
	}
	return nil
}

func (rp claimedCotaNftKvPairRepo) DeleteClaimedCotaNftKvPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(ClaimedCotaNftKvPair{}).Error; err != nil {
		return err
	}
	return nil
}
