package model

import "time"

type Notification struct {
	ID         string    `json:"id"`
	Date       time.Time `json:"-"`
	DateExport string    `json:"date"`
	Text       string    `json:"text"`
	UserID     string    `json:"user_id"`
}
