package repo

import (
	"sync"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
	"github.com/google/uuid"
)

type RequestStatus struct {
	// Id uuid.UUID key in map
	Status  string
	Results []string
}

type WorkerTasks struct {
	// TaskId uuid.UUID key in map
	RequestId  uuid.UUID
	StartIndex int64
	Count      int64
	Alphabet   string
	MaxLen     int64
	TargetHash string
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

	var results []string
	if entry.Status == models.StatusREADY {
		results = entry.Results
	} else {
		results = nil
	}

	return entry.Status, results, nil
}

func (r *TaskRepo) UpdateResult(id uuid.UUID, results []string, isEnd bool) error {
	r.mxStatus.Lock()
	defer r.mxStatus.Unlock()

	requestStatus := r.cacheStatus[id]
	requestStatus.Results = append(requestStatus.Results, results...)

	if isEnd {
		requestStatus.Status = models.StatusREADY
	}

	return nil
}
