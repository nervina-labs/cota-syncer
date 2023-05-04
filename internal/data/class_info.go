package data

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"gorm.io/gorm/clause"
)

var _ biz.ClassInfoRepo = (*classInfoRepo)(nil)

var ErrInvalidClassInfo = errors.New("class info is invalid")

const CotaIdLen = 42

type Audio struct {
	ID        uint `gorm:"primaryKey"`
	Url       string
	Name      string
	CotaId    string
	Idx       uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ClassInfo struct {
	ID             uint `gorm:"primaryKey"`
	BlockNumber    uint64
	CotaId         string
	Version        string
	Name           string
	Symbol         string
	Description    string
	Image          string
	Audio          string
	Video          string
	Model          string
	Characteristic string
	Properties     string
	Localization   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ClassInfoVersion struct {
	ID                uint `gorm:"primaryKey"`
	OldBlockNumber    uint64
	BlockNumber       uint64
	CotaId            string
	OldVersion        string
	Version           string
	OldName           string
	Name              string
	OldSymbol         string
	Symbol            string
	OldDescription    string
	Description       string
	OldImage          string
	Image             string
	OldAudio          string
	Audio             string
	OldVideo          string
	Video             string
	OldModel          string
	Model             string
	OldCharacteristic string
	Characteristic    string
	OldProperties     string
	Properties        string
	OldLocalization   string
	Localization      string
	ActionType        uint8 //	0-create 1-update 2-delete
	TxIndex           uint32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type classInfoRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewClassInfoRepo(data *Data, logger *logger.Logger) biz.ClassInfoRepo {
	return &classInfoRepo{
		data:   data,
		logger: logger,
	}
}

func (repo classInfoRepo) CreateClassInfo(ctx context.Context, classInfo *biz.ClassInfo) error {
	if err := repo.data.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "cota_id"}},
		UpdateAll: true,
	}).Create(classInfo).Error; err != nil {
		return err
	}
	return nil
}

func (repo classInfoRepo) DeleteClassInfo(ctx context.Context, blockNumber uint64) error {
	if err := repo.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(ClassInfo{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo classInfoRepo) ParseClassInfo(blockNumber uint64, txIndex uint32, classMeta map[string]any) (class biz.ClassInfo, err error) {
	var classInfo biz.ClassInfoJson
	err = mapstructure.Decode(classMeta, &classInfo)
	if err != nil {
		return
	}
	characteristic, err := json.Marshal(classInfo.Characteristic)
	if len(classInfo.CotaId) != CotaIdLen {
		repo.logger.Infof(context.TODO(), "class info: %+v", classInfo)
		err = ErrInvalidClassInfo
	}
	if err != nil {
		return
	}
	characteristicStr := string(characteristic)
	if characteristicStr == "null" {
		characteristicStr = ""
	}
	properties, err := json.Marshal(classInfo.Properties)
	if err != nil {
		return
	}
	propertiesStr := string(properties)
	if propertiesStr == "null" {
		propertiesStr = ""
	}

	audios := make([]biz.Audio, len(classInfo.Audios))
	for i, audio := range classInfo.Audios {
		audios[i] = biz.Audio{
			CotaId: classInfo.CotaId[2:],
			Url:    audio.Url,
			Name:   audio.Name,
			Idx:    uint32(i),
		}
	}

	localization, err := json.Marshal(classInfo.Localization)
	if err != nil {
		return
	}
	localizationStr := string(localization)
	if localizationStr == "{}" {
		localizationStr = ""
	}
	class = biz.ClassInfo{
		BlockNumber:    blockNumber,
		CotaId:         classInfo.CotaId[2:],
		Version:        classInfo.Version,
		Name:           classInfo.Name,
		Symbol:         classInfo.Symbol,
		Description:    classInfo.Description,
		Image:          classInfo.Image,
		Audio:          classInfo.Audio,
		Audios:         audios,
		Video:          classInfo.Video,
		Model:          classInfo.Model,
		Properties:     propertiesStr,
		Localization:   localizationStr,
		Characteristic: characteristicStr,
		TxIndex:        txIndex,
	}
	return
}
