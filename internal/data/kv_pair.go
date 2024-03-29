package data

import (
	"context"
	"errors"
	"hash/crc32"
	"time"

	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (rp kvPairRepo) CreateCotaEntryKvPairs(ctx context.Context, checkInfo biz.CheckInfo, kvPair *biz.KvPair) error {
	return rp.data.db.Transaction(func(tx *gorm.DB) error {
		// create register cotas
		if kvPair.HasRegisters() {
			registers := make([]RegisterCotaKvPair, len(kvPair.Registers))
			for i, register := range kvPair.Registers {
				registers[i] = RegisterCotaKvPair{
					BlockNumber:  register.BlockNumber,
					LockHash:     register.LockHash,
					CotaCellID:   register.CotaCellID,
					LockScriptId: register.LockScriptId,
				}
			}
			if err := tx.Model(RegisterCotaKvPair{}).WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "lock_hash"}},
				DoUpdates: clause.AssignmentColumns([]string{"cota_cell_id", "lock_script_id", "updated_at"}),
			}).Create(registers).Error; err != nil {
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
			if err := tx.Debug().Model(DefineCotaNftKvPair{}).WithContext(ctx).Create(defineCotas).Error; err != nil {
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
					UpdatedAt:   cota.UpdatedAt,
				}
			}
			if err := tx.Model(DefineCotaNftKvPair{}).WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cota_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"issued", "block_number", "updated_at"}),
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
					TxHash:               cota.TxHash,
					State:                cota.State,
					Configure:            cota.Configure,
					Characteristic:       cota.Characteristic,
					ReceiverLockScriptId: cota.ReceiverLockScriptId,
					LockHash:             cota.LockHash,
					LockHashCrc:          cota.LockHashCrc,
					LockScriptId:         cota.LockScriptId,
					Version:              cota.Version,
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
					UpdatedAt:      cota.UpdatedAt,
				}
			}
			if err := tx.Debug().Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cota_id"}, {Name: "token_index"}},
				DoUpdates: clause.AssignmentColumns([]string{"block_number", "state", "characteristic", "lock_hash", "lock_hash_crc", "updated_at"}),
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

		if kvPair.HasExtensionPairs() {
			// create extension pairs
			extensionPairs := make([]ExtensionKvPair, len(kvPair.ExtensionPairs))
			for i, extension := range kvPair.ExtensionPairs {
				extensionPairs[i] = ExtensionKvPair{
					BlockNumber: extension.BlockNumber,
					Key:         extension.Key,
					Value:       extension.Value,
					LockHash:    extension.LockHash,
					LockHashCRC: extension.LockHashCRC,
				}
			}
			if err := tx.Debug().Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "key"}, {Name: "lock_hash"}},
				DoUpdates: clause.AssignmentColumns([]string{"block_number", "value", "updated_at"}),
			}).Create(extensionPairs).Error; err != nil {
				return err
			}
			extensionPairVersions := make([]ExtensionKvPairVersion, len(kvPair.ExtensionPairs))
			for i, extension := range kvPair.ExtensionPairs {
				extensionPairVersions[i] = ExtensionKvPairVersion{
					BlockNumber: extension.BlockNumber,
					Key:         extension.Key,
					Value:       extension.Value,
					LockHash:    extension.LockHash,
					TxIndex:     extension.TxIndex,
					ActionType:  0,
				}
			}
			// create extension pair versions
			if err := tx.Model(ExtensionKvPairVersion{}).WithContext(ctx).Create(extensionPairVersions).Error; err != nil {
				return err
			}
		}
		if kvPair.HasUpdatedExtensionPairs() {
			updatedExtensionPairVersions := make([]ExtensionKvPairVersion, len(kvPair.UpdatedExtensionPairs))
			for i, extension := range kvPair.UpdatedExtensionPairs {
				var oldExtension ExtensionKvPair
				if err := tx.Model(ExtensionKvPair{}).WithContext(ctx).Where("`key` = ?", extension.Key).First(&oldExtension).Error; err != nil {
					return err
				}
				updatedExtensionPairVersions[i] = ExtensionKvPairVersion{
					OldBlockNumber: oldExtension.BlockNumber,
					BlockNumber:    extension.BlockNumber,
					Key:            extension.Key,
					Value:          extension.Value,
					OldValue:       oldExtension.Value,
					LockHash:       extension.LockHash,
					TxIndex:        extension.TxIndex,
					ActionType:     1,
				}
			}
			// create updated extension pair versions
			if err := tx.Model(ExtensionKvPairVersion{}).WithContext(ctx).Create(updatedExtensionPairVersions).Error; err != nil {
				return err
			}

			// update extension pairs
			updatedExtensionPairs := make([]ExtensionKvPair, len(kvPair.UpdatedExtensionPairs))
			for i, extension := range kvPair.UpdatedExtensionPairs {
				updatedExtensionPairs[i] = ExtensionKvPair{
					BlockNumber: extension.BlockNumber,
					Key:         extension.Key,
					Value:       extension.Value,
					LockHash:    extension.LockHash,
					LockHashCRC: extension.LockHashCRC,
					UpdatedAt:   extension.UpdatedAt,
				}
			}
			if err := tx.Debug().Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "key"}, {Name: "lock_hash"}},
				DoUpdates: clause.AssignmentColumns([]string{"block_number", "value", "updated_at"}),
			}).Create(updatedExtensionPairs).Error; err != nil {
				return err
			}
		}

		if kvPair.HasSubKeyPairs() {
			subKeyPairs := make([]SubKeyKvPair, len(kvPair.SubKeyPairs))
			subKeyPairVersions := make([]SubKeyKvPairVersion, len(kvPair.SubKeyPairs))
			for i, subKey := range kvPair.SubKeyPairs {
				subKeyPairs[i] = SubKeyKvPair{
					BlockNumber: subKey.BlockNumber,
					LockHash:    subKey.LockHash,
					SubType:     subKey.SubType,
					ExtData:     subKey.ExtData,
					AlgIndex:    subKey.AlgIndex,
					PubkeyHash:  subKey.PubkeyHash,
				}
				subKeyPairVersions[i] = SubKeyKvPairVersion{
					BlockNumber: subKey.BlockNumber,
					LockHash:    subKey.LockHash,
					SubType:     subKey.SubType,
					ExtData:     subKey.ExtData,
					AlgIndex:    subKey.AlgIndex,
					PubkeyHash:  subKey.PubkeyHash,
					ActionType:  0,
				}
			}
			if err := tx.Debug().WithContext(ctx).Create(subKeyPairs).Error; err != nil {
				return err
			}
			if err := tx.Debug().WithContext(ctx).Create(subKeyPairVersions).Error; err != nil {
				return err
			}
		}
		if kvPair.HasUpdatedSubKeyPairs() {
			updatedSubKeyPairVersions := make([]SubKeyKvPairVersion, len(kvPair.UpdatedSubKeyPairs))
			updatedSubKeyPairs := make([]SubKeyKvPair, len(kvPair.UpdatedSubKeyPairs))
			for i, subKey := range kvPair.UpdatedSubKeyPairs {
				var oldSubKey SubKeyKvPair
				if err := tx.Model(SubKeyKvPair{}).WithContext(ctx).Where("lock_hash = ? and ext_data = ?", subKey.LockHash, subKey.ExtData).First(&oldSubKey).Error; err != nil {
					return err
				}
				updatedSubKeyPairVersions[i] = SubKeyKvPairVersion{
					OldBlockNumber: oldSubKey.BlockNumber,
					BlockNumber:    subKey.BlockNumber,
					LockHash:       subKey.LockHash,
					SubType:        subKey.SubType,
					ExtData:        subKey.ExtData,
					OldAlgIndex:    oldSubKey.AlgIndex,
					AlgIndex:       subKey.AlgIndex,
					OldPubkeyHash:  oldSubKey.PubkeyHash,
					PubkeyHash:     subKey.PubkeyHash,
					ActionType:     1,
				}
				updatedSubKeyPairs[i] = SubKeyKvPair{
					BlockNumber: subKey.BlockNumber,
					LockHash:    subKey.LockHash,
					SubType:     subKey.SubType,
					ExtData:     subKey.ExtData,
					AlgIndex:    subKey.AlgIndex,
					PubkeyHash:  subKey.PubkeyHash,
					UpdatedAt:   subKey.UpdatedAt,
				}
			}
			if err := tx.Debug().WithContext(ctx).Create(updatedSubKeyPairVersions).Error; err != nil {
				return err
			}
			if err := tx.Debug().Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "lock_hash"}, {Name: "ext_data"}},
				DoUpdates: clause.AssignmentColumns([]string{"block_number", "alg_index", "pubkey_hash", "updated_at"}),
			}).Create(updatedSubKeyPairs).Error; err != nil {
				return err
			}
		}

		if kvPair.HasSocialPairs() {
			socialPairs := make([]SocialKvPair, len(kvPair.SocialPairs))
			socialPairVersions := make([]SocialKvPairVersion, len(kvPair.SocialPairs))
			for i, social := range kvPair.SocialPairs {
				socialPairs[i] = SocialKvPair{
					BlockNumber:  social.BlockNumber,
					LockHash:     social.LockHash,
					LockHashCRC:  social.LockHashCRC,
					RecoveryMode: social.RecoveryMode,
					Must:         social.Must,
					Total:        social.Total,
					Signers:      social.Signers,
				}
				socialPairVersions[i] = SocialKvPairVersion{
					BlockNumber:  social.BlockNumber,
					LockHash:     social.LockHash,
					RecoveryMode: social.RecoveryMode,
					Must:         social.Must,
					Total:        social.Total,
					Signers:      social.Signers,
					ActionType:   0,
				}
			}
			if err := tx.Debug().WithContext(ctx).Create(socialPairs).Error; err != nil {
				return err
			}
			if err := tx.Debug().WithContext(ctx).Create(socialPairVersions).Error; err != nil {
				return err
			}
		}
		if kvPair.HasUpdatedSocialPairs() {
			updatedSocialPairs := make([]SocialKvPair, len(kvPair.UpdatedSocialPairs))
			updatedSocialPairVersions := make([]SocialKvPairVersion, len(kvPair.UpdatedSocialPairs))
			for i, social := range kvPair.UpdatedSocialPairs {
				var oldSocial SocialKvPair
				if err := tx.Model(SocialKvPair{}).WithContext(ctx).Where("lock_hash = ?", social.LockHash).First(&oldSocial).Error; err != nil {
					return err
				}
				updatedSocialPairs[i] = SocialKvPair{
					BlockNumber:  social.BlockNumber,
					LockHash:     social.LockHash,
					LockHashCRC:  social.LockHashCRC,
					RecoveryMode: social.RecoveryMode,
					Must:         social.Must,
					Total:        social.Total,
					Signers:      social.Signers,
				}
				updatedSocialPairVersions[i] = SocialKvPairVersion{
					OldBlockNumber:  oldSocial.BlockNumber,
					BlockNumber:     social.BlockNumber,
					LockHash:        social.LockHash,
					OldRecoveryMode: oldSocial.RecoveryMode,
					RecoveryMode:    social.RecoveryMode,
					OldMust:         oldSocial.Must,
					Must:            social.Must,
					OldTotal:        oldSocial.Total,
					Total:           social.Total,
					OldSigners:      oldSocial.Signers,
					Signers:         social.Signers,
					ActionType:      1,
				}
				if err := tx.Debug().WithContext(ctx).Create(updatedSocialPairVersions).Error; err != nil {
					return err
				}

				if err := tx.Debug().Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "lock_hash"}},
					DoUpdates: clause.AssignmentColumns([]string{"block_number", "recovery_mode", "must", "total", "signers"}),
				}).Create(updatedSocialPairs).Error; err != nil {
					return err
				}
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

