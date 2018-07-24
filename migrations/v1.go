package migrations

import (
	"github.com/0x111/telegram-rss-bot/db"
	log "github.com/sirupsen/logrus"
)

func V1() (bool, error) {
	DB := db.GetDB()
	sqlStmt := `
	create table if not exists feeds(id integer primary key autoincrement, name text, url text, chatid integer, userid integer);
	create table if not exists feedData(id integer primary key autoincrement, feedid integer, title text, link text, publishedDate timestamp, published integer);
	`

	_, err := DB.Exec(sqlStmt)

	if err != nil {
		log.Error("Query error", err)
		return false, err
	}
	return true, nil
}
