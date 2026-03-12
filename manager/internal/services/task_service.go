package services

import (
	"fmt"
	"math"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
	"github.com/TKaterinna/CrackHash/manager/internal/repo"
	"github.com/google/uuid"
)

type TaskService struct {
	repo        *repo.TaskRepo
	taskSender  *TaskSender
	CombForTask int64
}

func NewTaskService(repo *repo.TaskRepo, taskSender *TaskSender, combForTask int64) *TaskService {
	return &TaskService{
		repo:        repo,
		taskSender:  taskSender,
		CombForTask: combForTask,
	}
}

func (s *TaskService) Crack(req *models.HashCrackRequest) (uuid.UUID, error) {
	var err error

	id := uuid.New()

	tasks := s.CreateTasks(req, id)

	if err = s.repo.SaveRequest(id, tasks); err != nil {
		return uuid.Nil, err
	}

	go s.taskSender.Send(tasks)

	return id, nil
}

func (s *TaskService) CreateTasks(req *models.HashCrackRequest, requestId uuid.UUID) []*models.CrackTaskRequest {
	tasks := make([]*models.CrackTaskRequest, 0)

	wordsCount := s.WordsCount(int64(len(req.Alphabet)), req.MaxLength)
	partsCount := wordsCount / s.CombForTask
	if wordsCount%s.CombForTask != 0 {
		partsCount += 1
	}

	for i := range partsCount {
		task := s.CreateTask(req, requestId, i)
		tasks = append(tasks, task)
	}

	return tasks
}

func (s *TaskService) WordsCount(alphabetLen int64, maxLen int64) int64 {
	return int64(float64(alphabetLen) * (math.Pow(float64(alphabetLen), float64(maxLen)) - 1) / float64(alphabetLen-1))
}

func (s *TaskService) CreateTask(req *models.HashCrackRequest, requestId uuid.UUID, partNumber int64) *models.CrackTaskRequest {
	return &models.CrackTaskRequest{
		TaskId:     uuid.New(),
		RequestId:  requestId,
		StartIndex: partNumber * s.CombForTask,
		Count:      s.CombForTask,
		MaxLen:     req.MaxLength,
		TargetHash: req.Hash,
		Alphabet:   req.Alphabet,
	}
}

func (s *TaskService) GetStatus(requestId uuid.UUID) (string, []string, error) {
	var status string
	var results []string
	var err error

	if status, results, err = s.repo.GetStatus(requestId); err != nil {
		return "", nil, err
	}

	return status, results, nil
}

func (s *TaskService) UpdateResult(req *models.CrackTaskResult) error {
	if req.Status == models.StatusERROR {
		return fmt.Errorf("Worker error during crack hash")
	}
	if err := s.repo.UpdateResult(req.RequestId, req.TaskId, req.Results); err != nil {
		return err
	}

	return nil
}
