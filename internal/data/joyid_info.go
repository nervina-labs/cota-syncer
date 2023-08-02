package data

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

var _ biz.JoyIDInfoRepo = (*joyIDInfoRepo)(nil)

var ErrInvalidJoyIDInfo = errors.New("JoyID info is invalid")

type JoyIDInfo struct {
	ID                   uint `gorm:"primaryKey"`
	BlockNumber          uint64
	LockHash             string
	Version              string
	PubKey               string
	CredentialId         string
	Alg                  string
	FrontEnd             string
	DeviceName           string
	DeviceType           string
	CotaCellId           string
	Name                 string
	Avatar               string
	Description          string
	Extension            string
	DerivationCId        string
	DerivationCommitment string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type JoyIDInfoVersion struct {
	ID                      uint `gorm:"primaryKey"`
	OldBlockNumber          uint64
	BlockNumber             uint64
	LockHash                string
	OldVersion              string
	Version                 string
	PubKey                  string
	CredentialId            string
	Alg                     string
	OldFrontEnd             string
	FrontEnd                string
	OldDeviceName           string
	DeviceName              string
	OldDeviceType           string
	DeviceType              string
	CotaCellId              string
	OldName                 string
	Name                    string
	OldAvatar               string
	Avatar                  string
	OldDescription          string
	Description             string
	OldExtension            string
	Extension               string
	ActionType              uint8 //	0-create 1-update 2-delete
	TxIndex                 uint32
	DerivationCId           string
	OldDerivationCId        string
	DerivationCommitment    string
	OldDerivationCommitment string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type SubKeyInfo struct {
	ID                   uint `gorm:"primaryKey"`
	LockHash             string
	BlockNumber          uint64
	PubKey               string
	CredentialId         string
	Alg                  string
	FrontEnd             string
	DeviceName           string
	DeviceType           string
	DerivationCId        string
	DerivationCommitment string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type SubKeyInfoVersion struct {
	ID                      uint `gorm:"primaryKey"`
	OldBlockNumber          uint64
	BlockNumber             uint64
	LockHash                string
	PubKey                  string
	CredentialId            string
	Alg                     string
	OldFrontEnd             string
	FrontEnd                string
	OldDeviceName           string
	DeviceName              string
	OldDeviceType           string
	DeviceType              string
	ActionType              uint8 //	0-create 1-update 2-delete
	TxIndex                 uint32
	DerivationCId           string
	OldDerivationCId        string
	DerivationCommitment    string
	OldDerivationCommitment string
	CreatedAt               time.Time
	UpdatedAt               time.Time
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

func (repo joyIDInfoRepo) DeleteJoyIDInfo(ctx context.Context, blockNumber uint64) error {
	if err := repo.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(JoyIDInfo{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo joyIDInfoRepo) ParseJoyIDInfo(ctx context.Context, blockNumber uint64, txIndex uint32, lockScript *ckbTypes.Script, joyIDMeta map[string]any) (joyID biz.JoyIDInfo, err error) {
	lockHash, err := lockScript.Hash()
	if err != nil {
		return
	}
	lockHashStr := lockHash.String()[2:]
	var joyIDInfo biz.JoyIDInfoJson
	err = mapstructure.Decode(joyIDMeta, &joyIDInfo)
	if err != nil {
		return
	}
	if lenWithout0x(joyIDInfo.PubKey) > 128 || lenWithout0x(joyIDInfo.CotaCellId) > 16 || lenWithout0x(joyIDInfo.Alg) > 2 {
		err = ErrInvalidJoyIDInfo
		return
	}
	if len(joyIDInfo.Name) > 240 || len(joyIDInfo.Avatar) > 500 || len(joyIDInfo.Description) > 1000 {
		err = ErrInvalidJoyIDInfo
		return
	}
	subKeys := make([]biz.SubKeyInfo, len(joyIDInfo.SubKeys))
	for i, v := range joyIDInfo.SubKeys {
		if lenWithout0x(v.PubKey) > 128 || lenWithout0x(v.Alg) > 2 {
			err = ErrInvalidJoyIDInfo
		}
		subKeys[i] = biz.SubKeyInfo{
			PubKey:               remove0x(v.PubKey),
			CredentialId:         remove0x(v.CredentialId),
			Alg:                  remove0x(v.Alg),
			FrontEnd:             v.FrontEnd,
			DeviceName:           v.DeviceName,
			DeviceType:           v.DeviceType,
			DerivationCId:        remove0x(v.Derivation.CredentialId),
			DerivationCommitment: remove0x(v.Derivation.Commitment),
		}
	}
	joyID = biz.JoyIDInfo{
		BlockNumber:          blockNumber,
		LockHash:             lockHashStr,
		Version:              joyIDInfo.Version,
		Name:                 joyIDInfo.Name,
		Avatar:               joyIDInfo.Avatar,
		Description:          joyIDInfo.Description,
		PubKey:               remove0x(joyIDInfo.PubKey),
		CredentialId:         remove0x(joyIDInfo.CredentialId),
		Alg:                  remove0x(joyIDInfo.Alg),
		FrontEnd:             joyIDInfo.FrontEnd,
		DeviceName:           joyIDInfo.DeviceName,
		DeviceType:           joyIDInfo.DeviceType,
		CotaCellId:           remove0x(joyIDInfo.CotaCellId),
		Extension:            joyIDInfo.Extension,
		SubKeys:              subKeys,
		TxIndex:              txIndex,
		DerivationCId:        remove0x(joyIDInfo.Derivation.CredentialId),
		DerivationCommitment: remove0x(joyIDInfo.Derivation.Commitment),
	}
	return
}

func remove0x(value string) string {
	if strings.HasPrefix(value, "0x") {
		return value[2:]
	}
	return value
}

func lenWithout0x(value string) int {
	return len(remove0x(value))
}
