package data

import (
	"context"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"gorm.io/gorm/clause"
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

type IssuerInfoVersion struct {
	ID              uint `gorm:"primaryKey"`
	OldBlockNumber  uint64
	BlockNumber     uint64
	LockHash        string
	OldVersion      string
	Version         string
	OldName         string
	Name            string
	OldAvatar       string
	Avatar          string
	OldDescription  string
	Description     string
	OldLocalization string
	Localization    string
	ActionType      uint8 //	0-create 1-update 2-delete
	TxIndex         uint32
	CreatedAt       time.Time
	UpdatedAt       time.Time
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

func (repo issuerInfoRepo) ParseIssuerInfo(blockNumber uint64, txIndex uint32, lockScript *ckbTypes.Script, issuerMeta map[string]any) (issuer biz.IssuerInfo, err error) {
	lockHash, err := lockScript.Hash()
	if err != nil {
		return
	}
	lockHashStr := lockHash.String()[2:]
	var issuerInfo biz.IssuerInfoJson
	err = mapstructure.Decode(issuerMeta, &issuerInfo)
	if err != nil {
		return
	}
	localization, err := json.Marshal(issuerInfo.Localization)
	if err != nil {
		return
	}
	localizationStr := string(localization)
	if localizationStr == "{}" {
		localizationStr = ""
	}
	issuer = biz.IssuerInfo{
		BlockNumber:  blockNumber,
		LockHash:     lockHashStr,
		Version:      issuerInfo.Version,
		Name:         issuerInfo.Name,
		Avatar:       issuerInfo.Avatar,
		Description:  issuerInfo.Description,
		Localization: localizationStr,
		TxIndex:      txIndex,
	}
	return
}
