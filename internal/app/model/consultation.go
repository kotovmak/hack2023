package model

type TypeList struct {
	NadzonOrgans  []NadzonOrgan  `json:"nadzor_organs"`
	Services      []Service      `json:"services"`
	ControlTypes  []ControlType  `json:"control_types"`
	ConsultTopics []ConsultTopic `json:"consult_topics"`
	PravActs      []PravAct      `json:"prav_acts"`
}

type NadzonOrgan struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	ConsultTopic []ConsultTopic `json:"consult_topics,omitempty"`
}

type Service struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ControlType struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	NadzonOrgan NadzonOrgan `json:"nadzor_organ"`
}

type ConsultTopic struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	ControlType ControlType `json:"control_type"`
}

type PravAct struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	ControlType ControlType `json:"control_type"`
}
