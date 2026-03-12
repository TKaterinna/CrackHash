package config

import "os"

type Config struct {
	RabbitMQURL string
}

func NewConfig() *Config {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	if len(rabbitMQURL) == 0 {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	return &Config{
		RabbitMQURL: rabbitMQURL,
	}
}
