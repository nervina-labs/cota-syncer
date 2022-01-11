package data

import (
	"context"
	"encoding/hex"
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
	"os"
)

// ProviderSet is data providers
var ProviderSet = wire.NewSet(NewData, NewDBMigration, NewCheckInfoRepo, NewRegisterCotaKvPairRepo,
	NewDefineCotaNftKvPairRepo, NewHoldCotaNftKvPairRepo, NewWithdrawCotaNftKvPairRepo, NewClaimedCotaNftKvPairRepo,
	NewKvPairRepo, NewSystemScripts, NewCkbNodeClient, NewBlockParser, NewCotaWitnessArgsParser, NewMintCotaKvPairRepo)

type Data struct {
	db *gorm.DB
}

type Option func(*Data)

func NewData(conf *config.Database, logger *logger.Logger) (*Data, func(), error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = conf.Dsn
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
	Args     []byte
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
		args, _ := hex.DecodeString("")
		return SystemScript{
			CodeHash: ckbTypes.HexToHash("0x"),
			HashType: ckbTypes.HashTypeType,
			Args:     args,
			OutPoint: ckbTypes.OutPoint{
				TxHash: ckbTypes.HexToHash("0x"),
				Index:  0,
			},
			DepType: ckbTypes.DepTypeDepGroup,
		}
	}
	args, _ := hex.DecodeString("9da28b58954f6d710333b43832a151c5c3c47476")
	return SystemScript{
		CodeHash: ckbTypes.HexToHash("0x243e92edb5767b445560260b838261a2c79b7b40b806d6f86fa6f40a427b879c"),
		HashType: ckbTypes.HashTypeType,
		Args:     args,
		OutPoint: ckbTypes.OutPoint{
			TxHash: ckbTypes.HexToHash("0x6efe711f781d801aca8bf10378037ef18313837908f63eb6a4d1be382eaa4e55"),
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
		CodeHash: ckbTypes.HexToHash("0x0057b351bf489c3b93649566a5d5511d845f1744b2c1b6599f1198ed9d0d4378"),
		HashType: ckbTypes.HashTypeType,
		OutPoint: ckbTypes.OutPoint{
			TxHash: ckbTypes.HexToHash("0x6efe711f781d801aca8bf10378037ef18313837908f63eb6a4d1be382eaa4e55"),
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
