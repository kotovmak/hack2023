package model

import "time"

type Slot struct {
	ID           int           `json:"id"`
	Time         string        `json:"time"`
	Date         time.Time     `json:"-"`
	DateExport   string        `json:"date"`
	Consultation *Consultation `json:"consultation,omitempty"`
}
