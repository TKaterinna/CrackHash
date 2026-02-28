package models

import "github.com/google/uuid"

type CrackTaskRequest struct {
	RequestId  uuid.UUID `json:"requestId"`
	PartNumber int       `json:"partNumber"`
	PartCount  int       `json:"partCount"`
	MaxLen     int       `json:"maxLen"`
	CheckHash  string    `json:"checkHash"`
	Alphabet   string    `json:"alphabet"`
}

type GetResultRequest struct {
	RequestId uuid.UUID `json:"requestId"`
}

type CrackTaskResult struct {
	RequestId uuid.UUID `json:"requestId"`
	Results   []string  `json:"results,omitempty"`
	Status    string    `json:"status"`
}

const (
	StatusDONE  string = "DONE"
	StatusERROR string = "ERROR"
)
