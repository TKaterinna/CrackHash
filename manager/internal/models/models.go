package models

import "github.com/google/uuid"

type HashCrackRequest struct {
	Hash      string `json:"hash"`
	MaxLength int    `json:"maxLength"`
	Alphabet  string `json:"alphabet"`
}

type HashCrackResponse struct {
	RequestId uuid.UUID `json:"requestId"`
}

type HashStatusResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data,omitempty"`
}

type CrackTaskRequest struct {
	RequestId  uuid.UUID `json:"requestId"`
	PartNumber int       `json:"partNumber"`
	PartCount  int       `json:"partCount"`
	MaxLen     int       `json:"maxLen"`
	CheckHash  string    `json:"checkHash"`
	Alphabet   string    `json:"alphabet"`
}

type CrackTaskResult struct {
	RequestId uuid.UUID `json:"requestId"`
	Results   []string  `json:"results,omitempty"`
	Status    string    `json:"status"`
}

const (
	StatusInProgress string = "IN_PROGRESS"
	StatusREADY      string = "READY"
	StatusERROR      string = "ERROR"
)
