package config

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Log                 *logrus.Logger
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	SetMaxIdleConns     int           `mapstructure:"SET_MAX_IDLE_CONNS"`
	SetMaxOpenConns     int           `mapstructure:"SET_MAX_OPEN_CONNS"`
	SetConnMaxLifeTime  int           `mapstructure:"SET_CONN_MAX_LIFE_TIME"`
	SetConnMaxIdleTime  int           `mapstructure:"SET_CONN_MAX_IDLE_TIME"`
}

func LoadConfigDev(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(`env_dev`)
	viper.SetConfigType(`env`)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

func LoadConfigProd(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(`env_prod`)
	viper.SetConfigType(`env`)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
