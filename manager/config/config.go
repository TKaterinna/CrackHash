package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ManagerPort  string
	WorkersCount int64
	RabbitMQURL  string
	CombForTask  int64
	ErrorDelay   time.Duration
	MongoURI     string
	MongoDBName  string
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

	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	if len(rabbitMQURL) == 0 {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
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

	errorDelayStr := os.Getenv("ERROR_DELAY")
	var errorDelay time.Duration

	if len(errorDelayStr) == 0 {
		if errorDelay, err = time.ParseDuration("5m"); err != nil {
			log.Panic("Can not parse default error delay to duration")
		}
	} else {
		if errorDelay, err = time.ParseDuration(errorDelayStr); err != nil {
			log.Panic("Can not parse error delay to duration")
		}
	}

	mongoURI := os.Getenv("MONGODB_URI")

	if len(mongoURI) == 0 {
		mongoURI = "mongodb://mongodb-0.mongodb-headless.default.svc.cluster.local:27017,mongodb-1.mongodb-headless.default.svc.cluster.local:27017,mongodb-2.mongodb-headless.default.svc.cluster.local:27017/crackhash_db?replicaSet=rs0"
	}

	mongoDBName := os.Getenv("MONGO_DB_NAME")

	if len(mongoDBName) == 0 {
		mongoDBName = "crackhash_db"
	}

	return &Config{
		ManagerPort:  managerPort,
		WorkersCount: int64(workersCount),
		RabbitMQURL:  rabbitMQURL,
		CombForTask:  int64(combForTask),
		ErrorDelay:   errorDelay,
		MongoURI:     mongoURI,
		MongoDBName:  mongoDBName,
	}
}
