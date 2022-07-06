package data

import (
	"context"
	"hash/crc32"

	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
)

const pageSize = 1000

// var _ biz.WithdrawExtraInfoRepo = (*withdrawExtraInfoRepo)(nil)

type withdrawExtraInfoRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewWithdrawExtraInfoRepo(data *Data, logger *logger.Logger) biz.WithdrawExtraInfoRepo {
	return &withdrawExtraInfoRepo{
		data:   data,
		logger: logger,
	}
}

func (rp withdrawExtraInfoRepo) CreateExtraInfo(ctx context.Context, outPoint string, txHash string, lockScriptId uint) error {
	var withdraw WithdrawCotaNftKvPair
	if err := rp.data.db.WithContext(ctx).Where("out_point = ?", outPoint).First(withdraw).Error; err == nil {
		withdraw.TxHash = txHash
		withdraw.LockScriptId = lockScriptId
		if err = rp.data.db.WithContext(ctx).Save(&withdraw).Error; err != nil {
			return err
		}
	}
	return nil
}

func (rp withdrawExtraInfoRepo) FindAllQueryInfos(ctx context.Context) ([]biz.WithdrawQueryInfo, error) {
	var (
		withdrawals []WithdrawCotaNftKvPair
		queryInfos  []biz.WithdrawQueryInfo
	)
	offset := 0
	for {
		result := rp.data.db.WithContext(ctx).Select("out_point").Where("tx_hash IS NULL").Limit(pageSize).Offset(offset * pageSize).Find(&withdrawals)
		if result.Error != nil {
			return queryInfos, result.Error
		}
		if result.RowsAffected == 0 {
			break
		}
		offset++
	}
	for _, v := range withdrawals {
		queryInfos = append(queryInfos, biz.WithdrawQueryInfo{
			BlockNumber: v.BlockNumber,
			OutPoint:    v.OutPoint,
		})
	}
	return queryInfos, nil
}

func (rp withdrawExtraInfoRepo) FindOrCreateScript(ctx context.Context, script *biz.Script) error {
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
