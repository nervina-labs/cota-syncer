package data

import (
	"context"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"gorm.io/gorm/clause"
)

var _ biz.JoyIDInfoRepo = (*joyIDInfoRepo)(nil)

type JoyIDInfo struct {
	ID           uint `gorm:"primaryKey"`
	BlockNumber  uint64
	LockHash     string
	Version      string
	PubKey       string
	CredentialId string
	Alg          string
	CotaCellId   string
	Name         string
	Avatar       string
	Description  string
	Extension    string
	TxIndex      uint32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type JoyIDInfoVersion struct {
	ID             uint `gorm:"primaryKey"`
	OldBlockNumber uint64
	BlockNumber    uint64
	LockHash       string
	OldVersion     string
	Version        string
	PubKey         string
	CredentialId   string
	Alg            string
	CotaCellId     string
	OldName        string
	Name           string
	OldAvatar      string
	Avatar         string
	OldDescription string
	Description    string
	OldExtension   string
	Extension      string
	ActionType     uint8 //	0-create 1-update 2-delete
	TxIndex        uint32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type SubKeyInfo struct {
	ID           uint `gorm:"primaryKey"`
	LockHash     string
	PubKey       string
	CredentialId string
	Alg          string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type joyIDInfoRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewJoyIDInfoRepo(data *Data, logger *logger.Logger) biz.JoyIDInfoRepo {
	return &joyIDInfoRepo{
		data:   data,
		logger: logger,
	}
}

func (repo joyIDInfoRepo) CreateJoyIDInfo(ctx context.Context, joyIDInfo *biz.JoyIDInfo) error {
	if err := repo.data.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "lock_hash"}},
		UpdateAll: true,
	}).Create(joyIDInfo).Error; err != nil {
		return err
	}

	var subKeys []biz.SubKeyInfo
	for _, v := range joyIDInfo.SubKeys {
		subKeys = append(subKeys, biz.SubKeyInfo{
			LockHash:     joyIDInfo.LockHash,
			PubKey:       v.PubKey,
			CredentialId: v.CredentialId,
			Alg:          v.Alg,
		})
	}
	if err := repo.data.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "pub_key"}},
		UpdateAll: true,
	}).Create(subKeys).Error; err != nil {
		return err
	}
	return nil
}

func (repo joyIDInfoRepo) DeleteJoyIDInfo(ctx context.Context, blockNumber uint64) error {
	if err := repo.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(JoyIDInfo{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo joyIDInfoRepo) ParseJoyIDInfo(blockNumber uint64, txIndex uint32, lockScript *ckbTypes.Script, joyIDMeta map[string]any) (joyID biz.JoyIDInfo, err error) {
	lockHash, err := lockScript.Hash()
	if err != nil {
		return
	}
	lockHashStr := lockHash.String()[2:]
	var joyIDInfo biz.JoyIDInfo
	err = mapstructure.Decode(joyIDMeta, &joyIDInfo)
	if err != nil {
		return
	}
	joyID = biz.JoyIDInfo{
		BlockNumber:  blockNumber,
		LockHash:     lockHashStr,
		Version:      joyIDInfo.Version,
		Name:         joyIDInfo.Name,
		Avatar:       joyIDInfo.Avatar,
		Description:  joyIDInfo.Description,
		PubKey:       joyIDInfo.PubKey,
		CredentialId: joyIDInfo.CredentialId,
		Alg:          joyIDInfo.Alg,
		CotaCellId:   joyIDInfo.Alg,
		Extension:    joyIDInfo.Extension,
		SubKeys:      joyID.SubKeys,
		TxIndex:      txIndex,
	}
	return
}
