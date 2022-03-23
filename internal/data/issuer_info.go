package data

import (
	"context"
	"encoding/json"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"gorm.io/gorm/clause"
	"hash/crc32"
	"time"
)

var _ biz.IssuerInfoRepo = (*issuerInfoRepo)(nil)

type IssuerInfo struct {
	ID           uint `gorm:"primaryKey"`
	BlockNumber  uint64
	LockHash     string
	Version      string
	Name         string
	Avatar       string
	Description  string
	Localization string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

func (repo issuerInfoRepo) CreateIssuerInfo(ctx context.Context, issuerInfo *biz.IssuerInfo) error {
	if err := repo.data.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "lock_hash"}},
		UpdateAll: true,
	}).Create(issuerInfo).Error; err != nil {
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

func (repo issuerInfoRepo) ParseIssuerInfo(blockNumber uint64, lockScript *ckbTypes.Script, issuerMeta []byte) (issuer biz.IssuerInfo, err error) {
	lockHash, err := lockScript.Hash()
	if err != nil {
		return
	}
	lockHashStr := lockHash.String()[2:]
	lockHashCRC32 := crc32.ChecksumIEEE([]byte(lockHashStr))
	var issuerJson biz.IssuerInfoJson
	err = json.Unmarshal(issuerMeta, &issuerJson)
	if err != nil {
		return
	}
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
