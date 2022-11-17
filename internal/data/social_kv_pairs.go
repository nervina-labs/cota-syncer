package data

import (
	"context"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"time"
)

type SocialKvPair struct {
	ID           uint `gorm:"primaryKey"`
	BlockNumber  uint64
	LockHash     string
	LockHashCRC  uint32
	RecoveryMode uint8
	Must         uint8
	Total        uint8
	Signers      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SocialKvPairVersion struct {
	ID              uint `gorm:"primaryKey"`
	OldBlockNumber  uint64
	BlockNumber     uint64
	LockHash        string
	OldRecoveryMode uint8
	RecoveryMode    uint8
	OldMust         uint8
	Must            uint8
	OldTotal        uint8
	Total           uint8
	OldSigners      string
	Signers         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

var _ biz.SocialPairRepo = (*socialKeyPairRepo)(nil)

type socialKeyPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewSocialKeyKvPairRepo(data *Data, logger *logger.Logger) biz.SubKeyPairRepo {
	return &subKeyPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp socialKeyPairRepo) CreateSocialKeyPair(ctx context.Context, social *biz.SocialKvPair) error {
	if err := rp.data.db.WithContext(ctx).Create(social).Error; err != nil {
		return err
	}
	return nil
}

func (rp socialKeyPairRepo) DeleteSocialKeyPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Error; err != nil {
		return err
	}
	return nil
}
