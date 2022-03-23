package data

import (
	"context"
	"encoding/json"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"gorm.io/gorm/clause"
	"time"
)

var _ biz.ClassInfoRepo = (*classInfoRepo)(nil)

type ClassInfo struct {
	ID           uint `gorm:"primaryKey"`
	BlockNumber  uint64
	CotaId       string
	Version      string
	Name         string
	Symbol       string
	Description  string
	Image        string
	Audio        string
	Video        string
	Model        string
	Schema       string
	Properties   string
	Localization string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

func (repo classInfoRepo) ParseClassInfo(blockNumber uint64, classMeta []byte) (class biz.ClassInfo, err error) {
	var classJson biz.ClassInfoJson
	err = json.Unmarshal(classMeta, &classJson)
	if err != nil {
		return
	}
	class = biz.ClassInfo{
		BlockNumber:  blockNumber,
		CotaId:       classJson.CotaId,
		Version:      classJson.Version,
		Name:         classJson.Name,
		Symbol:       classJson.Symbol,
		Description:  classJson.Description,
		Image:        classJson.Image,
		Audio:        classJson.Audio,
		Video:        classJson.Video,
		Model:        classJson.Model,
		Properties:   classJson.Properties,
		Localization: classJson.Localization,
	}
	return
}
