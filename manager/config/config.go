package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ManagerPort  string
	WorkersCount int64
	WorkersPort  []string
	CombForTask  int64
}

func NewConfig() *Config {
	managerPort := os.Getenv("MANAGER_PORT")

	if len(managerPort) == 0 {
		managerPort = "8080"
	}

	managerPort = ":" + managerPort

	workersCountStr := os.Getenv("WORKERS_COUNT")
	var workersCount int
	var err error

	if len(workersCountStr) == 0 {
		log.Println("aaaaaaaaaaaaaaaaaaaaa WORKERS_COUNT empty")
		workersCount = 1
	} else {
		if workersCount, err = strconv.Atoi(workersCountStr); err != nil {
			workersCount = 1
		}
	}

	var workersPort []string
	for i := range workersCount {
		workerPort := os.Getenv("WORKER_PORT_" + strconv.Itoa(i))

		if len(workerPort) == 0 {
			workerPort = "8081"
		}

		workerPort = ":" + workerPort

		workersPort = append(workersPort, workerPort)
	}

	combForTaskStr := os.Getenv("COMB_FOR_TASK")
	var combForTask int

	if len(combForTaskStr) == 0 {
		combForTask = 100000
	} else {
		if combForTask, err = strconv.Atoi(combForTaskStr); err != nil {
			combForTask = 100000
		}
	}

	return &Config{
		ManagerPort:  managerPort,
		WorkersCount: int64(workersCount),
		WorkersPort:  workersPort,
		CombForTask:  int64(combForTask),
	}
}
