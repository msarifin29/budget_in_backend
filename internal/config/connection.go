package config

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Connection(log *logrus.Logger) *sql.DB {
	configuration, err := LoadConfig("../.", "prod")
	if err != nil {
		log.Fatal(err.Error())
	}
	// db, err := sql.Open(configuration.DBDriver, configuration.DBSource) for MySQL
	db, err := sql.Open(configuration.DBPostgresDriver, configuration.DBPostgresSource)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(configuration.SetMaxIdleConns)
	db.SetMaxOpenConns(configuration.SetMaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(configuration.SetConnMaxLifeTime))
	db.SetConnMaxIdleTime(time.Duration(configuration.SetConnMaxLifeTime))

	return db
}
