package config

import (
	"errors"
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
	SenderName          string        `mapstructure:"SENDER_NAME"`
	AuthEmail           string        `mapstructure:"AUTH_EMAIL"`
	AuthPassword        string        `mapstructure:"AUTH_PASSWORD"`
}

func LoadConfig(path string, fileName string) (config Config, err error) {
	viper.AddConfigPath(path)

	if fileName == "dev" {
		viper.SetConfigName("dev")
	} else if fileName == "prod" {
		viper.SetConfigName("prod")
	} else {
		err = errors.New("invalid environment variable")
		return
	}
	viper.SetConfigType(`env`)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
