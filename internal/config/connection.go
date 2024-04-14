package config

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func Connection(log *logrus.Logger) *sql.DB {
	configuration, err := LoadConfig("../.", "env_prod")
	if err != nil {
		log.Fatal(err.Error())
	}
	db, err := sql.Open(configuration.DBDriver, configuration.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(configuration.SetMaxIdleConns)
	db.SetMaxOpenConns(configuration.SetMaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(configuration.SetConnMaxLifeTime))
	db.SetConnMaxIdleTime(time.Duration(configuration.SetConnMaxLifeTime))

	return db
}
