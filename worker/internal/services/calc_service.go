package services

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"sync/atomic"
	"time"

	"github.com/TKaterinna/CrackHash/worker/internal/metrics"
	"github.com/TKaterinna/CrackHash/worker/internal/models"
	"github.com/TKaterinna/CrackHash/worker/internal/repo"
	"github.com/google/uuid"
)

type CalcService struct {
	repo         *repo.CalcRepo
	resultSender *ResultSender
	sleepMs      time.Duration
	activeTasks  atomic.Int64
}

func NewCalcService(repo *repo.CalcRepo, resultSender *ResultSender, sleepMs time.Duration) *CalcService {
	return &CalcService{
		repo:         repo,
		resultSender: resultSender,
		sleepMs:      sleepMs,
	}
}

func (s *CalcService) Save(req *models.CrackTaskRequest) error {
	// if err = s.repo.SaveTask(req.RequestId, req.PartNumber, req.PartCount, req.MaxLen, req.CheckHash); err != nil {
	// 	return err
	// }

	print("Save")

	s.work(req)

	return nil
}

func (s *CalcService) GetResult(requestId uuid.UUID) (string, error) {
	var data string
	var err error

	if data, err = s.repo.GetResult(requestId); err != nil {
		return "", err
	}

	return data, nil
}

func (s *CalcService) getMD5Hash(word string) string {
	hash := md5.Sum([]byte(word))
	return hex.EncodeToString(hash[:])
}

func (s *CalcService) checkWord(word string, checkHash string) bool {
	curHash := s.getMD5Hash(word)

	if curHash == checkHash {
		return true
	}

	return false
}

func (s *CalcService) updateWorkerStatus(activeCount int64) {
	metrics.ActiveTasks.Set(float64(activeCount))

	if activeCount > 0 {
		metrics.WorkerStatus.Set(1)
	} else {
		metrics.WorkerStatus.Set(0)
	}
}

func (s *CalcService) work(req *models.CrackTaskRequest) {
	start := time.Now()

	current := s.activeTasks.Add(1)
	s.updateWorkerStatus(current)

	time.Sleep(s.sleepMs)

	defer func() {
		current := s.activeTasks.Add(-1)
		s.updateWorkerStatus(current)
		metrics.TaskDuration.Observe(time.Since(start).Seconds())
	}()

	var err error
	var words []string
	var wg *WordGenerator

	if wg, err = NewWordGenerator(req); err != nil {
		log.Printf("Task %s failed at generator init: %v", req.TaskId, err)

		res := &models.CrackTaskResult{
			TaskId:    req.TaskId,
			RequestId: req.RequestId,
			Results:   nil,
			Status:    models.StatusERROR,
		}
		s.resultSender.Send(res)

		metrics.TasksTotal.WithLabelValues("error").Inc()
		return
	}

	log.Printf("START work on task %s", req.TaskId)

	for {
		var word string
		var isNotEnd bool
		if word, isNotEnd = wg.Next(); !isNotEnd {
			break
		}

		if s.checkWord(word, req.TargetHash) {
			log.Printf("✓ FOUND: %s", word)
			words = append(words, word)
		}
	}

	res := &models.CrackTaskResult{
		TaskId:    req.TaskId,
		RequestId: req.RequestId,
		Results:   words,
		Status:    models.StatusDONE,
	}
	s.resultSender.Send(res)

	metrics.TasksTotal.WithLabelValues("done").Inc()
	log.Printf("✓ Task %s completed", req.TaskId)
}