func (rp kvPairRepo) RestoreCotaEntryKvPairs(ctx context.Context, blockNumber uint64) error {
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
				UpdatedAt:   time.Now().UTC(),
			})
		}
		if len(updatedDefineCotas) > 0 {
			if err := tx.Debug().WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cota_id"}},
				UpdateAll: true,
			}).Create(updatedDefineCotas).Error; err != nil {
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
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Group("cota_id, token_index").Order("tx_index").Find(&updatedHoldCotaVersions).Error; err != nil {
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

		// delete all extension pairs by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(ExtensionKvPair{}).Error; err != nil {
			return err
		}
		// delete all created extension pair versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 0).Delete(ExtensionKvPairVersion{}).Error; err != nil {
			return err
		}
		// restore all deleted extension pairs by the block number
		var deletedExtensionPairVersions []ExtensionKvPairVersion
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 2).Group("`key`").Order("tx_index").Find(&deletedExtensionPairVersions).Error; err != nil {
			return err
		}
		var deletedExtensionPairs []ExtensionKvPair
		for _, version := range deletedExtensionPairVersions {
			deletedExtensionPairs = append(deletedExtensionPairs, ExtensionKvPair{
				BlockNumber: version.OldBlockNumber,
				Key:         version.Key,
				Value:       version.OldValue,
				LockHash:    version.LockHash,
				LockHashCRC: crc32.ChecksumIEEE([]byte(version.LockHash)),
			})
		}
		if len(deletedExtensionPairs) > 0 {
			if err := tx.WithContext(ctx).Create(deletedExtensionPairs).Error; err != nil {
				return err
			}
		}
		// delete all deleted extension pair versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 2).Delete(ExtensionKvPairVersion{}).Error; err != nil {
			return err
		}
		// restore all updated extension pairs by the block number
		var updatedExtensionPairVersions []ExtensionKvPairVersion
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Group("`key`").Order("tx_index").Find(&updatedExtensionPairVersions).Error; err != nil {
			return err
		}
		var updatedExtensionPairs []ExtensionKvPair
		for _, version := range updatedExtensionPairVersions {
			updatedExtensionPairs = append(updatedExtensionPairs, ExtensionKvPair{
				BlockNumber: version.OldBlockNumber,
				Key:         version.Key,
				Value:       version.OldValue,
				LockHash:    version.LockHash,
				LockHashCRC: crc32.ChecksumIEEE([]byte(version.LockHash)),
			})
		}
		if len(updatedExtensionPairVersions) > 0 {
			if err := tx.WithContext(ctx).Create(updatedExtensionPairs).Error; err != nil {
				return err
			}
		}
		// delete all updated extension pair versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Delete(ExtensionKvPairVersion{}).Error; err != nil {
			return err
		}

		// delete all sub key pair by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(SubKeyKvPair{}).Error; err != nil {
			return err
		}
		// delete all created extension pair versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 0).Delete(SubKeyKvPairVersion{}).Error; err != nil {
			return err
		}
		// restore all updated sub key pairs by the block number
		var updatedSubKeyPairKvVersions []SubKeyKvPairVersion
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Find(&updatedSubKeyPairKvVersions).Error; err != nil {
			return err
		}
		var updatedSubKeyKvPairs []SubKeyKvPair
		for _, version := range updatedSubKeyPairKvVersions {
			updatedSubKeyKvPairs = append(updatedSubKeyKvPairs, SubKeyKvPair{
				BlockNumber: version.OldBlockNumber,
				LockHash:    version.LockHash,
				SubType:     version.SubType,
				ExtData:     version.ExtData,
				AlgIndex:    version.OldAlgIndex,
				PubkeyHash:  version.OldPubkeyHash,
			})
		}
		if len(updatedSubKeyKvPairs) > 0 {
			if err := tx.WithContext(ctx).Create(updatedSubKeyKvPairs).Error; err != nil {
				return err
			}
		}
		// delete all updated sub key pair versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Delete(SubKeyKvPairVersion{}).Error; err != nil {
			return err
		}

		// delete all social key pair by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(SocialKvPair{}).Error; err != nil {
			return err
		}
		// delete all created social key pair versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 0).Delete(SocialKvPairVersion{}).Error; err != nil {
			return err
		}
		// restore all updated sub key pairs by the block number
		var updatedSocialPairVersions []SocialKvPairVersion
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Find(&updatedSocialPairVersions).Error; err != nil {
			return err
		}
		var updatedSocialKvPairs []SocialKvPair
		for _, version := range updatedSocialPairVersions {
			updatedSocialKvPairs = append(updatedSocialKvPairs, SocialKvPair{
				BlockNumber:  version.OldBlockNumber,
				LockHash:     version.LockHash,
				LockHashCRC:  crc32.ChecksumIEEE([]byte(version.LockHash)),
				RecoveryMode: version.OldRecoveryMode,
				Must:         version.OldMust,
				Total:        version.OldTotal,
				Signers:      version.OldSigners,
			})
		}
		if len(updatedSocialKvPairs) > 0 {
			if err := tx.WithContext(ctx).Create(updatedSocialKvPairs).Error; err != nil {
				return err
			}
		}
		// delete all updated sub key pair versions by the block number
		if err := tx.WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Delete(SocialKvPairVersion{}).Error; err != nil {
			return err
		}

		// delete check info
		if err := tx.Debug().WithContext(ctx).Where("block_number = ? and check_type = ?", blockNumber, biz.SyncBlock).Delete(CheckInfo{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (rp kvPairRepo) CreateMetadataKvPairs(ctx context.Context, checkInfo biz.CheckInfo, kvPair *biz.KvPair) error {
	return rp.data.db.Transaction(func(tx *gorm.DB) error {
		if kvPair.HasIssuerInfos() {
			// save issuer info versions
			issuerInfoVersions := make([]IssuerInfoVersion, len(kvPair.IssuerInfos))
			for i, info := range kvPair.IssuerInfos {
				var oldInfo IssuerInfo
				err := tx.Model(IssuerInfo{}).WithContext(ctx).Where("lock_hash = ?", info.LockHash).First(&oldInfo).Error
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
					issuerInfoVersions[i] = IssuerInfoVersion{
						BlockNumber:  info.BlockNumber,
						LockHash:     info.LockHash,
						Version:      info.Version,
						Name:         info.Name,
						Avatar:       info.Avatar,
						Description:  info.Description,
						Localization: info.Localization,
						ActionType:   0,
						TxIndex:      info.TxIndex,
					}
				} else {
					issuerInfoVersions[i] = IssuerInfoVersion{
						OldBlockNumber:  oldInfo.BlockNumber,
						BlockNumber:     info.BlockNumber,
						LockHash:        info.LockHash,
						OldVersion:      oldInfo.Version,
						Version:         info.Version,
						OldName:         oldInfo.Name,
						Name:            info.Name,
						OldAvatar:       oldInfo.Avatar,
						Avatar:          info.Avatar,
						OldDescription:  oldInfo.Description,
						Description:     info.Description,
						OldLocalization: oldInfo.Localization,
						Localization:    info.Localization,
						ActionType:      1,
						TxIndex:         info.TxIndex,
					}
				}
			}
			if err := tx.Model(IssuerInfoVersion{}).WithContext(ctx).Create(issuerInfoVersions).Error; err != nil {
				return err
			}
			// insert issuer info
			issuerInfos := make([]IssuerInfo, len(kvPair.IssuerInfos))
			for i, issuer := range kvPair.IssuerInfos {
				issuerInfos[i] = IssuerInfo{
					BlockNumber:  issuer.BlockNumber,
					LockHash:     issuer.LockHash,
					Version:      issuer.Version,
					Name:         issuer.Name,
					Avatar:       issuer.Avatar,
					Description:  issuer.Description,
					Localization: issuer.Localization,
				}
			}
			if err := tx.Model(IssuerInfo{}).WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "lock_hash"}},
				UpdateAll: true,
			}).Create(issuerInfos).Error; err != nil {
				return err
			}
		}
		if kvPair.HasClassInfos() {
			// save class info versions
			classInfoVersions := make([]ClassInfoVersion, len(kvPair.ClassInfos))
			for i, info := range kvPair.ClassInfos {
				var oldInfo ClassInfo
				err := tx.Model(ClassInfo{}).WithContext(ctx).Where("cota_id = ?", info.CotaId).First(&oldInfo).Error
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
					classInfoVersions[i] = ClassInfoVersion{
						BlockNumber:    0,
						CotaId:         info.CotaId,
						Version:        info.Version,
						Name:           info.Name,
						Symbol:         info.Symbol,
						Description:    info.Description,
						Image:          info.Image,
						Audio:          info.Audio,
						Video:          info.Video,
						Model:          info.Model,
						Characteristic: info.Characteristic,
						Properties:     info.Properties,
						Localization:   info.Localization,
						ActionType:     0,
						TxIndex:        info.TxIndex,
					}
				} else {
					classInfoVersions[i] = ClassInfoVersion{
						OldBlockNumber:    oldInfo.BlockNumber,
						BlockNumber:       info.BlockNumber,
						CotaId:            info.CotaId,
						OldVersion:        oldInfo.Version,
						Version:           info.Version,
						OldName:           oldInfo.Name,
						Name:              info.Name,
						OldSymbol:         oldInfo.Symbol,
						Symbol:            info.Symbol,
						OldDescription:    oldInfo.Description,
						Description:       info.Description,
						OldImage:          oldInfo.Image,
						Image:             info.Image,
						OldAudio:          oldInfo.Audio,
						Audio:             info.Audio,
						OldVideo:          oldInfo.Video,
						Video:             info.Video,
						OldModel:          oldInfo.Model,
						Model:             info.Model,
						OldCharacteristic: oldInfo.Characteristic,
						Characteristic:    info.Characteristic,
						OldProperties:     oldInfo.Properties,
						Properties:        info.Properties,
						OldLocalization:   oldInfo.Localization,
						Localization:      info.Localization,
						ActionType:        1,
						TxIndex:           info.TxIndex,
					}
				}
			}

			if err := tx.Model(ClassInfoVersion{}).WithContext(ctx).Create(classInfoVersions).Error; err != nil {
				return err
			}
			// insert class info
			classInfos := make([]ClassInfo, len(kvPair.ClassInfos))
			audios := make([]TokenClassAudio, 0)

			for i, class := range kvPair.ClassInfos {
				classInfos[i] = ClassInfo{
					BlockNumber:    class.BlockNumber,
					CotaId:         class.CotaId,
					Version:        class.Version,
					Name:           class.Name,
					Symbol:         class.Symbol,
					Description:    class.Description,
					Image:          class.Image,
					Audio:          class.Audio,
					Video:          class.Video,
					Model:          class.Model,
					Characteristic: class.Characteristic,
					Properties:     class.Properties,
					Localization:   class.Localization,
				}

				for i, audio := range class.Audios {
					audios = append(audios, TokenClassAudio{
						CotaId: class.CotaId,
						Url:    audio.Url,
						Name:   audio.Name,
						Idx:    uint32(i),
					})
				}
			}

			if err := tx.Model(ClassInfo{}).WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cota_id"}},
				UpdateAll: true,
			}).Create(classInfos).Error; err != nil {
				return err
			}

			// insert audios
			if len(audios) > 0 {
				if err := rp.upsertAudios(tx, audios, ctx); err != nil {
					return err
				}
			}
		}
		if kvPair.HasJoyIDInfos() {
			// save joyID info versions and subkey info versions
			joyIDInfoVersions := make([]JoyIDInfoVersion, len(kvPair.JoyIDInfos))
			var subKeyVersions []SubKeyInfoVersion
			for i, info := range kvPair.JoyIDInfos {
				var oldInfo JoyIDInfo
				err := tx.Model(JoyIDInfo{}).WithContext(ctx).Where("lock_hash = ?", info.LockHash).First(&oldInfo).Error
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				if len(oldInfo.Name) > 240 || len(info.Name) > 240 {
					continue
				}
				if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
					joyIDInfoVersions[i] = JoyIDInfoVersion{
						BlockNumber:          0,
						Version:              info.Version,
						Name:                 info.Name,
						Description:          info.Description,
						Avatar:               info.Avatar,
						PubKey:               info.PubKey,
						CredentialId:         info.CredentialId,
						Alg:                  info.Alg,
						FrontEnd:             info.FrontEnd,
						DeviceName:           info.DeviceName,
						DeviceType:           info.DeviceType,
						CotaCellId:           info.CotaCellId,
						LockHash:             info.LockHash,
						Extension:            info.Extension,
						ActionType:           0,
						TxIndex:              info.TxIndex,
						DerivationCId:        info.DerivationCId,
						DerivationCommitment: info.DerivationCommitment,
					}
				} else {
					joyIDInfoVersions[i] = JoyIDInfoVersion{
						OldBlockNumber:          oldInfo.BlockNumber,
						BlockNumber:             info.BlockNumber,
						LockHash:                info.LockHash,
						OldVersion:              oldInfo.Version,
						Version:                 info.Version,
						OldName:                 oldInfo.Name,
						Name:                    info.Name,
						OldAvatar:               oldInfo.Avatar,
						Avatar:                  info.Avatar,
						OldDescription:          oldInfo.Description,
						Description:             info.Description,
						OldExtension:            oldInfo.Extension,
						Extension:               info.Extension,
						PubKey:                  info.PubKey,
						CredentialId:            info.CredentialId,
						Alg:                     info.Alg,
						OldFrontEnd:             oldInfo.FrontEnd,
						FrontEnd:                info.FrontEnd,
						OldDeviceName:           oldInfo.DeviceName,
						DeviceName:              info.DeviceName,
						OldDeviceType:           oldInfo.DeviceType,
						DeviceType:              info.DeviceType,
						CotaCellId:              info.CotaCellId,
						ActionType:              1,
						TxIndex:                 info.TxIndex,
						OldDerivationCId:        oldInfo.DerivationCId,
						DerivationCId:           info.DerivationCId,
						OldDerivationCommitment: oldInfo.DerivationCommitment,
						DerivationCommitment:    info.DerivationCommitment,
					}
				}

				for _, sub := range info.SubKeys {
					var oldSub SubKeyInfo
					err := tx.Model(SubKeyInfo{}).WithContext(ctx).Where("lock_hash = ? and pub_key = ?", sub.LockHash, sub.PubKey).First(&oldSub).Error
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return err
					}
					if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
						subKeyVersions = append(subKeyVersions, SubKeyInfoVersion{
							BlockNumber:          sub.BlockNumber,
							PubKey:               sub.PubKey,
							CredentialId:         sub.CredentialId,
							Alg:                  sub.Alg,
							FrontEnd:             sub.FrontEnd,
							DeviceName:           sub.DeviceName,
							DeviceType:           sub.DeviceType,
							LockHash:             sub.LockHash,
							ActionType:           0,
							TxIndex:              info.TxIndex,
							DerivationCId:        sub.DerivationCId,
							DerivationCommitment: sub.DerivationCommitment,
						})
					} else {
						subKeyVersions = append(subKeyVersions, SubKeyInfoVersion{
							OldBlockNumber:          oldSub.BlockNumber,
							BlockNumber:             sub.BlockNumber,
							LockHash:                sub.LockHash,
							PubKey:                  sub.PubKey,
							CredentialId:            sub.CredentialId,
							Alg:                     sub.Alg,
							OldFrontEnd:             oldSub.FrontEnd,
							FrontEnd:                sub.FrontEnd,
							OldDeviceName:           oldSub.DeviceName,
							DeviceName:              sub.DeviceName,
							OldDeviceType:           oldSub.DeviceType,
							DeviceType:              sub.DeviceType,
							ActionType:              1,
							TxIndex:                 info.TxIndex,
							OldDerivationCId:        oldSub.DerivationCId,
							DerivationCId:           sub.DerivationCId,
							OldDerivationCommitment: oldSub.DerivationCommitment,
							DerivationCommitment:    sub.DerivationCommitment,
						})
					}
				}
			}
			if err := tx.Model(JoyIDInfoVersion{}).WithContext(ctx).Create(joyIDInfoVersions).Error; err != nil {
				return err
			}
			if len(subKeyVersions) > 0 {
				if err := tx.Model(SubKeyInfoVersion{}).WithContext(ctx).Create(subKeyVersions).Error; err != nil {
					return err
				}
			}

			// insert joyID info and subkey info
			var subKeys []SubKeyInfo
			joyIDInfos := make([]JoyIDInfo, len(kvPair.JoyIDInfos))
			for i, joyID := range kvPair.JoyIDInfos {
				joyIDInfos[i] = JoyIDInfo{
					BlockNumber:          joyID.BlockNumber,
					LockHash:             joyID.LockHash,
					Version:              joyID.Version,
					Name:                 joyID.Name,
					Description:          joyID.Description,
					Avatar:               joyID.Avatar,
					PubKey:               joyID.PubKey,
					CredentialId:         joyID.CredentialId,
					Alg:                  joyID.Alg,
					FrontEnd:             joyID.FrontEnd,
					DeviceName:           joyID.DeviceName,
					DeviceType:           joyID.DeviceType,
					CotaCellId:           joyID.CotaCellId,
					Extension:            joyID.Extension,
					DerivationCId:        joyID.DerivationCId,
					DerivationCommitment: joyID.DerivationCommitment,
				}
				for _, subKey := range joyID.SubKeys {
					subKeys = append(subKeys, SubKeyInfo{
						LockHash:             joyID.LockHash,
						BlockNumber:          joyID.BlockNumber,
						PubKey:               subKey.PubKey,
						CredentialId:         subKey.CredentialId,
						Alg:                  subKey.Alg,
						FrontEnd:             subKey.FrontEnd,
						DeviceName:           subKey.DeviceName,
						DeviceType:           subKey.DeviceType,
						DerivationCId:        subKey.DerivationCId,
						DerivationCommitment: subKey.DerivationCommitment,
					})
				}
			}
			if err := tx.Model(JoyIDInfo{}).WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "lock_hash"}},
				UpdateAll: true,
			}).Create(joyIDInfos).Error; err != nil {
				return err
			}

			if len(subKeys) > 0 {
				for _, subKey := range subKeys {
					var oldSubkey SubKeyInfo
					err := tx.Model(SubKeyInfo{}).WithContext(ctx).Where("lock_hash = ? and pub_key = ? and credential_id = ?", subKey.LockHash, subKey.PubKey, subKey.CredentialId).First(&oldSubkey).Error
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return err
					}
					if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
						if err := tx.Model(SubKeyInfo{}).WithContext(ctx).Create(&subKey).Error; err != nil {
							return err
						}
					}
					if err == nil {
						if err := tx.Model(&oldSubkey).WithContext(ctx).Updates(subKey).Error; err != nil {
							return err
						}
					}
				}
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

