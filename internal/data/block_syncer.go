package data

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

type BlockSyncer struct {
	claimedCotaUsecase    *biz.ClaimedCotaNftKvPairUsecase
	defineCotaUsecase     *biz.DefineCotaNftKvPairUsecase
	holdCotaUsecase       *biz.HoldCotaNftKvPairUsecase
	registerCotaUsecase   *biz.RegisterCotaKvPairUsecase
	withdrawCotaUsecase   *biz.WithdrawCotaNftKvPairUsecase
	cotaWitnessArgsParser CotaWitnessArgsParser
	kvPairUsecase         *biz.SyncKvPairUsecase
}

func NewBlockParser(claimedCotaUsecase *biz.ClaimedCotaNftKvPairUsecase, defineCotaUsecase *biz.DefineCotaNftKvPairUsecase, holdCotaUsecase *biz.HoldCotaNftKvPairUsecase, registerCotaUsecase *biz.RegisterCotaKvPairUsecase, withdrawCotaUsecase *biz.WithdrawCotaNftKvPairUsecase, kvPairUsecase *biz.SyncKvPairUsecase) BlockSyncer {
	return BlockSyncer{
		claimedCotaUsecase:  claimedCotaUsecase,
		defineCotaUsecase:   defineCotaUsecase,
		holdCotaUsecase:     holdCotaUsecase,
		registerCotaUsecase: registerCotaUsecase,
		withdrawCotaUsecase: withdrawCotaUsecase,
		kvPairUsecase:       kvPairUsecase,
	}
}

func (bp BlockSyncer) Sync(ctx context.Context, block *ckbTypes.Block, systemScripts SystemScripts) error {
	for index, tx := range block.Transactions {
		kvPair := biz.KvPair{}
		// ParseRegistryEntries TODO 拆到独立到 repo 中
		if bp.hasCotaRegistryCell(tx.Outputs, systemScripts.CotaRegistryType) {
			registers, err := bp.registerCotaUsecase.ParseRegistryEntries(ctx, block.Header.Number, tx)
			if err != nil {
				return err
			}
			kvPair.Registers = append(kvPair.Registers, registers...)
		}
		entries, err := bp.cotaWitnessArgsParser.Parse(tx, systemScripts.CotaType)
		if err != nil {
			return err
		}
		pairs, err := bp.parseCotaEntries(entries)
		kvPair.DefineCotas = pairs.DefineCotas
		kvPair.WithdrawCotas = pairs.WithdrawCotas
		err = bp.kvPairUsecase.CreateKvPairs(ctx, index, &pairs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bp BlockSyncer) isCotaRegistryCell(output *ckbTypes.CellOutput, registryType SystemScript) bool {
	return output.Type.CodeHash == registryType.CodeHash && output.Type.HashType == registryType.HashType
}

func (bp BlockSyncer) hasCotaRegistryCell(outputs []*ckbTypes.CellOutput, registryType SystemScript) (result bool) {
	for _, output := range outputs {
		if result = bp.isCotaRegistryCell(output, registryType); result {
			break
		}
	}
	return result
}

func (bp BlockSyncer) Rollback(ctx context.Context, blockNumber uint64) error {
	return bp.kvPairUsecase.RestoreKvPairs(ctx, blockNumber)
}

func (bp BlockSyncer) parseCotaEntries(entries [][]byte) (biz.KvPair, error) {
	return biz.KvPair{}, nil
}
