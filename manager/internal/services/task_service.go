package services

import (
	"github.com/TKaterinna/CrackHash/manager/internal/models"
	"github.com/TKaterinna/CrackHash/manager/internal/repo"
	"github.com/google/uuid"
)

type TaskService struct {
	repo               *repo.TaskRepo
	workerCountService *WorkerCountService
	taskSender         *TaskSender
}

func NewTaskService(repo *repo.TaskRepo, workerCountService *WorkerCountService, taskSender *TaskSender) *TaskService {
	return &TaskService{
		repo:               repo,
		workerCountService: workerCountService,
		taskSender:         taskSender,
	}
}

func (s *TaskService) Crack(req *models.HashCrackRequest) (uuid.UUID, error) {
	var err error

	id := uuid.New()

	if err = s.repo.SaveRequest(id); err != nil {
		return uuid.Nil, err
	}

	tasks := s.CreateTasks(req, id)

	go s.taskSender.Send(tasks)

	return id, nil
}

func (s *TaskService) CreateTasks(req *models.HashCrackRequest, requestId uuid.UUID) []*models.CrackTaskRequest {
	tasks := make([]*models.CrackTaskRequest, 0)

	partCount := s.workerCountService.GetWorkerCount()
	for i := range partCount {
		task := s.CreateTask(req, requestId, i, partCount)
		tasks = append(tasks, task)
	}

	return tasks
}

func (s *TaskService) CreateTask(req *models.HashCrackRequest, requestId uuid.UUID, partNumber int, partCount int) *models.CrackTaskRequest {
	return &models.CrackTaskRequest{
		RequestId:  requestId,
		PartNumber: partNumber,
		PartCount:  partCount,
		MaxLen:     req.MaxLength,
		CheckHash:  req.Hash,
		Alphabet:   req.Alphabet,
	}
}

func (s *TaskService) GetStatus(requestId uuid.UUID) (string, []string, error) {
	var status string
	var data []string
	var err error

	if status, data, err = s.repo.GetStatus(requestId); err != nil {
		return "", nil, err
	}

	return status, data, nil
}

func (s *TaskService) UpdateResult(req *models.CrackTaskResult) error {
	if err := s.repo.UpdateResult(req.RequestId, req.Results, true); err != nil {
		return err
	}

	return nil
}
