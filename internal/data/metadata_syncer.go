package data

import (
	"context"
	"errors"

	"github.com/nervina-labs/cota-syncer/internal/biz"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

type MetadataSyncer struct {
	kvPairUsecase         *biz.SyncKvPairUsecase
	cotaWitnessArgsParser CotaWitnessArgsParser
	issuerInfoUsecase     *biz.IssuerInfoUsecase
	classInfoUsecase      *biz.ClassInfoUsecase
	joyIDInfoUsecase      *biz.JoyIDInfoUsecase
}

func NewMetadataSyncer(
	kvPairUsecase *biz.SyncKvPairUsecase, cotaWitnessArgsParser CotaWitnessArgsParser, issuerInfoUsecase *biz.IssuerInfoUsecase,
	classInfoUsecase *biz.ClassInfoUsecase, joyIDInfoUsecase *biz.JoyIDInfoUsecase) MetadataSyncer {

	return MetadataSyncer{
		kvPairUsecase:         kvPairUsecase,
		cotaWitnessArgsParser: cotaWitnessArgsParser,
		issuerInfoUsecase:     issuerInfoUsecase,
		classInfoUsecase:      classInfoUsecase,
		joyIDInfoUsecase:      joyIDInfoUsecase,
	}
}

func (bp MetadataSyncer) Sync(ctx context.Context, block *ckbTypes.Block, checkInfo biz.CheckInfo, systemScripts SystemScripts) error {
	var entryVec []biz.Entry
	for index, tx := range block.Transactions {
		entries, err := bp.cotaWitnessArgsParser.Parse(tx, uint32(index), systemScripts.CotaType)
		if err != nil && err.Error() == "No data" {
			continue
		} else if err != nil {
			return err
		}
		entryVec = append(entryVec, entries...)
	}
	pairs, err := bp.parseMetadata(ctx, block.Header.Number, entryVec)
	if err != nil {
		return err
	}
	err = bp.kvPairUsecase.CreateMetadataKvPairs(ctx, checkInfo, &pairs)
	if err != nil {
		return err
	}
	return nil
}

func (bp MetadataSyncer) Rollback(ctx context.Context, blockNumber uint64) error {
	return bp.kvPairUsecase.RestoreMetadataKvPairs(ctx, blockNumber)
}

func (bp MetadataSyncer) parseMetadata(ctx context.Context, blockNumber uint64, entries []biz.Entry) (biz.KvPair, error) {
	var (
		kvPair biz.KvPair
		ctMeta biz.CTMeta
		err    error
	)
	for _, entry := range entries {
		// Parse Issuer/Class/JoyID Metadata
		if len(entry.OutputType) > 0 {
			ctMeta, err = biz.ParseMetadata(entry.OutputType)
			if err != nil && len(entry.ExtraWitness) > 0 {
				ctMeta, err = biz.ParseMetadata(entry.ExtraWitness)
				if err != nil {
					continue
				}
			}
		} else if len(entry.ExtraWitness) > 0 {
			ctMeta, err = biz.ParseMetadata(entry.ExtraWitness)
			if err != nil {
				continue
			}
		} else {
			continue
		}
		switch ctMeta.Metadata.Type {
		case "issuer":
			issuerInfo, err := bp.issuerInfoUsecase.ParseMetadata(blockNumber, entry.TxIndex, entry.LockScript, ctMeta.Metadata.Data)
			if err != nil {
				return kvPair, err
			}
			kvPair.IssuerInfos = append(kvPair.IssuerInfos, issuerInfo)
		case "cota":
			classInfo, err := bp.classInfoUsecase.ParseMetadata(blockNumber, entry.TxIndex, ctMeta.Metadata.Data)
			if errors.Is(err, ErrInvalidClassInfo) {
				continue
			}

			if err != nil {
				return kvPair, err
			}
			kvPair.ClassInfos = append(kvPair.ClassInfos, classInfo)
		case "joy_id":
			joyIDInfo, err := bp.joyIDInfoUsecase.ParseMetadata(ctx, blockNumber, entry.TxIndex, entry.LockScript, ctMeta.Metadata.Data)
			if errors.Is(err, ErrInvalidJoyIDInfo) {
				continue
			}
			if err != nil {
				return kvPair, err
			}
			kvPair.JoyIDInfos = append(kvPair.JoyIDInfos, joyIDInfo)
		}
	}
	return kvPair, nil
}
