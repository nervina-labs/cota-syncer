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

func (rp kvPairRepo) CreateKvPairs(ctx context.Context, checkInfo biz.CheckInfo, kvPair *biz.KvPair) error {
	return rp.data.db.Transaction(func(tx *gorm.DB) error {
		// create register cotas
		if kvPair.HasRegisters() {
			registers := make([]RegisterCotaKvPair, len(kvPair.Registers))
			for i, register := range kvPair.Registers {
				registers[i] = RegisterCotaKvPair{
					BlockNumber: register.BlockNumber,
					LockHash:    register.LockHash,
				}
			}
			if err := tx.Model(RegisterCotaKvPair{}).WithContext(ctx).Create(registers).Error; err != nil {
				return err
			}
		}
		// create define cotas
		if kvPair.HasDefineCotas() {
			defineCotas := make([]DefineCotaNftKvPair, len(kvPair.DefineCotas))
			for i, cota := range kvPair.DefineCotas {
				defineCotas[i] = DefineCotaNftKvPair{
					BlockNumber: cota.BlockNumber,
					CotaId:      cota.CotaId,
					Total:       cota.Total,
					Issued:      cota.Issued,
					Configure:   cota.Configure,
					LockHash:    cota.LockHash,
					LockHashCRC: cota.LockHashCRC,
				}
			}
			if err := tx.Model(DefineCotaNftKvPair{}).WithContext(ctx).Create(defineCotas).Error; err != nil {
				return err
			}
			defineCotaVersions := make([]DefineCotaNftKvPairVersion, len(kvPair.DefineCotas))
			for i, define := range kvPair.DefineCotas {
				defineCotaVersion := DefineCotaNftKvPairVersion{
					BlockNumber: define.BlockNumber,
					CotaId:      define.CotaId,
					Total:       define.Total,
					Issued:      define.Issued,
					OldIssued:   define.Issued,
					Configure:   define.Configure,
					LockHash:    define.LockHash,
					TxIndex:     define.TxIndex,
					ActionType:  0,
				}
				defineCotaVersions[i] = defineCotaVersion
			}
			// create define cotas versions
			if err := tx.Model(DefineCotaNftKvPairVersion{}).WithContext(ctx).Create(defineCotaVersions).Error; err != nil {
				return err
			}
		}
		if kvPair.HasUpdatedDefineCotas() {
			updatedDefineCotaVersions := make([]DefineCotaNftKvPairVersion, len(kvPair.UpdatedDefineCotas))
			for i, define := range kvPair.UpdatedDefineCotas {
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
					TxIndex:        define.TxIndex,
					ActionType:     1,
				}
				updatedDefineCotaVersions[i] = defineCotaVersion
			}
			// create updated define cotas versions
			if err := tx.Model(DefineCotaNftKvPairVersion{}).WithContext(ctx).Create(updatedDefineCotaVersions).Error; err != nil {
				return err
			}
			// update define cotas
			updatedDefineCotas := make([]DefineCotaNftKvPair, len(kvPair.UpdatedDefineCotas))
			for i, cota := range kvPair.UpdatedDefineCotas {
				updatedDefineCotas[i] = DefineCotaNftKvPair{
					BlockNumber: cota.BlockNumber,
					CotaId:      cota.CotaId,
					Total:       cota.Total,
					Issued:      cota.Issued,
					Configure:   cota.Configure,
					LockHash:    cota.LockHash,
					LockHashCRC: cota.LockHashCRC,
				}
			}
			if err := tx.Model(DefineCotaNftKvPair{}).WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cota_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"issued"}),
			}).Create(updatedDefineCotas).Error; err != nil {
				return err
			}
		}
		if kvPair.HasWithdrawCotas() {
			// create withdraw cotas
			withdrawCotas := make([]WithdrawCotaNftKvPair, len(kvPair.WithdrawCotas))
			for i, cota := range kvPair.WithdrawCotas {
				withdrawCotas[i] = WithdrawCotaNftKvPair{
					BlockNumber:          cota.BlockNumber,
					CotaId:               cota.CotaId,
					CotaIdCRC:            cota.CotaIdCRC,
					TokenIndex:           cota.TokenIndex,
					OutPoint:             cota.OutPoint,
					OutPointCrc:          cota.OutPointCrc,
					State:                cota.State,
					Configure:            cota.Configure,
					Characteristic:       cota.Characteristic,
					ReceiverLockScriptId: cota.ReceiverLockScriptId,
					LockHash:             cota.LockHash,
					LockHashCrc:          cota.LockHashCrc,
				}
			}
			if err := tx.Model(WithdrawCotaNftKvPair{}).WithContext(ctx).Create(withdrawCotas).Error; err != nil {
				return err
			}
			holdCotasSize := len(kvPair.WithdrawCotas)
			removedHoldCotas := make([]biz.HoldCotaNftKvPair, holdCotasSize)
			removedHoldCotaIds := make([]uint, holdCotasSize)
			for i, withdrawCota := range kvPair.WithdrawCotas {
				var holdCota biz.HoldCotaNftKvPair
				if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Select("*").Where("cota_id = ? and token_index = ?", withdrawCota.CotaId, withdrawCota.TokenIndex).Find(&holdCota).Error; err != nil {
					return err
				}
				// 上面把对象初始化出来了，所以需要通过具体值来判断是否存在
				if holdCota.CotaId == "" {
					continue
				}
				removedHoldCotas[i] = holdCota
				removedHoldCotaIds[i] = holdCota.ID
			}
			if removedHoldCotas[0].CotaId != "" {
				removedHoldCotaVersions := make([]HoldCotaNftKvPairVersion, holdCotasSize)
				blockNumber := kvPair.WithdrawCotas[0].BlockNumber
				for i, cota := range removedHoldCotas {
					removedHoldCotaVersions[i] = HoldCotaNftKvPairVersion{
						OldBlockNumber:    cota.BlockNumber,
						BlockNumber:       blockNumber,
						CotaId:            cota.CotaId,
						TokenIndex:        cota.TokenIndex,
						OldState:          cota.State,
						Configure:         cota.Configure,
						OldCharacteristic: cota.Characteristic,
						OldLockHash:       cota.LockHash,
						TxIndex:           cota.TxIndex,
						ActionType:        2,
					}
				}
				// create removed hold cota versions
				if err := tx.Model(HoldCotaNftKvPairVersion{}).WithContext(ctx).Create(removedHoldCotaVersions).Error; err != nil {
					return err
				}
				// remove those hold cotas that are equal with withdraw cotas
				if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Delete(&removedHoldCotas, removedHoldCotaIds).Error; err != nil {
					return err
				}
			}
		}
		if kvPair.HasHoldCotas() {
			// create hold cotas
			holdCotas := make([]HoldCotaNftKvPair, len(kvPair.HoldCotas))
			for i, cota := range kvPair.HoldCotas {
				holdCotas[i] = HoldCotaNftKvPair{
					BlockNumber:    cota.BlockNumber,
					CotaId:         cota.CotaId,
					TokenIndex:     cota.TokenIndex,
					State:          cota.State,
					Configure:      cota.Configure,
					Characteristic: cota.Characteristic,
					LockHash:       cota.LockHash,
					LockHashCRC:    cota.LockHashCRC,
				}
			}
			if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Create(holdCotas).Error; err != nil {
				return err
			}
			newHoldCotaVersions := make([]HoldCotaNftKvPairVersion, len(kvPair.HoldCotas))
			for i, cota := range kvPair.HoldCotas {
				newHoldCotaVersions[i] = HoldCotaNftKvPairVersion{
					BlockNumber:    cota.BlockNumber,
					CotaId:         cota.CotaId,
					TokenIndex:     cota.TokenIndex,
					State:          cota.State,
					Configure:      cota.Configure,
					Characteristic: cota.Characteristic,
					LockHash:       cota.LockHash,
					TxIndex:        cota.TxIndex,
					ActionType:     0,
				}
			}
			// create hold cota versions
			if err := tx.Model(HoldCotaNftKvPairVersion{}).WithContext(ctx).Create(newHoldCotaVersions).Error; err != nil {
				return err
			}
		}
		if kvPair.HasUpdatedHoldCotas() {
			updatedHoldCotaVersions := make([]HoldCotaNftKvPairVersion, len(kvPair.UpdatedHoldCotas))
			for i, cota := range kvPair.UpdatedHoldCotas {
				var oldHoldCota HoldCotaNftKvPair
				if err := tx.Model(HoldCotaNftKvPair{}).WithContext(ctx).Where("cota_id = ? and token_index = ?", cota.CotaId, cota.TokenIndex).First(&oldHoldCota).Error; err != nil {
					return err
				}
				updatedHoldCotaVersions[i] = HoldCotaNftKvPairVersion{
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
					TxIndex:           cota.TxIndex,
					ActionType:        1,
				}
			}
			// create updated hold cotas versions
			if err := tx.Model(HoldCotaNftKvPairVersion{}).WithContext(ctx).Create(updatedHoldCotaVersions).Error; err != nil {
				return err
			}
			// update hold cotas
			updatedHoldCotas := make([]HoldCotaNftKvPair, len(kvPair.UpdatedHoldCotas))
			for i, cota := range kvPair.UpdatedHoldCotas {
				updatedHoldCotas[i] = HoldCotaNftKvPair{
					BlockNumber:    cota.BlockNumber,
					CotaId:         cota.CotaId,
					TokenIndex:     cota.TokenIndex,
					State:          cota.State,
					Configure:      cota.Configure,
					Characteristic: cota.Characteristic,
					LockHash:       cota.LockHash,
					LockHashCRC:    cota.LockHashCRC,
				}
			}
			if err := tx.Debug().Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cota_id"}, {Name: "token_index"}},
				DoUpdates: clause.AssignmentColumns([]string{"block_number", "state", "characteristic", "lock_hash", "lock_hash_crc"}),
			}).Create(updatedHoldCotas).Error; err != nil {
				return err
			}
		}
		if kvPair.HasClaimedCotas() {
			// create claimed cotas
			claimedCotas := make([]ClaimedCotaNftKvPair, len(kvPair.ClaimedCotas))
			for i, cota := range kvPair.ClaimedCotas {
				claimedCotas[i] = ClaimedCotaNftKvPair{
					BlockNumber: cota.BlockNumber,
					CotaId:      cota.CotaId,
					CotaIdCRC:   cota.CotaIdCRC,
					TokenIndex:  cota.TokenIndex,
					OutPoint:    cota.OutPoint,
					OutPointCrc: cota.OutPointCrc,
					LockHash:    cota.LockHash,
					LockHashCrc: cota.LockHashCrc,
				}
			}
			if err := tx.Model(ClaimedCotaNftKvPair{}).WithContext(ctx).Create(claimedCotas).Error; err != nil {
				return err
			}
		}
		// create check info
		if err := tx.Debug().Model(CheckInfo{}).WithContext(ctx).Create(&CheckInfo{
			BlockNumber: checkInfo.BlockNumber,
			BlockHash:   checkInfo.BlockHash,
			CheckType:   checkInfo.CheckType,
		}).Error; err != nil {
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
		if len(updatedDefineCotas) > 0 {
			if err := tx.Model(DefineCotaNftKvPair{}).WithContext(ctx).Create(updatedDefineCotas).Error; err != nil {
				return err
			}
		}
		// delete all updated define versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Delete(DefineCotaNftKvPairVersion{}).Error; err != nil {
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
		if len(deletedHoldCotas) > 0 {
			if err := tx.WithContext(ctx).Create(deletedHoldCotas).Error; err != nil {
				return err
			}
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
		if len(updatedHoldCotaVersions) > 0 {
			if err := tx.WithContext(ctx).Create(updatedHoldCotas).Error; err != nil {
				return err
			}
		}
		// delete all updated hold cota versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Delete(HoldCotaNftKvPairVersion{}).Error; err != nil {
			return err
		}
		// delete all claimed cotas by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(ClaimedCotaNftKvPair{}).Error; err != nil {
			return err
		}

		// delete check info
		if err := tx.Debug().WithContext(ctx).Where("block_number = ? and check_type = ?", blockNumber, 0).Delete(CheckInfo{}).Error; err != nil {
			return err
		}

		return nil
	})
}
