package main

import (
	"context"
	"log"

	"github.com/TKaterinna/CrackHash/worker/config"
	"github.com/TKaterinna/CrackHash/worker/internal/repo"
	"github.com/TKaterinna/CrackHash/worker/internal/services"
)

func main() {
	config := config.NewConfig()

	calcRepo := repo.NewCalcRepo()
	rabbit_conn, err := services.RabbitMQConnect(config.RabbitMQURL)
	if err != nil {
		log.Panicf("Failed to connect RabbitMQ: %s", err)
	}
	defer rabbit_conn.Conn.Close()
	defer rabbit_conn.Channel.Close()

	err = rabbit_conn.SetupTopology()
	if err != nil {
		log.Panicf("Failed to setup topology: %s", err)
	}

	resultSender := services.NewResultSender(rabbit_conn)
	calcService := services.NewCalcService(calcRepo, resultSender)
	listener := services.NewCalcListener(rabbit_conn, calcService)

	listener.Listen(context.Background())
}
