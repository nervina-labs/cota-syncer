package data

import (
	"context"
	"encoding/hex"
	"hash/crc32"
	"time"

	"github.com/nervina-labs/cota-smt-go/smt"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
)

var _ biz.ExtensionPairRepo = (*extensionPairRepo)(nil)

type ExtensionPair struct {
	ID          uint `gorm:"primaryKey"`
	BlockNumber uint64
	LockHash    string
	LockHashCRC uint32
	Key         string
	Value       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ExtensionPairVersion struct {
	ID             uint `gorm:"primaryKey"`
	OldBlockNumber uint64
	BlockNumber    uint64
	Key            string
	Value          string
	OldValue       string
	LockHash       string
	ActionType     uint8 //	0-create 1-update 2-delete
	TxIndex        uint32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type extensionPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewExtensionPairRepo(data *Data, logger *logger.Logger) biz.ExtensionPairRepo {
	return &extensionPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp extensionPairRepo) CreateExtensionPair(ctx context.Context, ext *biz.ExtensionPair) error {
	if err := rp.data.db.WithContext(ctx).Create(ext).Error; err != nil {
		return err
	}
	return nil
}

func (rp extensionPairRepo) DeleteExtensionPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Error; err != nil {
		return err
	}
	return nil
}

func (rp extensionPairRepo) ParseExtensionPairs(blockNumber uint64, entry biz.Entry) (pairs []biz.ExtensionPair, err error) {
	entries := smt.ExtensionEntriesFromSliceUnchecked(entry.InputType[1:])
	extensionLeafKeys := entries.Leaves().Keys()
	extensionLeafValues := entries.Leaves().Values()
	lockHash, err := entry.LockScript.Hash()
	if err != nil {
		return
	}
	lockHashStr := lockHash.String()[2:]
	lockHashCRC32 := crc32.ChecksumIEEE([]byte(lockHashStr))
	for i := uint(0); i < extensionLeafKeys.Len(); i++ {
		key := extensionLeafKeys.Get(i)
		value := extensionLeafValues.Get(i)
		pairs = append(pairs, biz.ExtensionPair{
			BlockNumber: blockNumber,
			Key:         hex.EncodeToString(key.RawData()),
			Value:       hex.EncodeToString(value.RawData()),
			LockHash:    lockHashStr,
			LockHashCRC: lockHashCRC32,
			TxIndex:     entry.TxIndex,
		})
	}
	return
}
