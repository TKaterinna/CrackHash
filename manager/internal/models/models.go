package models

import "github.com/google/uuid"

type HashCrackRequest struct {
	Hash      string `json:"hash"`
	MaxLength int64  `json:"maxLength"`
	Alphabet  string `json:"alphabet"`
}

type HashCrackResponse struct {
	RequestId uuid.UUID `json:"requestId"`
}

type HashStatusResponse struct {
	Status  string    `json:"status"`
	Results *[]string `json:"results,omitempty"`
}

type CrackTaskRequest struct {
	TaskId     uuid.UUID `json:"taskId"`
	RequestId  uuid.UUID `json:"requestId"`
	StartIndex int64     `json:"startIndex"`
	Count      int64     `json:"count"`
	TargetHash string    `json:"targetHash"`
	MaxLen     int64     `json:"maxLen"`
	Alphabet   string    `json:"alphabet"`
}

type CrackTaskResult struct {
	TaskId    uuid.UUID `json:"taskId"`
	RequestId uuid.UUID `json:"requestId"`
	Results   []string  `json:"results,omitempty"`
	Status    string    `json:"status"`
}

const (
	StatusInProgress string = "IN_PROGRESS"
	StatusREADY      string = "READY"
	StatusERROR      string = "ERROR"
)
