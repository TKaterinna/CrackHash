package main

import (
	"context"
	"log"

	"github.com/TKaterinna/CrackHash/manager/config"
	"github.com/TKaterinna/CrackHash/manager/internal/handlers"
	"github.com/TKaterinna/CrackHash/manager/internal/repo"
	"github.com/TKaterinna/CrackHash/manager/internal/routers"
	"github.com/TKaterinna/CrackHash/manager/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {

	config := config.NewConfig()

	router := gin.Default()

	taskRepo := repo.NewTaskRepo(config.ErrorDelay)

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

	taskSender := services.NewTaskSender(rabbit_conn)
	taskService := services.NewTaskService(taskRepo, taskSender, config.CombForTask)
	taskHandler := handlers.NewTaskHandler(taskService)
	listener := services.NewCalcListener(rabbit_conn, taskService)
	listener.Listen(context.Background())

	routers.NewTaskRouter(router, taskHandler)

	router.Run(config.ManagerPort)
}
