package data

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	mMsql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/wire"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/config"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers
var ProviderSet = wire.NewSet(NewData, NewDBMigration, NewCheckInfoRepo, NewRegisterCotaKvPairRepo, NewDefineCotaNftKvPairRepo, NewHoldCotaNftKvPairRepo, NewWithdrawCotaNftKvPairRepo, NewClaimedCotaNftKvPairRepo, NewKvPairRepo, NewSystemScripts, NewCkbNodeClient, NewBlockParser, NewCotaWitnessArgsParser)

type Data struct {
	db *gorm.DB
}

type Option func(*Data)

func NewData(conf *config.Database, logger *logger.Logger) (*Data, func(), error) {
	db, err := gorm.Open(mysql.Open(conf.Dsn), &gorm.Config{})
	if err != nil {
		logger.Errorf(context.TODO(), "failed opening connection to mysql: %v", err)
		return nil, nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorf(context.TODO(), "failed get sql db: %v", err)
		return nil, nil, err
	}
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(conf.ConnMaxLifeTime)

	return &Data{
			db: db,
		}, func() {
			if err := sqlDB.Close(); err != nil {
				logger.Error(context.TODO(), err)
			}
			logger.Info(context.TODO(), "successfully closed the database")
		}, nil
}

type CkbNodeClient struct {
	Rpc rpc.Client
}

func NewCkbNodeClient(conf *config.CkbNode, logger *logger.Logger) (*CkbNodeClient, error) {
	client, err := rpc.Dial(conf.RpcUrl)
	if err != nil {
		logger.Errorf(context.TODO(), "failed to connect to the ckb node")
		return nil, err
	}
	return &CkbNodeClient{
		Rpc: client,
	}, nil
}

type DBMigration struct {
	data   *Data
	logger *logger.Logger
}

func (m *DBMigration) Up() error {
	sqlDB, err := m.data.db.DB()
	if err != nil {
		m.logger.Errorf(context.TODO(), "failed get sql db: %v", err)
		return err
	}
	driver, err := mMsql.WithInstance(sqlDB, &mMsql.Config{})
	if err != nil {
		return err
	}
	migration, err := migrate.NewWithDatabaseInstance("file://./internal/db/migrations", m.data.db.Migrator().CurrentDatabase(), driver)
	if err != nil {
		return err
	}
	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		m.logger.Infof(context.TODO(), "migration failed: %v", err)
		return err
	}
	return nil
}

type SystemScriptOption func(o *SystemScripts)

type SystemScript struct {
	CodeHash ckbTypes.Hash
	HashType ckbTypes.ScriptHashType
	OutPoint ckbTypes.OutPoint
	DepType  ckbTypes.DepType
}

type SystemScripts struct {
	CotaRegistryType SystemScript
	CotaType         SystemScript
}

func NewSystemScripts(client *CkbNodeClient, logger *logger.Logger) SystemScripts {
	chainInfo, err := client.Rpc.GetBlockchainInfo(context.Background())
	if err != nil {
		logger.Fatalf(context.Background(), "RPC get_blockchain_info error")
	}
	return SystemScripts{
		CotaRegistryType: cotaRegistryScript(chainInfo.Chain),
		CotaType:         cotaTypeScript(chainInfo.Chain),
	}
}

func cotaRegistryScript(chain string) SystemScript {
	if chain == "ckb" {
		return SystemScript{
			CodeHash: ckbTypes.HexToHash("0x"),
			HashType: ckbTypes.HashTypeType,
			OutPoint: ckbTypes.OutPoint{
				TxHash: ckbTypes.HexToHash("0x"),
				Index:  0,
			},
			DepType: ckbTypes.DepTypeDepGroup,
		}
	}
	return SystemScript{
		CodeHash: ckbTypes.HexToHash("0x3a6897ab78ad10d028d0c5ef375545e66bfdffd01f3a369b5b07906078e04f6d"),
		HashType: ckbTypes.HashTypeType,
		OutPoint: ckbTypes.OutPoint{
			TxHash: ckbTypes.HexToHash("0x4410efbdfb83c58198a10eae621a3169c4f8f776cb4c2dd61b69947b1f4b922a"),
			Index:  0,
		},
		DepType: ckbTypes.DepTypeDepGroup,
	}
}

func cotaTypeScript(chain string) SystemScript {
	if chain == "ckb" {
		return SystemScript{
			CodeHash: ckbTypes.HexToHash("0x"),
			HashType: ckbTypes.HashTypeType,
			OutPoint: ckbTypes.OutPoint{
				TxHash: ckbTypes.HexToHash("0x"),
				Index:  0,
			},
			DepType: ckbTypes.DepTypeDepGroup,
		}
	}
	return SystemScript{
		CodeHash: ckbTypes.HexToHash("0x"),
		HashType: ckbTypes.HashTypeType,
		OutPoint: ckbTypes.OutPoint{
			TxHash: ckbTypes.HexToHash("0x"),
			Index:  0,
		},
		DepType: ckbTypes.DepTypeDepGroup,
	}
}

func (m *DBMigration) Down() error {
	sqlDB, err := m.data.db.DB()
	if err != nil {
		m.logger.Errorf(context.TODO(), "failed get sql db: %v", err)
		return err
	}
	driver, err := mMsql.WithInstance(sqlDB, &mMsql.Config{})
	if err != nil {
		return err
	}
	migration, err := migrate.NewWithDatabaseInstance("file://../db/migrations", m.data.db.Migrator().CurrentDatabase(), driver)
	if err != nil {
		return err
	}
	return migration.Down()
}

func NewDBMigration(data *Data, logger *logger.Logger) *DBMigration {
	return &DBMigration{
		data:   data,
		logger: logger,
	}
}
