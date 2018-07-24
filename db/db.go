package db

import (
	"database/sql"
	"github.com/0x111/telegram-rss-bot/conf"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	db *sql.DB
)

func GetDB() *sql.DB {
	return db
}

func ConnectDB() *sql.DB {
	var err error
	config := conf.GetConfig()
	log.Debug("Connecting to the database")
	db, err = sql.Open("sqlite3", "file:"+config.GetString("db_path")+"?cache=shared")

	// Check error for database connection
	if err != nil {
		panic(err)
	}

	return db
}
