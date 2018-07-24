package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

func main() {

	db, err := sql.Open("sqlite3", "file:foo.db?cache=shared")

	// Check error for database connection
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create base tables
	sqlStmt := `
	drop table foo;
	
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	bot, err := tgbotapi.NewBotAPI("605498632:AAE_Eb_CJp_3k2m8t-1LPXrd0DdQoyCrnY0")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
