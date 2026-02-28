package repo

import (
	"sync"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
	"github.com/google/uuid"
)

type RequestStatus struct {
	// Id uuid.UUID key in map
	Status string
	Data   []string
}

type WorkerTasks struct {
	// WorkerId uuid.UUID key in map
	RequestId  uuid.UUID
	PartNumber int
	PartCount  int
	Alphabet   string
	MaxLen     int
	CheckHash  string
}

type TaskRepo struct {
	// db mongo
	mxStatus         sync.RWMutex
	cacheStatus      map[uuid.UUID]*RequestStatus
	mxTasks          sync.RWMutex
	cacheWorkerTasks map[uuid.UUID]*WorkerTasks
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		mxStatus:         sync.RWMutex{},
		cacheStatus:      make(map[uuid.UUID]*RequestStatus),
		mxTasks:          sync.RWMutex{},
		cacheWorkerTasks: make(map[uuid.UUID]*WorkerTasks),
	}
}

func (r *TaskRepo) SaveRequest(id uuid.UUID) error {
	r.mxStatus.Lock()
	defer r.mxStatus.Unlock()

	r.cacheStatus[id] = &RequestStatus{Status: models.StatusInProgress}

	return nil
}

func (r *TaskRepo) GetStatus(id uuid.UUID) (string, []string, error) {
	r.mxStatus.RLock()
	defer r.mxStatus.RUnlock()

	entry := r.cacheStatus[id]

	var data []string
	if entry.Status == models.StatusREADY {
		data = entry.Data
	} else {
		data = nil
	}

	return entry.Status, data, nil
}

func (r *TaskRepo) UpdateResult(id uuid.UUID, results []string, isEnd bool) error {
	r.mxStatus.Lock()
	defer r.mxStatus.Unlock()

	requestStatus := r.cacheStatus[id]
	requestStatus.Data = append(requestStatus.Data, results...)

	if isEnd {
		requestStatus.Status = models.StatusREADY
	}

	return nil
}
