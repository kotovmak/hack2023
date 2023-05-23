package model

import "time"

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
	ID             int       `json:"id"`
	NadzonOrganID  int       `json:"nadzor_organ_ID"`
	ControlTypeID  int       `json:"control_type_id"`
	ConsultTopicID int       `json:"consult_topic_id"`
	UserID         int       `json:"user_id"`
	Time           string    `json:"time"`
	Date           time.Time `json:"date"`
	IsNeedLetter   bool      `json:"is_need_letter"`
	IsConfirmed    bool      `json:"is_confirmed"`
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
