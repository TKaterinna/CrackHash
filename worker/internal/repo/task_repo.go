package repo

import (
	"sync"

	"github.com/google/uuid"
)

type CalcTasks struct {
	// RequestId  uuid.UUID key in map
	PartNumber int
	PartCount  int
	MaxLen     int
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

func (r *CalcRepo) SaveTask(id uuid.UUID, partNumber int, partCount int, maxLen int, checkHash string) error {
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
