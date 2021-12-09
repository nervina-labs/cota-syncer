package biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCheckInfoUsecase, NewRegisterCotaKvPairUsecase, NewDefineCotaNftKvPairUsecase, NewHoldCotaNftKvPairUsecase, NewWithdrawCotaNftKvPairUsecase, NewClaimedCotaNftKvPairUsecase)
