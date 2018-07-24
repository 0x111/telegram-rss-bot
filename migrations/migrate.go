package migrations

import (
	log "github.com/sirupsen/logrus"
)

func Migrate() {
	// This is as bad as it seems, but it looked weird to use some complex tool for something as simple
	// as creating a few tables for this small app.
	_, err := V1()

	log.Debug("Running Migrations")

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Panic("There was a failure while creating the database structure. Please check your config.")
	}
}
