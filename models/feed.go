package models

// Basic Feed Model for selects
type Feed struct {
	ID     int
	Name   string
	Url    string
	ChatID int64
	UserID int
}
