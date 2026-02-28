package config

import "os"

type Config struct {
	ManagerPort string
	WorkerPort  string
}

func NewConfig() *Config {
	managerPort := os.Getenv("MANAGER_PORT")

	if len(managerPort) == 0 {
		managerPort = "8080"
	}

	managerPort = ":" + managerPort

	workerPort := os.Getenv("WORKER_PORT")

	if len(workerPort) == 0 {
		workerPort = "8081"
	}

	workerPort = ":" + workerPort

	return &Config{
		WorkerPort:  workerPort,
		ManagerPort: managerPort,
	}
}
