package models

import "time"

// Basic Feed Data Model for the Feed Items which we get from the Feed
type FeedData struct {
	ID            int
	Name          string
	FeedID        int
	Title         string
	Link          string
	Published     bool
	PublishedDate time.Time
	ChatID        int64
}
