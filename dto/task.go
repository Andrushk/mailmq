package dto

type Task struct {
	To      []string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}