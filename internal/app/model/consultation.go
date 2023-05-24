package model

import (
	"time"
)

type TypeList struct {
	NadzonOrgans []NadzonOrgan `json:"nadzor_organs"`
	Services     []Service     `json:"services"`
	PravActs     []PravAct     `json:"prav_acts"`
}

type NadzonOrgan struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	ControlTypes []ControlType `json:"control_types,omitempty"`
}

type Service struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ControlType struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	ConsultTopics []ConsultTopic `json:"consult_topics,omitempty"`
}

type ConsultTopic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PravAct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Consultation struct {
	NadzonOrganID  int       `json:"nadzor_organ_id" form:"nadzor_organ_id"`
	ControlTypeID  int       `json:"control_type_id" form:"control_type_id"`
	ConsultTopicID int       `json:"consult_topic_id" form:"consult_topic_id"`
	UserID         int       `json:"user_id" form:"user_id"`
	Time           string    `json:"time" form:"time"`
	Date           time.Time `json:"-" form:"date"`
	DateExport     string    `json:"date"`
	Question       string    `json:"question" form:"question"`
	IsNeedLetter   bool      `json:"is_need_letter" form:"is_need_letter"`
	IsConfirmed    bool      `json:"is_confirmed" form:"is_confirmed"`
}

type Consultations struct {
	Active   []Consultation `json:"active"`
	Finished []Consultation `json:"finished"`
}

type Slot struct {
	ID   int    `json:"id"`
	Time string `json:"time"`
	Date string `json:"date"`
}
