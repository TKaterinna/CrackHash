package repo

import (
	"sync"

	"github.com/google/uuid"
)

type CalcTasks struct {
	// RequestId  uuid.UUID key in map
	PartNumber int64
	PartCount  int64
	MaxLen     int64
	checkHash  string
	Data       string
}

type CalcRepo struct {
	// db mongo
	mx             sync.RWMutex
	cacheCalcTasks map[uuid.UUID]CalcTasks
}

func NewCalcRepo() *CalcRepo {
	return &CalcRepo{
		cacheCalcTasks: make(map[uuid.UUID]CalcTasks),
	}
}

func (r *CalcRepo) SaveTask(id uuid.UUID, partNumber int64, partCount int64, maxLen int64, checkHash string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.cacheCalcTasks[id] = CalcTasks{
		PartNumber: partNumber,
		PartCount:  partCount,
		MaxLen:     maxLen,
		checkHash:  checkHash,
		Data:       "",
	}

	return nil
}

func (r *CalcRepo) GetResult(id uuid.UUID) (string, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	data := r.cacheCalcTasks[id].Data

	return data, nil
}
