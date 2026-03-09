package models

import "github.com/google/uuid"

type CrackTaskRequest struct {
	TaskId     uuid.UUID `json:"taskId"`
	RequestId  uuid.UUID `json:"requestId"`
	StartIndex int64     `json:"startIndex"`
	Count      int64     `json:"count"`
	TargetHash string    `json:"targetHash"`
	MaxLen     int64     `json:"maxLen"`
	Alphabet   string    `json:"alphabet"`
}

type GetResultRequest struct {
	RequestId uuid.UUID `json:"requestId"`
}

type CrackTaskResult struct {
	TaskId    uuid.UUID `json:"taskId"`
	RequestId uuid.UUID `json:"requestId"`
	Results   []string  `json:"results,omitempty"`
	Status    string    `json:"status"`
}

const (
	StatusDONE  string = "DONE"
	StatusERROR string = "ERROR"
)
