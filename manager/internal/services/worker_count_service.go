package services

type WorkerCountService struct {
	defaultWorkerCount int
}

func NewWorkerCountService() *WorkerCountService {
	return &WorkerCountService{
		defaultWorkerCount: 1,
	}
}

func (w *WorkerCountService) GetWorkerCount() int {
	return w.defaultWorkerCount
}
