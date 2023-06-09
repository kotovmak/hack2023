package model

import "time"

type Notification struct {
	ID             string    `json:"id"`
	Date           time.Time `json:"-"`
	DateExport     string    `json:"date"`
	Text           string    `json:"text"`
	UserID         int       `json:"user_id"`
	ConsultationID int       `json:"consultation_id"`
}
