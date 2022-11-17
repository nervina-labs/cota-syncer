package data

import (
	"context"
	"encoding/hex"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	"github.com/nervina-labs/cota-smt-go/smt"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
)

var _ biz.ExtensionPairRepo = (*extensionPairRepo)(nil)

type ExtensionKvPair struct {
	ID          uint `gorm:"primaryKey"`
	BlockNumber uint64
	LockHash    string
	LockHashCRC uint32
	Key         string
	Value       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ExtensionKvPairVersion struct {
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

func NewExtensionKvPairRepo(data *Data, logger *logger.Logger) biz.ExtensionPairRepo {
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

func (rp extensionPairRepo) ParseExtensionPairs(blockNumber uint64, entry biz.Entry) (pairs biz.ExtensionPairs, err error) {
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
		pairs.Extensions = append(pairs.Extensions, biz.ExtensionPair{
			BlockNumber: blockNumber,
			Key:         hex.EncodeToString(key.RawData()),
			Value:       hex.EncodeToString(value.RawData()),
			LockHash:    lockHashStr,
			LockHashCRC: lockHashCRC32,
			TxIndex:     entry.TxIndex,
		})
	}

	switch string(entries.SubType().RawData()) {
	case "subkey":
		var subKeys []biz.SubKeyPair
		if subKeys, err = rp.parseSubKeyPairs(entries, blockNumber, lockHashStr); err != nil {
			return biz.ExtensionPairs{}, err
		}
		pairs.SubKeys = append(pairs.SubKeys, subKeys...)
	case "social":
		var socialKey *biz.SocialKvPair
		if socialKey, err = rp.parseSocialKeyPairs(entries, blockNumber, lockHashStr); err != nil {
			return biz.ExtensionPairs{}, err
		}

		if socialKey != nil {
			pairs.SocialKeys = append(pairs.SocialKeys, *socialKey)
		}
	}

	return
}

func (rp extensionPairRepo) parseSubKeyPairs(entries *smt.ExtensionEntries, blockNumber uint64, lockHash string) ([]biz.SubKeyPair, error) {
	var (
		extData, algIndex int64
		subKeys           []biz.SubKeyPair
		err               error
	)

	subKeyEntries := smt.SubKeyEntriesFromSliceUnchecked(entries.RawData().RawData())
	subKeyLeafKeys := subKeyEntries.Keys()
	subKeyLeafValues := subKeyEntries.Values()
	for i := uint(0); i < subKeyLeafKeys.Len(); i++ {
		key := subKeyLeafKeys.Get(i)
		value := subKeyLeafValues.Get(i)

		if extData, err = strconv.ParseInt(hex.EncodeToString(key.ExtData().RawData()), 10, 32); err != nil {
			return nil, err
		}
		if algIndex, err = strconv.ParseInt(hex.EncodeToString(value.AlgIndex().RawData()), 10, 16); err != nil {
			return nil, err
		}
		subKeys = append(subKeys, biz.SubKeyPair{
			BlockNumber: blockNumber,
			LockHash:    lockHash,
			SubType:     string(key.SubType().RawData()),
			ExtData:     uint32(extData),
			AlgIndex:    uint16(algIndex),
			PubkeyHash:  remove0x(hex.EncodeToString(value.PubkeyHash().RawData())),
			UpdatedAt:   time.Now().UTC(),
		})
	}

	return subKeys, nil
}

func (rp extensionPairRepo) parseSocialKeyPairs(entries *smt.ExtensionEntries, blockNumber uint64, lockHash string) (*biz.SocialKvPair, error) {
	var (
		recoveryMode, must, total int64
		signers                   []string
		err                       error
	)

	socialEntry := smt.SocialEntryFromSliceUnchecked(entries.RawData().RawData())
	socialLeafValues := socialEntry.Value()
	if socialLeafValues == nil {
		return nil, nil
	}

	if recoveryMode, err = strconv.ParseInt(hex.EncodeToString(socialLeafValues.RecoveryMode().AsSlice()), 10, 8); err != nil {
		return nil, err
	}
	if must, err = strconv.ParseInt(hex.EncodeToString(socialLeafValues.Must().AsSlice()), 10, 8); err != nil {
		return nil, err
	}
	if total, err = strconv.ParseInt(hex.EncodeToString(socialLeafValues.Total().AsSlice()), 10, 8); err != nil {
		return nil, err
	}

	lockScriptVec := socialLeafValues.Signers()
	for i := uint(0); i < socialLeafValues.Signers().Len(); i++ {
		signer := lockScriptVec.Get(i).RawData()
		signers = append(signers, remove0x(hex.EncodeToString(signer)))
	}

	return &biz.SocialKvPair{
		BlockNumber:  blockNumber,
		LockHash:     lockHash,
		LockHashCRC:  crc32.ChecksumIEEE([]byte(lockHash)),
		RecoveryMode: uint8(recoveryMode),
		Must:         uint8(must),
		Total:        uint8(total),
		Signers:      strings.Join(signers, ","),
		UpdatedAt:    time.Now().UTC(),
	}, nil
}
