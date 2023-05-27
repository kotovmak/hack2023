package model

import "time"

type Message struct {
	ID         int       `json:"id"`
	Text       string    `json:"text" form:"text" validate:"required"`
	Date       time.Time `json:"-"`
	DateExport string    `json:"date_export"`
	SendByID   int       `json:"send_by_id"`
	UserID     int       `json:"user_id"`
}

type Button struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Link string `json:"link"`
}
