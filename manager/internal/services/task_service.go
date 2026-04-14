package services

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
	"github.com/google/uuid"
)

type TaskService struct {
	repo        TaskRepo
	taskSender  *TaskSender
	CombForTask int64
	mu          sync.Mutex
}

func NewTaskService(repo TaskRepo, taskSender *TaskSender, combForTask int64) *TaskService {
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
		return uuid.Nil, fmt.Errorf("save request to mongo: %w", err)
	}

	go func() {
		if err := s.taskSender.Send(tasks); err != nil {
			log.Printf("Initial send failed for request %s: %v. Tasks remain QUEUED.", id, err)
			return
		}

		if err := s.repo.UpdateRequestStatus(id, models.StatusInProgress); err != nil {
			log.Printf("Failed to update status for %s: %v", id, err)
		}
	}()

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

func (s *TaskService) ResendQueuedTasks() {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.repo.GetQueuedTasks()
	if err != nil {
		log.Printf("Reconnect recovery: failed to fetch queued tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		return
	}

	log.Printf("Connection restored: resending %d queued tasks...", len(tasks))
	if err := s.taskSender.Send(tasks); err != nil {
		log.Printf("Resend failed: %v. Tasks remain QUEUED for next attempt.", err)
		return
	}

	reqIds := make(map[uuid.UUID]struct{})
	for _, t := range tasks {
		reqIds[t.RequestId] = struct{}{}
	}
	for reqId := range reqIds {
		_ = s.repo.UpdateRequestStatus(reqId, models.StatusInProgress)
	}
	log.Printf("Queued tasks successfully resent and status updated.")
}
