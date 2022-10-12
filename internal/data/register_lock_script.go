package data

import (
	"context"
	"hash/crc32"

	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
)

var _ biz.RegisterLockScriptRepo = (*registerLockScriptRepo)(nil)

type registerLockScriptRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewRegisterLockScriptRepo(data *Data, logger *logger.Logger) biz.RegisterLockScriptRepo {
	return &registerLockScriptRepo{
		data:   data,
		logger: logger,
	}
}

func (rp registerLockScriptRepo) AddRegisterLock(ctx context.Context, lockHash string, lockScriptId uint) error {
	if err := rp.data.db.WithContext(ctx).Model(RegisterCotaKvPair{}).Where("lock_hash = ?", lockHash).Updates(RegisterCotaKvPair{LockScriptId: lockScriptId}).Error; err != nil {
		return err
	}
	return nil
}

func (rp registerLockScriptRepo) IsAllHaveLock(ctx context.Context) (bool, error) {
	var count int64
	if err := rp.data.db.WithContext(ctx).Model(RegisterCotaKvPair{}).Count(&count).Where("lock_script_id = 3094967296").Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func (rp registerLockScriptRepo) FindRegisterQueryInfos(ctx context.Context, page int, pageSize int) ([]biz.RegisterQueryInfo, error) {
	var (
		registerPairs        []RegisterCotaKvPair
		registerQueryInfos   []biz.RegisterQueryInfo
	)
	result := rp.data.db.WithContext(ctx).Model(RegisterCotaKvPair{}).Select("block_number, lock_hash").Where("lock_script_id = 3094967296").Limit(pageSize).Offset(page * pageSize).Find(&registerPairs);
	if result.Error != nil {
		return registerQueryInfos, result.Error
	}
	for _, v := range registerPairs {
		registerQueryInfos = append(registerQueryInfos, biz.RegisterQueryInfo{
			BlockNumber: v.BlockNumber,
			LockHash:    v.LockHash,
		})
	}
	return registerQueryInfos, nil
}

func (rp registerLockScriptRepo) FindOrCreateScript(ctx context.Context, script *biz.Script) error {
	ht, err := hashType(script.HashType)
	if err != nil {
		return err
	}
	s := Script{}
	if err = rp.data.db.WithContext(ctx).FirstOrCreate(&s, Script{
		CodeHash:    script.CodeHash,
		CodeHashCrc: crc32.ChecksumIEEE([]byte(script.CodeHash)),
		HashType:    ht,
		Args:        script.Args,
		ArgsCrc:     crc32.ChecksumIEEE([]byte(script.Args)),
	}).Error; err != nil {
		return err
	}
	script.ID = s.ID
	return nil
}
