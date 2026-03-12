package repo

import (
	"fmt"
	"sync"
	"time"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
	"github.com/google/uuid"
)

type RequestStatus struct {
	// Id uuid.UUID key in map
	TasksCount int
	TasksReady map[uuid.UUID]bool
	StartTime  time.Time
	Status     string
	Results    []string
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
	errorDelay       time.Duration
}

func NewTaskRepo(errorDelay time.Duration) *TaskRepo {
	return &TaskRepo{
		mxStatus:         sync.RWMutex{},
		cacheStatus:      make(map[uuid.UUID]*RequestStatus),
		mxTasks:          sync.RWMutex{},
		cacheWorkerTasks: make(map[uuid.UUID]*WorkerTasks),
		errorDelay:       errorDelay,
	}
}

func (r *TaskRepo) SaveRequest(id uuid.UUID, tasks []*models.CrackTaskRequest) error {
	tasksReady := make(map[uuid.UUID]bool)

	for _, t := range tasks {
		tasksReady[t.TaskId] = false
	}

	r.mxStatus.Lock()
	defer r.mxStatus.Unlock()

	r.cacheStatus[id] = &RequestStatus{
		TasksCount: len(tasks),
		TasksReady: tasksReady,
		StartTime:  time.Now(),
		Status:     models.StatusInProgress,
		Results:    make([]string, 0),
	}

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

		if time.Now().Compare(entry.StartTime.Add(r.errorDelay)) >= 0 {
			entry.Status = models.StatusERROR
		}
	}

	return entry.Status, results, nil
}

func (r *TaskRepo) UpdateResult(reqId uuid.UUID, taskId uuid.UUID, results []string) error {
	r.mxStatus.Lock()
	defer r.mxStatus.Unlock()

	requestStatus := r.cacheStatus[reqId]

	if requestStatus.Status == models.StatusERROR {
		return fmt.Errorf("This task was canceled by timeout")
	}

	if requestStatus.TasksReady[taskId] {
		return fmt.Errorf("dublicated task result")
	}

	requestStatus.TasksReady[taskId] = true
	requestStatus.TasksCount--

	requestStatus.Results = append(requestStatus.Results, results...)

	if requestStatus.TasksCount == 0 {
		requestStatus.Status = models.StatusREADY
	}

	return nil
}
