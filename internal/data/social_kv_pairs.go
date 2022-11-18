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
	ActionType      uint8 //	0-create 1-update
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

var _ biz.SocialPairRepo = (*socialPairRepo)(nil)

type socialPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewSocialKvPairRepo(data *Data, logger *logger.Logger) biz.SubKeyPairRepo {
	return &subKeyPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp socialPairRepo) CreateSocialPair(ctx context.Context, social *biz.SocialKvPair) error {
	if err := rp.data.db.WithContext(ctx).Create(social).Error; err != nil {
		return err
	}
	return nil
}

func (rp socialPairRepo) DeleteSocialPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(SocialKvPair{}).Error; err != nil {
		return err
	}
	return nil
}
