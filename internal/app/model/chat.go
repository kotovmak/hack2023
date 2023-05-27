package model

import "time"

type Message struct {
	ID       int       `json:"id"`
	Text     string    `json:"text"`
	Date     time.Time `json:"date"`
	SendByID int       `json:"send_by_id"`
	UserID   int       `json:"user_id"`
	Buttons  []Button  `json:"buttons"`
}

type Button struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}
