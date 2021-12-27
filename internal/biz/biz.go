package biz

import (
	"github.com/google/wire"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

var ProviderSet = wire.NewSet(NewCheckInfoUsecase, NewRegisterCotaKvPairUsecase, NewDefineCotaNftKvPairUsecase,
	NewHoldCotaNftKvPairUsecase, NewWithdrawCotaNftKvPairUsecase, NewClaimedCotaNftKvPairUsecase, NewSyncKvPairUsecase, NewMintCotaKvPairUsecase)

type Entry struct {
	Witness    []byte
	LockScript *ckbTypes.Script
	TxIndex    uint32
}
