package config

import (
	"os"
	"time"
)

type Config struct {
	RabbitMQURL string
	SleepMs     time.Duration
}

func NewConfig() *Config {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	if len(rabbitMQURL) == 0 {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	sleepMs, err := time.ParseDuration(os.Getenv("SLEEP_MS"))

	if err != nil {
		sleepMs = 0
	}

	return &Config{
		RabbitMQURL: rabbitMQURL,
		SleepMs:     sleepMs,
	}
}
