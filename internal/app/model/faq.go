package model

import "time"

type FAQ struct {
	ID            int       `json:"id"`
	Question      string    `json:"question"`
	Answer        string    `json:"answer"`
	NadzorOrganID int       `json:"nadzor_organ_id"`
	ControlTypeID int       `json:"control_type_id"`
	Likes         int       `json:"likes"`
	Date          time.Time `json:"date"`
}
