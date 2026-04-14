package services

import (
	"github.com/TKaterinna/CrackHash/manager/internal/models"
	"github.com/google/uuid"
)

type TaskRepo interface {
	SaveRequest(id uuid.UUID, tasks []*models.CrackTaskRequest) error
	GetStatus(id uuid.UUID) (string, []string, error)
	UpdateResult(reqId uuid.UUID, taskId uuid.UUID, results []string) error
	UpdateRequestStatus(id uuid.UUID, status string) error
	GetQueuedTasks() ([]*models.CrackTaskRequest, error)
}
