package model

type Status struct {
	Status string `json:"status,omitempty"` // Статус
}

type ResponseError struct {
	Error string `json:"error,omitempty"`
}

const (
	StatusOK string = "OK"
)

type Version struct {
	Tier     string `json:"tier"`
	Version  string `json:"version"`
	Revision string `json:"revision"`
}

type PushMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
