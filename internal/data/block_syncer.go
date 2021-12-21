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
	mintCotaUsecase       biz.MintCotaKvPairUsecase
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
		pairs, err := bp.parseCotaEntries(block.Header.Number, entries)
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
	return output.Type.CodeHash == registryType.CodeHash && output.Type.HashType == registryType.HashType && argsEq(output.Type.Args, registryType.Args)
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

func (bp BlockSyncer) parseCotaEntries(blockNumber uint64, entries []biz.Entry) (biz.KvPair, error) {
	var kvPair biz.KvPair
	for _, entry := range entries {
		switch entry.Witness[0] {
		//	创建 DefineCota Kv pairs
		case 1:
			defineCotas, err := bp.defineCotaUsecase.ParseDefineCotaEntries(blockNumber, entry)
			if err != nil {
				return kvPair, err
			}
			kvPair.DefineCotas = append(kvPair.DefineCotas, defineCotas...)
		//	更新 DefineCota Kv pairs 创建 withdrawCota kv pairs
		case 2:
			updatedDefineCotas, withdrawCotas, err := bp.mintCotaUsecase.ParseMintCotaEntries(blockNumber, entry)
			if err != nil {
				return kvPair, err
			}
			kvPair.UpdatedDefineCotas = append(kvPair.UpdatedDefineCotas, updatedDefineCotas...)
			kvPair.WithdrawCotas = append(kvPair.WithdrawCotas, withdrawCotas...)
		//	删除 HoldCota kv pairs 创建 withdrawCota kv pairs
		case 3:
			withdrawCotas, err := bp.withdrawCotaUsecase.ParseWithdrawCotaEntries(blockNumber, entry)
			if err != nil {
				return kvPair, err
			}
			kvPair.WithdrawCotas = append(kvPair.WithdrawCotas, withdrawCotas...)
		//	创建 HoldCota kv pairs 与 claimedCota kv pairs
		case 4:
			holdCotas, claimedCotas, err := bp.claimedCotaUsecase.ParseClaimedCotaEntries(blockNumber, entry)
			if err != nil {
				return kvPair, err
			}
			kvPair.ClaimedCotas = append(kvPair.ClaimedCotas, claimedCotas...)
			kvPair.HoldCotas = append(kvPair.HoldCotas, holdCotas...)
		//	更新 HoldCota kv pairs
		case 5:
			holdCotas, err := bp.holdCotaUsecase.ParseHoldCotaEntries(blockNumber, entry)
			if err != nil {
				return kvPair, err
			}
			kvPair.UpdatedHoldCotas = append(kvPair.HoldCotas, holdCotas...)
		}
	}
	return kvPair, nil
}

func argsEq(args1, args2 []byte) bool {
	if len(args1) != len(args2) {
		return false
	}
	for i := range args1 {
		if args1[i] != args2[i] {
			return false
		}
	}
	return true
}
