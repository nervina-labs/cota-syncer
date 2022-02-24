package data

import (
	"context"
	"encoding/json"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"hash/crc32"
	"time"
)

var _ biz.IssuerInfoRepo = (*issuerInfoRepo)(nil)

type IssuerInfo struct {
	ID           uint `gorm:"primaryKey"`
	BlockNumber  uint64
	LockHash     string
	LockHashCRC  uint32
	Version      string
	Name         string
	Avatar       string
	Description  string
	Localization string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type IssuerInfoJson struct {
	Version      string
	Name         string
	Avatar       string
	Description  string
	Localization string
}

type issuerInfoRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewIssuerInfoRepo(data *Data, logger *logger.Logger) biz.IssuerInfoRepo {
	return &issuerInfoRepo{
		data:   data,
		logger: logger,
	}
}

func (repo issuerInfoRepo) CreateIssuerInfo(ctx context.Context, issuer *biz.IssuerInfo) error {
	db := repo.data.db.WithContext(ctx)
	var dest biz.IssuerInfo
	rows := db.First(dest, "lock_hash_crc = ? AND lock_hash = ?", issuer.LockHashCRC, issuer.LockHash)
	if rows.RowsAffected > 0 {
		return nil
	}
	if err := repo.data.db.WithContext(ctx).Create(issuer).Error; err != nil {
		return err
	}
	return nil
}

func (repo issuerInfoRepo) DeleteIssuerInfo(ctx context.Context, blockNumber uint64) error {
	if err := repo.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(IssuerInfo{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo issuerInfoRepo) ParseIssuerInfo(blockNumber uint64, entry biz.Entry) (issuer biz.IssuerInfo, err error) {
	var issuerJson IssuerInfoJson
	err = json.Unmarshal(entry.OutputType, &issuerJson)
	if err != nil {
		return
	}
	lockHash, err := entry.LockScript.Hash()
	if err != nil {
		return
	}
	lockHashStr := lockHash.String()[2:]
	lockHashCRC32 := crc32.ChecksumIEEE([]byte(lockHashStr))
	issuer = biz.IssuerInfo{
		BlockNumber:  blockNumber,
		LockHash:     lockHashStr,
		LockHashCRC:  lockHashCRC32,
		Version:      issuerJson.Version,
		Name:         issuerJson.Name,
		Avatar:       issuerJson.Avatar,
		Description:  issuerJson.Description,
		Localization: issuerJson.Localization,
	}
	return
}
