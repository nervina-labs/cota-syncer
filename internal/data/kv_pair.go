package data

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hash/crc32"
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

func (rp kvPairRepo) CreateKvPairs(ctx context.Context, txIndex int, kvPair *biz.KvPair) error {
	return rp.data.db.Transaction(func(tx *gorm.DB) error {
		// create register cotas
		for _, register := range kvPair.Registers {
			if err := tx.Model(RegisterCotaKvPair{}).WithContext(ctx).Create(register).Error; err != nil {
				return err
			}
		}
		// create define cotas
		if err := tx.Model(DefineCotaNftKvPair{}).WithContext(ctx).Create(kvPair.DefineCotas).Error; err != nil {
			return err
		}
		defineCotaVersions := make([]DefineCotaNftKvPairVersion, len(kvPair.DefineCotas))
		for _, define := range kvPair.DefineCotas {
			defineCotaVersion := DefineCotaNftKvPairVersion{
				BlockNumber: define.BlockNumber,
				CotaId:      define.CotaId,
				Total:       define.Total,
				Issued:      define.Issued,
				OldIssued:   define.Issued,
				Configure:   define.Configure,
				LockHash:    define.LockHash,
				TxIndex:     uint32(txIndex),
				ActionType:  0,
			}
			defineCotaVersions = append(defineCotaVersions, defineCotaVersion)
		}
		// create define cotas versions
		if err := tx.Model(DefineCotaNftKvPairVersion{}).WithContext(ctx).Create(defineCotaVersions).Error; err != nil {
			return err
		}

		updatedDefineCotaVersions := make([]DefineCotaNftKvPairVersion, len(kvPair.UpdatedDefineCotas))
		for _, define := range kvPair.UpdatedDefineCotas {
			var defineCota DefineCotaNftKvPair
			if err := tx.Model(DefineCotaNftKvPair{}).WithContext(ctx).Where("cota_id = ?", define.CotaId).First(&defineCota).Error; err != nil {
				return err
			}
			defineCotaVersion := DefineCotaNftKvPairVersion{
				OldBlockNumber: defineCota.BlockNumber,
				BlockNumber:    define.BlockNumber,
				CotaId:         define.CotaId,
				Total:          define.Total,
				Issued:         define.Issued,
				OldIssued:      defineCota.Issued,
				Configure:      define.Configure,
				LockHash:       define.LockHash,
				TxIndex:        uint32(txIndex),
				ActionType:     1,
			}
			updatedDefineCotaVersions = append(updatedDefineCotaVersions, defineCotaVersion)
		}

		// create updated define cotas versions
		if err := tx.Model(DefineCotaNftKvPairVersion{}).WithContext(ctx).Create(updatedDefineCotaVersions).Error; err != nil {
			return err
		}
		// update define cotas
		if err := tx.Model(DefineCotaNftKvPair{}).WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "cota_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"issued"}),
		}).Create(kvPair.UpdatedDefineCotas).Error; err != nil {
			return err
		}
		// create withdraw cotas
		if err := tx.Model(WithdrawCotaNftKvPair{}).WithContext(ctx).Create(kvPair.WithdrawCotas).Error; err != nil {
			return err
		}

		holdCotasSize := len(kvPair.WithdrawCotas)
		removedHoldCotas := make([]HoldCotaNftKvPair, holdCotasSize)
		removedHoldCotaIds := make([]uint, holdCotasSize)
		for _, withdrawCota := range kvPair.WithdrawCotas {
			var holdCota HoldCotaNftKvPair
			if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Select("id").Where("cota_id = ? and token_index = ?", withdrawCota.CotaId, withdrawCota.TokenIndex).Find(&holdCota).Error; err != nil {
				return err
			}
			removedHoldCotas = append(removedHoldCotas, holdCota)
			removedHoldCotaIds = append(removedHoldCotaIds, holdCota.ID)
		}
		removedHoldCotaVersions := make([]HoldCotaNftKvPairVersion, holdCotasSize)
		blockNumber := kvPair.WithdrawCotas[0].BlockNumber
		for _, cota := range removedHoldCotas {
			removedHoldCotaVersions = append(removedHoldCotaVersions, HoldCotaNftKvPairVersion{
				OldBlockNumber:    cota.BlockNumber,
				BlockNumber:       blockNumber,
				CotaId:            cota.CotaId,
				TokenIndex:        cota.TokenIndex,
				OldState:          cota.State,
				Configure:         cota.Configure,
				OldCharacteristic: cota.Characteristic,
				OldLockHash:       cota.LockHash,
				TxIndex:           uint32(txIndex),
				ActionType:        2,
			})
		}
		// create removed hold cota versions
		if err := tx.Model(HoldCotaNftKvPairVersion{}).WithContext(ctx).Create(removedHoldCotaVersions).Error; err != nil {
			return err
		}
		// remove those hold cotas that are equal with withdraw cotas
		if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Unscoped().Delete(&removedHoldCotas, removedHoldCotaIds).Error; err != nil {
			return err
		}
		// create hold cotas
		if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Create(kvPair.HoldCotas).Error; err != nil {
			return err
		}
		newHoldCotas := make([]HoldCotaNftKvPairVersion, len(kvPair.HoldCotas))
		for _, cota := range kvPair.HoldCotas {
			newHoldCotas = append(newHoldCotas, HoldCotaNftKvPairVersion{
				BlockNumber:    cota.BlockNumber,
				CotaId:         cota.CotaId,
				TokenIndex:     cota.TokenIndex,
				State:          cota.State,
				Configure:      cota.Configure,
				Characteristic: cota.Characteristic,
				LockHash:       cota.LockHash,
				TxIndex:        uint32(txIndex),
				ActionType:     0,
			})
		}
		// create hold cota versions
		if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Create(newHoldCotas).Error; err != nil {
			return err
		}

		updatedHoldCotaVersions := make([]HoldCotaNftKvPairVersion, len(kvPair.UpdatedHoldCotas))
		for _, cota := range kvPair.UpdatedHoldCotas {
			var oldHoldCota HoldCotaNftKvPair
			if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Where("cota_id = ? and token_index = ?", cota.CotaId, cota.TokenIndex).First(&oldHoldCota).Error; err != nil {
				return err
			}
			updatedHoldCotaVersions = append(updatedHoldCotaVersions, HoldCotaNftKvPairVersion{
				OldBlockNumber:    oldHoldCota.BlockNumber,
				BlockNumber:       cota.BlockNumber,
				CotaId:            cota.CotaId,
				TokenIndex:        cota.TokenIndex,
				OldState:          oldHoldCota.State,
				State:             cota.State,
				Configure:         cota.Configure,
				OldCharacteristic: oldHoldCota.Characteristic,
				Characteristic:    cota.Characteristic,
				OldLockHash:       oldHoldCota.LockHash,
				LockHash:          cota.LockHash,
				TxIndex:           uint32(txIndex),
				ActionType:        1,
			})
		}
		// create updated hold cotas versions
		if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Create(updatedHoldCotaVersions).Error; err != nil {
			return err
		}
		// update hold cotas
		if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Select("state", "characteristic", "block_number").Updates(kvPair.UpdatedHoldCotas).Error; err != nil {
			return err
		}
		// create claimed cotas
		if err := tx.Model(ClaimedCotaNftKvPair{}).WithContext(ctx).Create(kvPair.ClaimedCotas).Error; err != nil {
			return err
		}
		return nil
	})
}

