package config

import (
	"context"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"time"
)

type Data struct {
	Database Database
}

type Database struct {
	Driver          string
	Dsn             string
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifeTime time.Duration `mapstructure:"conn_max_lifetime"`
}

type App struct {
	LogSavePath string `mapstructure:"log_save_path"`
	LogFileName string `mapstructure:"log_file_name"`
	LogFileExt  string `mapstructure:"log_file_ext"`
}

type CkbNode struct {
	RpcUrl string `mapstructure:"rpc_url"`
	Mode   string `mapstructure:"mode"`
}

type SystemScriptOption func(o *SystemScripts)

type SystemScript struct {
	CodeHash ckbTypes.Hash
	HashType ckbTypes.ScriptHashType
	OutPoint ckbTypes.OutPoint
	DepType  ckbTypes.DepType
}

type SystemScripts struct {
	CompactRegistryType SystemScript
	CompactNFTType      SystemScript
}

func NewSystemScripts(client rpc.Client, options ...SystemScriptOption) (*SystemScripts, error) {
	chainInfo, err := client.GetBlockchainInfo(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "RPC get_blockchain_info error")
	}
	scripts := &SystemScripts{
		CompactRegistryType: cotaRegistryScript(chainInfo.Chain),
		CompactNFTType:      cotaTypeScript(chainInfo.Chain),
	}

	for _, option := range options {
		option(scripts)
	}

	return scripts, nil
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

type Config struct {
	vp *viper.Viper
}

func NewConfig() (*Config, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("configs/")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Config{vp: vp}, nil
}

func (s *Config) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}