func (rp kvPairRepo) upsertAudios(tx *gorm.DB, audios []TokenClassAudio, ctx context.Context) error {
	for _, audio := range audios {
		if err := tx.Model(TokenClassAudio{}).WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "cota_id"}, {Name: "idx"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "url", "updated_at"}),
		}).Create(&audio).Error; err != nil {
			return err
		}
	}
	return nil
}

func (rp kvPairRepo) RestoreMetadataKvPairs(ctx context.Context, blockNumber uint64) error {
	return rp.data.db.Transaction(func(tx *gorm.DB) error {
		// delete all issuer info by the block number
		if err := tx.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(IssuerInfo{}).Error; err != nil {
			return err
		}
		// update issuer info to the data before the last update
		var issuerInfoVersions []IssuerInfoVersion
		if err := tx.Model(IssuerInfoVersion{}).WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Group("lock_hash").Order("tx_index").Find(&issuerInfoVersions).Error; err != nil {
			return err
		}
		var updatedIssuerInfos []IssuerInfo
		for _, version := range issuerInfoVersions {
			updatedIssuerInfos = append(updatedIssuerInfos, IssuerInfo{
				BlockNumber:  version.OldBlockNumber,
				LockHash:     version.LockHash,
				Version:      version.OldVersion,
				Name:         version.OldName,
				Avatar:       version.OldAvatar,
				Description:  version.OldDescription,
				Localization: version.OldLocalization,
				UpdatedAt:    time.Now().UTC(),
			})
		}
		if len(updatedIssuerInfos) > 0 {
			if err := tx.Model(IssuerInfo{}).WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "lock_hash"}},
				UpdateAll: true,
			}).Create(updatedIssuerInfos).Error; err != nil {
				return err
			}
		}
		// delete all class info by the block number
		if err := tx.Debug().WithContext(ctx).Where("block_number = ?", blockNumber).Delete(ClassInfo{}).Error; err != nil {
			return err
		}
		var classInfoVersions []ClassInfoVersion
		if err := tx.Model(ClassInfoVersion{}).WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Group("cota_id").Order("tx_index").Find(&classInfoVersions).Error; err != nil {
			return err
		}
		var updatedClassInfos []ClassInfo
		for _, version := range classInfoVersions {
			updatedClassInfos = append(updatedClassInfos, ClassInfo{
				BlockNumber:    version.OldBlockNumber,
				CotaId:         version.CotaId,
				Version:        version.OldVersion,
				Name:           version.OldName,
				Symbol:         version.OldSymbol,
				Description:    version.OldDescription,
				Image:          version.OldImage,
				Audio:          version.OldAudio,
				Video:          version.OldVideo,
				Model:          version.OldModel,
				Characteristic: version.OldCharacteristic,
				Properties:     version.OldProperties,
				Localization:   version.OldLocalization,
				UpdatedAt:      time.Now().UTC(),
			})
		}
		if len(updatedClassInfos) > 0 {
			if err := tx.Debug().Model(ClassInfo{}).WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "cota_id"}},
				UpdateAll: true,
			}).Create(updatedClassInfos).Error; err != nil {
				return err
			}
		}
		// delete all joyID info by the block number
		if err := tx.Debug().WithContext(ctx).Where("block_number = ?", blockNumber).Delete(JoyIDInfo{}).Error; err != nil {
			return err
		}
		var joyIDInfoVersions []JoyIDInfoVersion
		if err := tx.Model(JoyIDInfoVersion{}).WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Group("lock_hash").Order("tx_index").Find(&joyIDInfoVersions).Error; err != nil {
			return err
		}
		var updatedJoyIDInfos []JoyIDInfo
		for _, version := range joyIDInfoVersions {
			updatedJoyIDInfos = append(updatedJoyIDInfos, JoyIDInfo{
				BlockNumber:          version.OldBlockNumber,
				LockHash:             version.LockHash,
				Version:              version.OldVersion,
				Name:                 version.OldName,
				Avatar:               version.OldAvatar,
				Description:          version.OldDescription,
				Extension:            version.OldExtension,
				PubKey:               version.PubKey,
				CredentialId:         version.CredentialId,
				Alg:                  version.Alg,
				FrontEnd:             version.OldFrontEnd,
				DeviceName:           version.OldDeviceName,
				DeviceType:           version.OldDeviceType,
				CotaCellId:           version.CotaCellId,
				DerivationCId:        version.OldDerivationCId,
				DerivationCommitment: version.OldDerivationCommitment,
				UpdatedAt:            time.Now().UTC(),
			})
		}
		if len(updatedJoyIDInfos) > 0 {
			if err := tx.Debug().Model(JoyIDInfo{}).WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "lock_hash"}},
				UpdateAll: true,
			}).Create(updatedJoyIDInfos).Error; err != nil {
				return err
			}
		}
		// delete all subkey info by the block number
		if err := tx.Debug().WithContext(ctx).Where("block_number = ?", blockNumber).Delete(SubKeyInfo{}).Error; err != nil {
			return err
		}
		var subKeyInfoVersions []SubKeyInfoVersion
		if err := tx.Model(SubKeyInfoVersion{}).WithContext(ctx).Where("block_number = ? and action_type = ?", blockNumber, 1).Group("lock_hash").Order("tx_index").Find(&subKeyInfoVersions).Error; err != nil {
			return err
		}
		var updatedSubKeyInfos []SubKeyInfo
		for _, version := range subKeyInfoVersions {
			updatedSubKeyInfos = append(updatedSubKeyInfos, SubKeyInfo{
				BlockNumber:          version.OldBlockNumber,
				LockHash:             version.LockHash,
				PubKey:               version.PubKey,
				CredentialId:         version.CredentialId,
				Alg:                  version.Alg,
				FrontEnd:             version.OldFrontEnd,
				DeviceName:           version.OldDeviceName,
				DeviceType:           version.OldDeviceType,
				DerivationCId:        version.OldDerivationCId,
				DerivationCommitment: version.OldDerivationCommitment,
				UpdatedAt:            time.Now().UTC(),
			})
		}
		if len(updatedSubKeyInfos) > 0 {
			if err := tx.Debug().Model(SubKeyInfo{}).WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "lock_hash"}},
				UpdateAll: true,
			}).Create(updatedSubKeyInfos).Error; err != nil {
				return err
			}
		}
		// delete check info
		if err := tx.Debug().WithContext(ctx).Where("block_number = ? and check_type = ?", blockNumber, biz.SyncMetadata).Delete(CheckInfo{}).Error; err != nil {
			return err
		}
		return nil
	})
}
