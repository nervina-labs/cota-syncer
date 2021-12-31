package config

import (
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
	Mode        string `mapstructure:"mode"`
}

type CkbNode struct {
	RpcUrl string `mapstructure:"rpc_url"`
	Mode   string `mapstructure:"mode"`
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
