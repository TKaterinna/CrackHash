package main

import (
	"context"
	"log"
	"net/http"

	"github.com/TKaterinna/CrackHash/manager/config"
	"github.com/TKaterinna/CrackHash/manager/internal/handlers"
	"github.com/TKaterinna/CrackHash/manager/internal/metrics"
	"github.com/TKaterinna/CrackHash/manager/internal/middleware"
	"github.com/TKaterinna/CrackHash/manager/internal/repo"
	"github.com/TKaterinna/CrackHash/manager/internal/routers"
	"github.com/TKaterinna/CrackHash/manager/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	config := config.NewConfig()

	router := gin.Default()

	router.Use(middleware.PrometheusMiddleware())

	taskRepo := repo.NewMongoTaskRepo(config.MongoURI, config.MongoDBName, config.ErrorDelay)

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

	log.Println("Recovering pending/queued tasks from MongoDB...")
	taskService.ResendQueuedTasks()

	listener := services.NewCalcListener(rabbit_conn, taskService)
	listener.Listen(context.Background())

	routers.NewTaskRouter(router, taskHandler)

	metrics.HealthStatus.Set(1)

	go func() {
		metricsRouter := http.NewServeMux()
		metricsRouter.Handle("/metrics", promhttp.Handler())

		log.Println("Metrics server starting on :8082")
		if err := http.ListenAndServe(":8082", metricsRouter); err != nil {
			log.Printf("Metrics server error: %v", err)
			metrics.HealthStatus.Set(0)
		}
	}()

	if err := router.Run(config.ManagerPort); err != nil {
		log.Panicf("Failed to start manager: %v", err)
	}
}
