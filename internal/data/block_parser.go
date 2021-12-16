package data

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

type BlockParser struct {
	claimedCotaUsecase    *biz.ClaimedCotaNftKvPairUsecase
	defineCotaUsecase     *biz.DefineCotaNftKvPairUsecase
	holdCotaUsecase       *biz.HoldCotaNftKvPairUsecase
	registerCotaUsecase   *biz.RegisterCotaKvPairUsecase
	withdrawCotaUsecase   *biz.WithdrawCotaNftKvPairUsecase
	cotaWitnessArgsParser CotaWitnessArgsParser
}

func NewBlockParser(claimedCotaUsecase *biz.ClaimedCotaNftKvPairUsecase, defineCotaUsecase *biz.DefineCotaNftKvPairUsecase, holdCotaUsecase *biz.HoldCotaNftKvPairUsecase, registerCotaUsecase *biz.RegisterCotaKvPairUsecase, withdrawCotaUsecase *biz.WithdrawCotaNftKvPairUsecase) BlockParser {
	return BlockParser{
		claimedCotaUsecase:  claimedCotaUsecase,
		defineCotaUsecase:   defineCotaUsecase,
		holdCotaUsecase:     holdCotaUsecase,
		registerCotaUsecase: registerCotaUsecase,
		withdrawCotaUsecase: withdrawCotaUsecase,
	}
}

func (bp BlockParser) Parse(block *ckbTypes.Block, systemScripts SystemScripts) (biz.KvPair, error) {
	kvPair := biz.KvPair{
		Registers:     []biz.RegisterCotaKvPair{},
		DefineCotas:   []biz.DefineCotaNftKvPair{},
		HoldCotas:     []biz.HoldCotaNftKvPair{},
		WithdrawCotas: []biz.WithdrawCotaNftKvPair{},
		ClaimedCotas:  []biz.ClaimedCotaNftKvPair{},
	}
	for _, tx := range block.Transactions {
		// ParseRegistryEntries
		if bp.hasCotaRegistryCell(tx.Outputs, systemScripts.CotaRegistryType) {
			registers, err := bp.registerCotaUsecase.ParseRegistryEntries(context.TODO(), block.Header.Number, tx)
			if err != nil {
				return kvPair, err
			}
			kvPair.Registers = append(kvPair.Registers, registers...)
		}
	}
	return kvPair, nil
}

func (bp BlockParser) isCotaRegistryCell(output *ckbTypes.CellOutput, registryType SystemScript) bool {
	return output.Type.CodeHash == registryType.CodeHash && output.Type.HashType == registryType.HashType
}

func (bp BlockParser) hasCotaRegistryCell(outputs []*ckbTypes.CellOutput, registryType SystemScript) (result bool) {
	for _, output := range outputs {
		if result = bp.isCotaRegistryCell(output, registryType); result {
			break
		}
	}
	return result
}
