package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"gorm.io/gorm/clause"
	"regexp"
	"strings"
	"time"
)

var _ biz.JoyIDInfoRepo = (*joyIDInfoRepo)(nil)

var ErrInvalidJoyIDInfo = errors.New("JoyID info is invalid")

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
	Nickname     string
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
	OldNickname    string
	Nickname       string
	ActionType     uint8 //	0-create 1-update 2-delete
	TxIndex        uint32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type SubKeyInfo struct {
	ID           uint `gorm:"primaryKey"`
	LockHash     string
	BlockNumber  uint64
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
			BlockNumber:  joyIDInfo.BlockNumber,
			LockHash:     joyIDInfo.LockHash,
			PubKey:       remove0x(v.PubKey),
			CredentialId: remove0x(v.CredentialId),
			Alg:          remove0x(v.Alg),
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

func (repo joyIDInfoRepo) ParseJoyIDInfo(ctx context.Context, blockNumber uint64, txIndex uint32, lockScript *ckbTypes.Script, joyIDMeta map[string]any) (joyID biz.JoyIDInfo, err error) {
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
		if lenWithout0x(v.PubKey) > 128 || len(v.Alg) > 4 {
			err = ErrInvalidJoyIDInfo
		}
		subKeys[i] = biz.SubKeyInfo{
			PubKey:       remove0x(v.PubKey),
			CredentialId: remove0x(v.CredentialId),
			Alg:          remove0x(v.Alg),
		}
	}
	nickname, err := repo.parseNickname(ctx, joyIDInfo.Name, lockHashStr)
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
		PubKey:       remove0x(joyIDInfo.PubKey),
		CredentialId: remove0x(joyIDInfo.CredentialId),
		Alg:          remove0x(joyIDInfo.Alg),
		CotaCellId:   remove0x(joyIDInfo.CotaCellId),
		Extension:    joyIDInfo.Extension,
		Nickname:     nickname,
		SubKeys:      subKeys,
		TxIndex:      txIndex,
	}
	return
}

func (repo joyIDInfoRepo) parseNickname(ctx context.Context, name string, lockHash string) (nickname string, err error) {
	var registry RegisterCotaKvPair
	if err = repo.data.db.WithContext(ctx).Select("cota_cell_id").Where("lock_hash = ?", lockHash).First(&registry).Error; err != nil {
		return
	}
	match, err := regexp.MatchString(`^[A-Za-z0-9]{4,240}$`, name)
	if err != nil {
		return
	}
	var realName = name
	if !match {
		realName = "noname"
	}
	nickname = realName + "#" + fmt.Sprintf("%04d", registry.CotaCellID%10000)

	var joyIDInfos []JoyIDInfo
	if err = repo.data.db.Model(JoyIDInfo{}).WithContext(ctx).Where("lock_hash <> ? and nickname = ?", lockHash, nickname).Find(&joyIDInfos).Error; err != nil {
		return
	}
	if len(joyIDInfos) == 0 {
		return
	} else {
		nickname = realName + "#" + fmt.Sprintf("%06d", registry.CotaCellID%1000000)
		if err = repo.data.db.Model(JoyIDInfo{}).WithContext(ctx).Where("lock_hash <> ? and nickname = ?", lockHash, nickname).Find(&joyIDInfos).Error; err != nil {
			return
		}
		if len(joyIDInfos) == 0 {
			return
		} else {
			nickname = realName + "#" + fmt.Sprintf("%08d", registry.CotaCellID%100000000)
			if err = repo.data.db.Model(JoyIDInfo{}).WithContext(ctx).Where("lock_hash <> ? and nickname = ?", lockHash, nickname).Find(&joyIDInfos).Error; err != nil {
				return
			}
			if len(joyIDInfos) == 0 {
				return
			} else {
				nickname = realName + "#" + fmt.Sprintf("%010d", registry.CotaCellID%10000000000)
			}
		}
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
