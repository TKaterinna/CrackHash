package main

import (
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

	taskRepo := repo.NewTaskRepo()
	taskSender := services.NewTaskSender(config.WorkersCount, config.WorkersPort)
	taskService := services.NewTaskService(taskRepo, taskSender, config.CombForTask)
	taskHandler := handlers.NewTaskHandler(taskService)

	routers.NewTaskRouter(router, taskHandler)

	router.Run(config.ManagerPort)
}
