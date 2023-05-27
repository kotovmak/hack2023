package model

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