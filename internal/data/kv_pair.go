package data

import (
	"context"
	"fmt"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"gorm.io/gorm"
)

var _ biz.KvPairRepo = (*kvPairRepo)(nil)

type kvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewKvPairRepo(data *Data, logger *logger.Logger) biz.KvPairRepo {
	return &kvPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp kvPairRepo) CreateKvPairs(ctx context.Context, kvPair *biz.KvPair) error {
	return rp.data.db.Transaction(func(tx *gorm.DB) error {
		// insert register cotas
		for _, register := range kvPair.Registers {
			if err := tx.WithContext(ctx).Create(register).Error; err != nil {
				return err
			}
		}
		// insert define cotas
		if err := tx.WithContext(ctx).Create(kvPair.DefineCotas).Error; err != nil {
			return err
		}
		// insert hold cotas
		for _, holdCota := range kvPair.HoldCotas {
			var exists bool
			if err := tx.WithContext(ctx).Model(HoldCotaNftKvPair{}).Where("cota_id_crc = ? and token_index = ? and lock_hash_crc = ? and cota_id = ?and lock_hash = ?", holdCota.CotaIdCRC, holdCota.TokenIndex, holdCota.LockHashCRC, holdCota.CotaId, holdCota.LockHash).Find(&exists).Error; err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("hold cota kv pair already exist cota_id: %v, token_index: %v, lock_hash: %v", holdCota.CotaId, holdCota.TokenIndex, holdCota.LockHash)
			}
			if err := tx.WithContext(ctx).Create(holdCota).Error; err != nil {
				return err
			}
		}

		// insert withdraw cotas
		if err := tx.WithContext(ctx).Create(kvPair.WithdrawCotas).Error; err != nil {
			return err
		}
		// remove those hold cotas that are equal with withdraw cotas
		holdCotasSize := len(kvPair.WithdrawCotas)
		holdCotas := make([]HoldCotaNftKvPair, holdCotasSize)
		for _, withdrawCota := range kvPair.WithdrawCotas {
			var holdCota HoldCotaNftKvPair
			if err := tx.WithContext(ctx).Select("id").Where("cota_id_crc = ? and token_index = ? and lock_hash_crc = ? and cota_id = ? and lock_hash = ?", withdrawCota.CotaIdCRC, withdrawCota.TokenIndex, withdrawCota.LockHashCrc, withdrawCota.CotaId, withdrawCota.LockHash).Find(&holdCota).Error; err != nil {
				return err
			}
			holdCotas = append(holdCotas, holdCota)
		}
		holdCotaIds := make([]int, holdCotasSize)
		if err := tx.WithContext(ctx).Delete(&holdCotas, holdCotaIds).Error; err != nil {
			return err
		}
		// insert claimed cotas
		if err := tx.WithContext(ctx).Create(kvPair.ClaimedCotas).Error; err != nil {
			return err
		}
		return nil
	})
}

func (rp kvPairRepo) DeleteKvPairs(ctx context.Context, blockNumber uint64) error {
	return rp.data.db.Transaction(func(tx *gorm.DB) error {
		// delete all register cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(RegisterCotaKvPair{}).Error; err != nil {
			return err
		}
		// delete all define cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(DefineCotaNftKvPair{}).Error; err != nil {
			return err
		}
		// delete all hold cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(HoldCotaNftKvPair{}).Error; err != nil {
			return err
		}
		// delete all withdraw cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(WithdrawCotaNftKvPair{}).Error; err != nil {
			return err
		}
		// delete all claimed cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(ClaimedCotaNftKvPair{}).Error; err != nil {
			return err
		}
		return nil
	})
}
