// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/nervina-labs/cota-syncer/internal/app"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/config"
	"github.com/nervina-labs/cota-syncer/internal/data"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"github.com/nervina-labs/cota-syncer/internal/service"
)

// Injectors from wire.go:

func initApp(database *config.Database, ckbNode *config.CkbNode, loggerLogger *logger.Logger) (*app.App, func(), error) {
	dataData, cleanup, err := data.NewData(database, loggerLogger)
	if err != nil {
		return nil, nil, err
	}
	checkInfoRepo := data.NewCheckInfoRepo(dataData, loggerLogger)
	checkInfoUsecase := biz.NewCheckInfoUsecase(checkInfoRepo, loggerLogger)
	ckbNodeClient, err := data.NewCkbNodeClient(ckbNode, loggerLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	systemScripts := data.NewSystemScripts(ckbNodeClient, loggerLogger)
	claimedCotaNftKvPairRepo := data.NewClaimedCotaNftKvPairRepo(dataData, loggerLogger)
	claimedCotaNftKvPairUsecase := biz.NewClaimedCotaNftKvPairUsecase(claimedCotaNftKvPairRepo, loggerLogger)
	defineCotaNftKvPairRepo := data.NewDefineCotaNftKvPairRepo(dataData, loggerLogger)
	defineCotaNftKvPairUsecase := biz.NewDefineCotaNftKvPairUsecase(defineCotaNftKvPairRepo, loggerLogger)
	holdCotaNftKvPairRepo := data.NewHoldCotaNftKvPairRepo(dataData, loggerLogger)
	holdCotaNftKvPairUsecase := biz.NewHoldCotaNftKvPairUsecase(holdCotaNftKvPairRepo, loggerLogger)
	registerCotaKvPairRepo := data.NewRegisterCotaKvPairRepo(dataData, loggerLogger)
	registerCotaKvPairUsecase := biz.NewRegisterCotaKvPairUsecase(registerCotaKvPairRepo, loggerLogger)
	withdrawCotaNftKvPairRepo := data.NewWithdrawCotaNftKvPairRepo(dataData, loggerLogger)
	withdrawCotaNftKvPairUsecase := biz.NewWithdrawCotaNftKvPairUsecase(withdrawCotaNftKvPairRepo, loggerLogger)
	cotaWitnessArgsParser := data.NewCotaWitnessArgsParser(ckbNodeClient)
	kvPairRepo := data.NewKvPairRepo(dataData, loggerLogger)
	syncKvPairUsecase := biz.NewSyncKvPairUsecase(kvPairRepo, loggerLogger)
	mintCotaKvPairRepo := data.NewMintCotaKvPairRepo(dataData, loggerLogger)
	mintCotaKvPairUsecase := biz.NewMintCotaKvPairUsecase(mintCotaKvPairRepo, loggerLogger)
	transferCotaKvPairRepo := data.NewTransferCotaKvPairRepo(dataData, loggerLogger)
	transferCotaKvPairUsecase := biz.NewTransferCotaKvPairUsecase(transferCotaKvPairRepo, loggerLogger)
	issuerInfoRepo := data.NewIssuerInfoRepo(dataData, loggerLogger)
	issuerInfoUsecase := biz.NewIssuerInfoUsecase(issuerInfoRepo, loggerLogger)
	classInfoRepo := data.NewClassInfoRepo(dataData, loggerLogger)
	classInfoUsecase := biz.NewClassInfoUsecase(classInfoRepo, loggerLogger)
	joyIDInfoRepo := data.NewJoyIDInfoRepo(dataData, loggerLogger)
	joyIDInfoUsecase := biz.NewJoyIDInfoUsecase(joyIDInfoRepo, loggerLogger)
	extensionPairRepo := data.NewExtensionKvPairRepo(dataData, loggerLogger)
	extensionPairUsecase := biz.NewExtensionPairUsecase(extensionPairRepo, loggerLogger)
	subKeyPairRepo := data.NewSubKeyKvPairRepo(dataData, loggerLogger)
	subKeyPairRepoUsecase := biz.NewSubKeyPairRepoUsecase(subKeyPairRepo, loggerLogger)
	blockSyncer := data.NewBlockSyncer(claimedCotaNftKvPairUsecase, defineCotaNftKvPairUsecase, holdCotaNftKvPairUsecase, registerCotaKvPairUsecase, withdrawCotaNftKvPairUsecase, cotaWitnessArgsParser, syncKvPairUsecase, mintCotaKvPairUsecase, transferCotaKvPairUsecase, issuerInfoUsecase, classInfoUsecase, joyIDInfoUsecase, extensionPairUsecase, subKeyPairRepoUsecase)
	blockSyncService := service.NewBlockSyncService(checkInfoUsecase, loggerLogger, ckbNodeClient, systemScripts, blockSyncer)
	checkInfoCleanerService := service.NewCheckInfoService(checkInfoUsecase, loggerLogger, ckbNodeClient)
	metadataSyncer := data.NewMetadataSyncer(syncKvPairUsecase, cotaWitnessArgsParser, issuerInfoUsecase, classInfoUsecase, joyIDInfoUsecase)
	metadataSyncService := service.NewMetadataSyncService(checkInfoUsecase, loggerLogger, ckbNodeClient, systemScripts, metadataSyncer)
	invalidDataRepo := data.NewInvalidDateRepo(dataData, loggerLogger)
	invalidDataUsecase := biz.NewInvalidDataUsecase(invalidDataRepo, loggerLogger)
	invalidDataCleaner := service.NewInvalidDataService(invalidDataUsecase, loggerLogger, ckbNodeClient)
	withdrawExtraInfoRepo := data.NewWithdrawExtraInfoRepo(dataData, loggerLogger)
	withdrawExtraInfoUsecase := biz.NewWithdrawExtraInfoUsecase(withdrawExtraInfoRepo, loggerLogger)
	withdrawExtraInfoService := service.NewWithdrawExtraInfoService(withdrawExtraInfoUsecase, loggerLogger, ckbNodeClient)
	registerLockScriptRepo := data.NewRegisterLockScriptRepo(dataData, loggerLogger)
	registerLockScriptUsecase := biz.NewRegisterLockScriptUsecase(registerLockScriptRepo, loggerLogger)
	registerLockService := service.NewRegisterLockService(registerLockScriptUsecase, loggerLogger, ckbNodeClient)
	dbMigration := data.NewDBMigration(dataData, loggerLogger)
	appApp := newApp(loggerLogger, blockSyncService, checkInfoCleanerService, metadataSyncService, invalidDataCleaner, withdrawExtraInfoService, registerLockService, dbMigration)
	return appApp, func() {
		cleanup()
	}, nil
}
