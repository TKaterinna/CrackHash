package main

import (
	"github.com/TKaterinna/CrackHash/worker/config"
	"github.com/TKaterinna/CrackHash/worker/internal/handlers"
	"github.com/TKaterinna/CrackHash/worker/internal/repo"
	"github.com/TKaterinna/CrackHash/worker/internal/routers"
	"github.com/TKaterinna/CrackHash/worker/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	config := config.NewConfig()

	router := gin.Default()

	calcRepo := repo.NewCalcRepo()
	resultSender := services.NewResultSender(config.ManagerPort)
	calcService := services.NewCalcService(calcRepo, resultSender)
	calcHandler := handlers.NewCalcHandler(calcService)

	routers.NewCalcRouter(router, calcHandler)

	router.Run(config.WorkerPort)
}