func (rp kvPairRepo) RestoreKvPairs(ctx context.Context, blockNumber uint64) error {
	return rp.data.db.Transaction(func(tx *gorm.DB) error {
		// delete all register cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(RegisterCotaKvPair{}).Error; err != nil {
			return err
		}
		// delete all new define cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(DefineCotaNftKvPair{}).Error; err != nil {
			return err
		}
		// delete all create define cota versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 0).Delete(DefineCotaNftKvPairVersion{}).Error; err != nil {
			return err
		}
		// 把需要回滚的 block 更新过的 define 恢复到更新前的状态
		var updatedDefineCotaVersions []DefineCotaNftKvPairVersion
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Group("cota_id").Order("tx_index").Find(&updatedDefineCotaVersions).Error; err != nil {
			return err
		}
		var updatedDefineCotas []DefineCotaNftKvPair
		for _, version := range updatedDefineCotaVersions {
			updatedDefineCotas = append(updatedDefineCotas, DefineCotaNftKvPair{
				BlockNumber: version.OldBlockNumber,
				CotaId:      version.CotaId,
				Total:       version.Total,
				Issued:      version.OldIssued,
				Configure:   version.Configure,
				LockHash:    version.LockHash,
				LockHashCRC: crc32.ChecksumIEEE([]byte(version.LockHash)),
			})
		}
		if err := tx.Model(DefineCotaNftKvPair{}).WithContext(ctx).Create(updatedDefineCotas).Error; err != nil {
			return err
		}
		// delete all updated define versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ？and action_type = ?", blockNumber, 1).Delete(DefineCotaNftKvPairVersion{}).Error; err != nil {
			return err
		}
		// delete all withdraw cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(WithdrawCotaNftKvPair{}).Error; err != nil {
			return err
		}
		// delete all hold cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(HoldCotaNftKvPair{}).Error; err != nil {
			return err
		}
		// delete all created hold cota versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 0).Delete(HoldCotaNftKvPairVersion{}).Error; err != nil {
			return err
		}
		// restore all deleted hold cotas by the block number
		var deletedHoldCotaVersions []HoldCotaNftKvPairVersion
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 2).Group("cota_id, token_index").Order("tx_index").Find(&deletedHoldCotaVersions).Error; err != nil {
			return err
		}
		var deletedHoldCotas []HoldCotaNftKvPair
		for _, version := range deletedHoldCotaVersions {
			deletedHoldCotas = append(deletedHoldCotas, HoldCotaNftKvPair{
				BlockNumber:    version.OldBlockNumber,
				CotaId:         version.CotaId,
				TokenIndex:     version.TokenIndex,
				State:          version.OldState,
				Configure:      version.Configure,
				Characteristic: version.OldCharacteristic,
				LockHash:       version.OldLockHash,
				LockHashCRC:    crc32.ChecksumIEEE([]byte(version.OldLockHash)),
			})
		}
		if err := tx.WithContext(ctx).Create(deletedHoldCotas).Error; err != nil {
			return err
		}
		// delete all deleted hold cota versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 2).Delete(HoldCotaNftKvPairVersion{}).Error; err != nil {
			return err
		}
		// restore all updated hold cotas by the block number
		var updatedHoldCotaVersions []HoldCotaNftKvPairVersion
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Group("cota_id, token_index").Order("tx_index").Find(&deletedHoldCotaVersions).Error; err != nil {
			return err
		}
		var updatedHoldCotas []HoldCotaNftKvPair
		for _, version := range updatedHoldCotaVersions {
			updatedHoldCotas = append(updatedHoldCotas, HoldCotaNftKvPair{
				BlockNumber:    version.OldBlockNumber,
				CotaId:         version.CotaId,
				TokenIndex:     version.TokenIndex,
				State:          version.OldState,
				Configure:      version.Configure,
				Characteristic: version.OldCharacteristic,
				LockHash:       version.OldLockHash,
				LockHashCRC:    crc32.ChecksumIEEE([]byte(version.OldLockHash)),
			})
		}
		if err := tx.WithContext(ctx).Create(updatedHoldCotas).Error; err != nil {
			return err
		}
		// delete all updated hold cota versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Delete(HoldCotaNftKvPair{}).Error; err != nil {
			return err
		}
		// delete all claimed cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(ClaimedCotaNftKvPair{}).Error; err != nil {
			return err
		}
		return nil
	})
}
